// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// +build go1.3

package lxd_test

import (
	gitjujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/1.25-upgrade/juju2/provider/lxd"
)

type environNetSuite struct {
	lxd.BaseSuite
}

var _ = gc.Suite(&environNetSuite{})

func (s *environNetSuite) TestGlobalFirewallName(c *gc.C) {
	uuid := s.Config.UUID()
	fwname := lxd.GlobalFirewallName(s.Env)

	c.Check(fwname, gc.Equals, "juju-"+uuid)
}

func (s *environNetSuite) TestOpenPortsOkay(c *gc.C) {
	err := s.Env.OpenPorts(s.Rules)

	c.Check(err, jc.ErrorIsNil)
}

func (s *environNetSuite) TestOpenPortsAPI(c *gc.C) {
	fwname := lxd.GlobalFirewallName(s.Env)
	err := s.Env.OpenPorts(s.Rules)
	c.Assert(err, jc.ErrorIsNil)

	s.Stub.CheckCalls(c, []gitjujutesting.StubCall{{
		FuncName: "OpenPorts",
		Args: []interface{}{
			fwname,
			s.Rules,
		},
	}})
}

func (s *environNetSuite) TestClosePortsOkay(c *gc.C) {
	err := s.Env.ClosePorts(s.Rules)

	c.Check(err, jc.ErrorIsNil)
}

func (s *environNetSuite) TestClosePortsAPI(c *gc.C) {
	fwname := lxd.GlobalFirewallName(s.Env)
	err := s.Env.ClosePorts(s.Rules)
	c.Assert(err, jc.ErrorIsNil)

	s.Stub.CheckCalls(c, []gitjujutesting.StubCall{{
		FuncName: "ClosePorts",
		Args: []interface{}{
			fwname,
			s.Rules,
		},
	}})
}

func (s *environNetSuite) TestPortsOkay(c *gc.C) {
	s.Firewaller.PortRanges = s.Rules

	ports, err := s.Env.IngressRules()
	c.Assert(err, jc.ErrorIsNil)

	c.Check(ports, jc.DeepEquals, s.Rules)
}

func (s *environNetSuite) TestPortsAPI(c *gc.C) {
	fwname := lxd.GlobalFirewallName(s.Env)
	_, err := s.Env.IngressRules()
	c.Assert(err, jc.ErrorIsNil)

	s.Stub.CheckCalls(c, []gitjujutesting.StubCall{{
		FuncName: "Ports",
		Args: []interface{}{
			fwname,
		},
	}})
}
