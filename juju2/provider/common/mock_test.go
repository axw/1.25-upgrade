// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"io"

	"github.com/juju/1.25-upgrade/juju2/cloudconfig/instancecfg"
	"github.com/juju/1.25-upgrade/juju2/constraints"
	"github.com/juju/1.25-upgrade/juju2/environs"
	"github.com/juju/1.25-upgrade/juju2/environs/config"
	"github.com/juju/1.25-upgrade/juju2/environs/simplestreams"
	"github.com/juju/1.25-upgrade/juju2/environs/storage"
	"github.com/juju/1.25-upgrade/juju2/instance"
	"github.com/juju/1.25-upgrade/juju2/network"
	"github.com/juju/1.25-upgrade/juju2/provider/common"
	jujustorage "github.com/juju/1.25-upgrade/juju2/storage"
	"github.com/juju/1.25-upgrade/juju2/tools"
)

type allInstancesFunc func() ([]instance.Instance, error)
type instancesFunc func([]instance.Id) ([]instance.Instance, error)
type startInstanceFunc func(string, constraints.Value, []string, tools.List, *instancecfg.InstanceConfig) (instance.Instance, *instance.HardwareCharacteristics, []network.InterfaceInfo, error)
type stopInstancesFunc func([]instance.Id) error
type getToolsSourcesFunc func() ([]simplestreams.DataSource, error)
type configFunc func() *config.Config
type setConfigFunc func(*config.Config) error

type mockEnviron struct {
	storage          storage.Storage
	allInstances     allInstancesFunc
	instances        instancesFunc
	startInstance    startInstanceFunc
	stopInstances    stopInstancesFunc
	getToolsSources  getToolsSourcesFunc
	config           configFunc
	setConfig        setConfigFunc
	storageProviders jujustorage.StaticProviderRegistry
	environs.Environ // stub out other methods with panics
}

func (env *mockEnviron) Storage() storage.Storage {
	return env.storage
}

func (env *mockEnviron) AllInstances() ([]instance.Instance, error) {
	return env.allInstances()
}

func (env *mockEnviron) Instances(ids []instance.Id) ([]instance.Instance, error) {
	return env.instances(ids)
}

func (env *mockEnviron) StartInstance(args environs.StartInstanceParams) (*environs.StartInstanceResult, error) {
	inst, hw, networkInfo, err := env.startInstance(
		args.Placement,
		args.Constraints,
		nil,
		args.Tools,
		args.InstanceConfig,
	)
	if err != nil {
		return nil, err
	}
	return &environs.StartInstanceResult{
		Instance:    inst,
		Hardware:    hw,
		NetworkInfo: networkInfo,
	}, nil
}

func (env *mockEnviron) StopInstances(ids ...instance.Id) error {
	return env.stopInstances(ids)
}

func (env *mockEnviron) Config() *config.Config {
	return env.config()
}

func (env *mockEnviron) SetConfig(cfg *config.Config) error {
	if env.setConfig != nil {
		return env.setConfig(cfg)
	}
	return nil
}

func (env *mockEnviron) GetToolsSources() ([]simplestreams.DataSource, error) {
	if env.getToolsSources != nil {
		return env.getToolsSources()
	}
	datasource := storage.NewStorageSimpleStreamsDataSource("test cloud storage", env.Storage(), storage.BaseToolsPath, simplestreams.SPECIFIC_CLOUD_DATA, false)
	return []simplestreams.DataSource{datasource}, nil
}

func (env *mockEnviron) StorageProviderTypes() ([]jujustorage.ProviderType, error) {
	return env.storageProviders.StorageProviderTypes()
}

func (env *mockEnviron) StorageProvider(t jujustorage.ProviderType) (jujustorage.Provider, error) {
	return env.storageProviders.StorageProvider(t)
}

type availabilityZonesFunc func() ([]common.AvailabilityZone, error)
type instanceAvailabilityZoneNamesFunc func([]instance.Id) ([]string, error)

type mockZonedEnviron struct {
	mockEnviron
	availabilityZones             availabilityZonesFunc
	instanceAvailabilityZoneNames instanceAvailabilityZoneNamesFunc
}

func (env *mockZonedEnviron) AvailabilityZones() ([]common.AvailabilityZone, error) {
	return env.availabilityZones()
}

func (env *mockZonedEnviron) InstanceAvailabilityZoneNames(ids []instance.Id) ([]string, error) {
	return env.instanceAvailabilityZoneNames(ids)
}

type mockInstance struct {
	id                string
	addresses         []network.Address
	addressesErr      error
	dnsName           string
	dnsNameErr        error
	status            instance.InstanceStatus
	instance.Instance // stub out other methods with panics
}

func (inst *mockInstance) Id() instance.Id {
	return instance.Id(inst.id)
}

func (inst *mockInstance) Status() instance.InstanceStatus {
	return inst.status
}

func (inst *mockInstance) Addresses() ([]network.Address, error) {
	return inst.addresses, inst.addressesErr
}

type mockStorage struct {
	storage.Storage
	putErr       error
	removeAllErr error
}

func (stor *mockStorage) Put(name string, reader io.Reader, size int64) error {
	if stor.putErr != nil {
		return stor.putErr
	}
	return stor.Storage.Put(name, reader, size)
}

func (stor *mockStorage) RemoveAll() error {
	if stor.removeAllErr != nil {
		return stor.removeAllErr
	}
	return stor.Storage.RemoveAll()
}

type mockAvailabilityZone struct {
	name      string
	available bool
}

func (z *mockAvailabilityZone) Name() string {
	return z.name
}

func (z *mockAvailabilityZone) Available() bool {
	return z.available
}
