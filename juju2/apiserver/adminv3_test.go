// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver_test

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v2"

	"github.com/juju/1.25-upgrade/juju2/api"
	"github.com/juju/1.25-upgrade/juju2/rpc"
	"github.com/juju/1.25-upgrade/juju2/testing/factory"
)

type loginV3Suite struct {
	loginSuite
}

var _ = gc.Suite(&loginV3Suite{})

func (s *loginV3Suite) TestClientLoginToEnvironment(c *gc.C) {
	info := s.APIInfo(c)
	apiState, err := api.Open(info, api.DialOpts{})
	c.Assert(err, jc.ErrorIsNil)
	defer apiState.Close()

	client := apiState.Client()
	_, err = client.GetModelConstraints()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *loginV3Suite) TestClientLoginToServer(c *gc.C) {
	info := s.APIInfo(c)
	info.ModelTag = names.ModelTag{}
	apiState, err := api.Open(info, api.DialOpts{})
	c.Assert(err, jc.ErrorIsNil)
	defer apiState.Close()

	client := apiState.Client()
	_, err = client.GetModelConstraints()
	c.Assert(errors.Cause(err), gc.DeepEquals, &rpc.RequestError{
		Message: `facade "Client" not supported for controller API connection`,
		Code:    "not supported",
	})
}

func (s *loginV3Suite) TestClientLoginToServerNoAccessToControllerEnv(c *gc.C) {
	password := "shhh..."
	user := s.Factory.MakeUser(c, &factory.UserParams{
		NoModelUser: true,
		Password:    password,
	})

	info := s.APIInfo(c)
	info.Tag = user.Tag()
	info.Password = password
	info.ModelTag = names.ModelTag{}
	apiState, err := api.Open(info, api.DialOpts{})
	c.Assert(err, jc.ErrorIsNil)
	defer apiState.Close()
	// The user now has last login updated.
	err = user.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	lastLogin, err := user.LastLogin()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(lastLogin, gc.NotNil)
}

func (s *loginV3Suite) TestClientLoginToRootOldClient(c *gc.C) {
	info := s.APIInfo(c)
	info.Tag = nil
	info.Password = ""
	info.ModelTag = names.ModelTag{}
	info.SkipLogin = true
	apiState, err := api.Open(info, api.DialOpts{})
	c.Assert(err, jc.ErrorIsNil)

	err = apiState.APICall("Admin", 2, "", "Login", struct{}{}, nil)
	c.Assert(err, gc.ErrorMatches, ".*this version of Juju does not support login from old clients.*")
}
