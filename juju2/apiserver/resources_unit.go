// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/juju/errors"
	"gopkg.in/juju/names.v2"

	"github.com/juju/1.25-upgrade/juju2/apiserver/params"
	"github.com/juju/1.25-upgrade/juju2/resource"
	"github.com/juju/1.25-upgrade/juju2/resource/api"
	"github.com/juju/1.25-upgrade/juju2/state"
)

// ResourcesHandler is the HTTP handler for unit agent downloads of
// resources.
type UnitResourcesHandler struct {
	NewOpener func(*http.Request, ...string) (resource.Opener, state.StatePoolReleaser, error)
}

// ServeHTTP implements http.Handler.
func (h *UnitResourcesHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		opener, closer, err := h.NewOpener(req, names.UnitTagKind)
		if err != nil {
			api.SendHTTPError(resp, err)
			return
		}
		defer closer()

		name := req.URL.Query().Get(":resource")
		opened, err := opener.OpenResource(name)
		if err != nil {
			logger.Errorf("cannot fetch resource reader: %v", err)
			api.SendHTTPError(resp, err)
			return
		}
		defer opened.Close()

		hdr := resp.Header()
		hdr.Set("Content-Type", params.ContentTypeRaw)
		hdr.Set("Content-Length", fmt.Sprint(opened.Size))
		hdr.Set("Content-Sha384", opened.Fingerprint.String())

		resp.WriteHeader(http.StatusOK)
		if _, err := io.Copy(resp, opened); err != nil {
			// We cannot use SendHTTPError here, so we log the error
			// and move on.
			logger.Errorf("unable to complete stream for resource: %v", err)
			return
		}
	default:
		api.SendHTTPError(resp, errors.MethodNotAllowedf("unsupported method: %q", req.Method))
	}
}
