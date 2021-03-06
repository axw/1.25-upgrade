// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package vsphere_test

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/juju/mutex"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/clock"
	gc "gopkg.in/check.v1"

	"github.com/juju/1.25-upgrade/juju2/environs"
	"github.com/juju/1.25-upgrade/juju2/environs/imagemetadata"
	"github.com/juju/1.25-upgrade/juju2/provider/vsphere"
	coretesting "github.com/juju/1.25-upgrade/juju2/testing"
)

type ProviderFixture struct {
	testing.IsolationSuite
	dialStub    testing.Stub
	client      *mockClient
	ovaCacheDir string
	provider    environs.EnvironProvider
}

func (s *ProviderFixture) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
	s.dialStub.ResetCalls()
	s.client = &mockClient{}
	s.ovaCacheDir = c.MkDir()
	s.provider = vsphere.NewEnvironProvider(vsphere.EnvironProviderConfig{
		Dial:        newMockDialFunc(&s.dialStub, s.client),
		OVACacheDir: s.ovaCacheDir,
		OVACacheLocker: vsphere.NewMutexCacheLocker(mutex.Spec{
			Name:  "juju-vsphere-testing",
			Clock: clock.WallClock,
			Delay: time.Second,
		}),
	})
}

type EnvironFixture struct {
	ProviderFixture
	imageServer         *httptest.Server
	imageServerRequests []*http.Request
	env                 environs.Environ
}

func (s *EnvironFixture) SetUpTest(c *gc.C) {
	s.ProviderFixture.SetUpTest(c)

	s.imageServerRequests = nil
	s.imageServer = serveImageMetadata(&s.imageServerRequests)
	s.AddCleanup(func(*gc.C) {
		s.imageServer.Close()
	})

	env, err := s.provider.Open(environs.OpenParams{
		Cloud: fakeCloudSpec(),
		Config: fakeConfig(c, coretesting.Attrs{
			"image-metadata-url": s.imageServer.URL,
		}),
	})
	c.Assert(err, jc.ErrorIsNil)
	s.env = env

	// Make sure we don't fall back to the public image sources.
	s.PatchValue(&imagemetadata.DefaultUbuntuBaseURL, "")
	s.PatchValue(&imagemetadata.DefaultJujuBaseURL, "")
}

func serveImageMetadata(requests *[]*http.Request) *httptest.Server {
	index := `{
  "index": {
     "com.ubuntu.cloud:released:download": {
      "datatype": "image-downloads", 
      "path": "streams/v1/com.ubuntu.cloud:released:download.json", 
      "updated": "Tue, 24 Feb 2015 10:16:54 +0000", 
      "products": ["com.ubuntu.cloud:server:14.04:amd64"], 
      "format": "products:1.0"
    }
  }, 
  "updated": "Tue, 24 Feb 2015 14:14:24 +0000", 
  "format": "index:1.0"
}`

	download := fmt.Sprintf(`{
  "updated": "Thu, 05 Mar 2015 12:14:40 +0000", 
  "license": "http://www.canonical.com/intellectual-property-policy", 
  "format": "products:1.0", 
  "datatype": "image-downloads", 
  "products": {
    "com.ubuntu.cloud:server:14.04:amd64": {
      "release": "trusty", 
      "version": "14.04", 
      "arch": "amd64", 
      "versions": {
        "20150305": {
          "items": {
            "ova": {
              "size": 7196, 
              "path": "server/releases/trusty/release-20150305/ubuntu-14.04-server-cloudimg-amd64.ova",
              "ftype": "ova", 
              "sha256": "%s", 
              "md5": "00662c59ca52558e7a3bb9a67d194730"
            }
          }      
        }
      }
    }
  }
}`, fakeOvaSha256)

	files := map[string][]byte{
		"/streams/v1/index.json":                                                          []byte(index),
		"/streams/v1/com.ubuntu.cloud:released:download.json":                             []byte(download),
		"/server/releases/trusty/release-20150305/ubuntu-14.04-server-cloudimg-amd64.ova": fakeOva,
	}
	mux := http.NewServeMux()
	for path := range files {
		mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			*requests = append(*requests, req)
			w.Write(files[req.URL.Path])
		})
	}
	return httptest.NewServer(mux)
}

var (
	fakeOva       []byte
	fakeOvaSha256 string
)

func init() {
	buf := new(bytes.Buffer)
	hash := sha256.New()
	tw := tar.NewWriter(io.MultiWriter(buf, hash))
	var files = []struct{ Name, Body string }{
		{"ubuntu-14.04-server-cloudimg-amd64.ovf", "FakeOvfContent"},
		{"ubuntu-14.04-server-cloudimg-amd64.vmdk", "FakeVmdkContent"},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte(file.Body))
	}
	tw.Close()
	fakeOva = buf.Bytes()
	fakeOvaSha256 = fmt.Sprintf("%x", hash.Sum(nil))
}
