// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package networkingcommon

import (
	"github.com/juju/errors"

	"github.com/juju/1.25-upgrade/juju2/apiserver/params"
	"github.com/juju/1.25-upgrade/juju2/network"
	providercommon "github.com/juju/1.25-upgrade/juju2/provider/common"
	"github.com/juju/1.25-upgrade/juju2/state"
	"github.com/juju/1.25-upgrade/juju2/state/stateenvirons"
)

// NOTE: All of the following code is only tested with a feature test.

// subnetShim forwards and adapts state.Subnets methods to BackingSubnet.
type subnetShim struct {
	subnet *state.Subnet
}

func (s *subnetShim) CIDR() string {
	return s.subnet.CIDR()
}

func (s *subnetShim) VLANTag() int {
	return s.subnet.VLANTag()
}

func (s *subnetShim) ProviderNetworkId() network.Id {
	return s.subnet.ProviderNetworkId()
}

func (s *subnetShim) ProviderId() network.Id {
	return s.subnet.ProviderId()
}

func (s *subnetShim) AvailabilityZones() []string {
	// TODO(dimitern): Add multiple zones to state.Subnet.
	return []string{s.subnet.AvailabilityZone()}
}

func (s *subnetShim) Life() params.Life {
	return params.Life(s.subnet.Life().String())
}

func (s *subnetShim) Status() string {
	// TODO(dimitern): This should happen in a cleaner way.
	if s.Life() != params.Alive {
		return "terminating"
	}
	return "in-use"
}

func (s *subnetShim) SpaceName() string {
	return s.subnet.SpaceName()
}

// spaceShim forwards and adapts state.Space methods to BackingSpace.
type spaceShim struct {
	space *state.Space
}

func (s *spaceShim) Name() string {
	return s.space.Name()
}

func (s *spaceShim) ProviderId() network.Id {
	return s.space.ProviderId()
}

func (s *spaceShim) Subnets() ([]BackingSubnet, error) {
	results, err := s.space.Subnets()
	if err != nil {
		return nil, errors.Trace(err)
	}
	subnets := make([]BackingSubnet, len(results))
	for i, result := range results {
		subnets[i] = &subnetShim{subnet: result}
	}
	return subnets, nil
}

func NewStateShim(st *state.State) *stateShim {
	return &stateShim{stateenvirons.EnvironConfigGetter{st}, st}
}

// stateShim forwards and adapts state.State methods to Backing
// method.
type stateShim struct {
	stateenvirons.EnvironConfigGetter
	st *state.State
}

func (s *stateShim) AddSpace(name string, providerId network.Id, subnetIds []string, public bool) error {
	_, err := s.st.AddSpace(name, providerId, subnetIds, public)
	return err
}

func (s *stateShim) AllSpaces() ([]BackingSpace, error) {
	// TODO(dimitern): Make this ListSpaces() instead.
	results, err := s.st.AllSpaces()
	if err != nil {
		return nil, errors.Trace(err)
	}
	spaces := make([]BackingSpace, len(results))
	for i, result := range results {
		spaces[i] = &spaceShim{space: result}
	}
	return spaces, nil
}

func (s *stateShim) AddSubnet(info BackingSubnetInfo) (BackingSubnet, error) {
	// TODO(babbageclunk): we only take the first zone because
	// state.Subnet currently only stores one.
	var firstZone string
	if len(info.AvailabilityZones) > 0 {
		firstZone = info.AvailabilityZones[0]
	}
	_, err := s.st.AddSubnet(state.SubnetInfo{
		CIDR:              info.CIDR,
		VLANTag:           info.VLANTag,
		ProviderId:        info.ProviderId,
		ProviderNetworkId: info.ProviderNetworkId,
		AvailabilityZone:  firstZone,
		SpaceName:         info.SpaceName,
	})
	return nil, err // Drop the first result, as it's unused.
}

func (s *stateShim) AllSubnets() ([]BackingSubnet, error) {
	results, err := s.st.AllSubnets()
	if err != nil {
		return nil, errors.Trace(err)
	}
	subnets := make([]BackingSubnet, len(results))
	for i, result := range results {
		subnets[i] = &subnetShim{subnet: result}
	}
	return subnets, nil
}

func (s *stateShim) AvailabilityZones() ([]providercommon.AvailabilityZone, error) {
	// TODO(dimitern): Fix this to get them from state when available!
	return nil, nil
}

func (s *stateShim) SetAvailabilityZones(zones []providercommon.AvailabilityZone) error {
	return nil
}
