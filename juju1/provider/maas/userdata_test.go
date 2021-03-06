// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package maas_test

import (
	"encoding/base64"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils"
	gc "gopkg.in/check.v1"

	"github.com/juju/1.25-upgrade/juju1/cloudconfig/providerinit/renderers"
	"github.com/juju/1.25-upgrade/juju1/provider/maas"
	"github.com/juju/1.25-upgrade/juju1/testing"
	"github.com/juju/1.25-upgrade/juju1/version"
)

type RenderersSuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&RenderersSuite{})

func (s *RenderersSuite) TestMAASUnix(c *gc.C) {
	renderer := maas.MAASRenderer{}
	data := []byte("test")
	result, err := renderer.EncodeUserdata(data, version.Ubuntu)
	c.Assert(err, jc.ErrorIsNil)
	expected := base64.StdEncoding.EncodeToString(utils.Gzip(data))
	c.Assert(string(result), jc.DeepEquals, expected)

	data = []byte("test")
	result, err = renderer.EncodeUserdata(data, version.CentOS)
	c.Assert(err, jc.ErrorIsNil)
	expected = base64.StdEncoding.EncodeToString(utils.Gzip(data))
	c.Assert(string(result), jc.DeepEquals, expected)
}

func (s *RenderersSuite) TestMAASWindows(c *gc.C) {
	renderer := maas.MAASRenderer{}
	data := []byte("test")
	result, err := renderer.EncodeUserdata(data, version.Windows)
	c.Assert(err, jc.ErrorIsNil)
	expected := base64.StdEncoding.EncodeToString(renderers.WinEmbedInScript(data))
	c.Assert(string(result), jc.DeepEquals, expected)
}

func (s *RenderersSuite) TestMAASUnknownOS(c *gc.C) {
	renderer := maas.MAASRenderer{}
	result, err := renderer.EncodeUserdata(nil, version.Arch)
	c.Assert(result, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Cannot encode userdata for OS: Arch")
}
