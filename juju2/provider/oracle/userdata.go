// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package oracle

import (
	"github.com/juju/errors"
	jujuos "github.com/juju/utils/os"

	"github.com/juju/1.25-upgrade/juju2/cloudconfig/cloudinit"
	"github.com/juju/1.25-upgrade/juju2/cloudconfig/providerinit/renderers"
)

// OracleRenderer implements the renderers.ProviderRenderer interface
type OracleRenderer struct{}

// Renderer is defined in the renderers.ProviderRenderer interface
func (OracleRenderer) Render(cfg cloudinit.CloudConfig, os jujuos.OSType) ([]byte, error) {
	switch os {
	case jujuos.Ubuntu:
		return renderers.RenderYAML(cfg)
	default:
		return nil, errors.Errorf("Cannot encode userdata for OS: %s", os.String())
	}
}
