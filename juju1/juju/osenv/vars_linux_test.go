// Copyright 2014 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv_test

import (
	"path/filepath"

	gc "gopkg.in/check.v1"

	"github.com/juju/1.25-upgrade/juju1/juju/osenv"
)

func (s *varsSuite) TestJujuHome(c *gc.C) {
	path := `/foo/bar/baz/`
	s.PatchEnvironment("HOME", path)
	c.Assert(osenv.JujuHomeLinux(), gc.Equals, filepath.Join(path, ".juju"))
}
