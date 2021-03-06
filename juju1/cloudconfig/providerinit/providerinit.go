// Copyright 2013, 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

// This package offers userdata in a gzipped format to be used by different
// cloud providers
package providerinit

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"

	"github.com/juju/1.25-upgrade/juju1/cloudconfig"
	"github.com/juju/1.25-upgrade/juju1/cloudconfig/cloudinit"
	"github.com/juju/1.25-upgrade/juju1/cloudconfig/instancecfg"
	"github.com/juju/1.25-upgrade/juju1/cloudconfig/providerinit/renderers"
	"github.com/juju/1.25-upgrade/juju1/version"
)

var logger = loggo.GetLogger("juju.cloudconfig.providerinit")

func configureCloudinit(icfg *instancecfg.InstanceConfig, cloudcfg cloudinit.CloudConfig) (cloudconfig.UserdataConfig, error) {
	// When bootstrapping, we only want to apt-get update/upgrade
	// and setup the SSH keys. The rest we leave to cloudinit/sshinit.
	udata, err := cloudconfig.NewUserdataConfig(icfg, cloudcfg)
	if err != nil {
		return nil, err
	}
	if icfg.Bootstrap {
		err = udata.ConfigureBasic()
		if err != nil {
			return nil, err
		}
		return udata, nil
	}
	err = udata.Configure()
	if err != nil {
		return nil, err
	}
	return udata, nil
}

// ComposeUserData fills out the provided cloudinit configuration structure
// so it is suitable for initialising a machine with the given configuration,
// and then renders it and encodes it using the supplied renderer.
// When calling ComposeUserData a encoding implementation must be chosen from
// the providerinit/encoders package according to the need of the provider.
//
// If the provided cloudcfg is nil, a new one will be created internally.
func ComposeUserData(icfg *instancecfg.InstanceConfig, cloudcfg cloudinit.CloudConfig, renderer renderers.ProviderRenderer) ([]byte, error) {
	if cloudcfg == nil {
		var err error
		cloudcfg, err = cloudinit.New(icfg.Series)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	_, err := configureCloudinit(icfg, cloudcfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	operatingSystem, err := version.GetOSFromSeries(icfg.Series)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// This might get replaced by a renderer.RenderUserdata which will either
	// render it as YAML or Bash since some CentOS images might ship without cloudnit
	udata, err := cloudcfg.RenderYAML()
	if err != nil {
		return nil, errors.Trace(err)
	}
	udata, err = renderer.EncodeUserdata(udata, operatingSystem)
	if err != nil {
		return nil, errors.Trace(err)
	}
	logger.Tracef("Generated cloud init:\n%s", string(udata))
	return udata, err
}
