package main

import (
	. "launchpad.net/gocheck"
	"launchpad.net/juju-core/charm"
	"launchpad.net/juju-core/testing"
)

type AddUnitSuite struct {
	repoSuite
}

var _ = Suite(&AddUnitSuite{})

func runAddUnit(c *C, args ...string) error {
	return testing.RunCommand(c, &AddUnitCommand{}, args)
}

func (s *AddUnitSuite) TestAddUnit(c *C) {
	testing.Charms.BundlePath(s.seriesPath, "dummy")
	err := runDeploy(c, "local:dummy", "some-service-name")
	c.Assert(err, IsNil)
	curl := charm.MustParseURL("local:precise/dummy-1")
	s.assertService(c, "some-service-name", curl, 1, 0)

	err = runAddUnit(c, "some-service-name")
	c.Assert(err, IsNil)
	s.assertService(c, "some-service-name", curl, 2, 0)

	err = runAddUnit(c, "--num-units", "2", "some-service-name")
	c.Assert(err, IsNil)
	s.assertService(c, "some-service-name", curl, 4, 0)
}
