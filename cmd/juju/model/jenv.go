// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/names"
	"gopkg.in/yaml.v2"

	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/environs/configstore"
	"github.com/juju/juju/juju/osenv"
)

// JenvCommand imports the given Juju generated jenv file into the local
// JUJU_HOME environments directory.
type JenvCommand struct {
	cmd.CommandBase
	jenvFile cmd.FileVar
	envName  string
}

const jenvHelpDoc = `
Copy the provided Juju generated jenv file to the proper location inside your
local JUJU_HOME. Also switch to using the resulting model as the default
one. This way it is possible to import and use jenv files generated by other
Juju commands such as "juju add-user".

Example:

  juju model jenv my-env.jenv (my-env.jenv is the path to the jenv file)
  juju model jenv my-env.jenv ec2 (copy and rename the model)
`

func (c *JenvCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "jenv",
		Args:    "<path/to/file.jenv> [<model name>]",
		Purpose: "import previously generated Juju model files",
		Doc:     strings.TrimSpace(jenvHelpDoc),
	}
}

func (c *JenvCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.New("no jenv file provided")
	}

	// Store the path to the jenv file.
	if err := c.jenvFile.Set(args[0]); err != nil {
		return errors.Annotate(err, "invalid jenv path")
	}
	args = args[1:]

	if len(args) > 0 {
		// Store and validate the provided environment name.
		c.envName, args = args[0], args[1:]
		if !names.IsValidUser(c.envName) {
			return errors.Errorf("invalid model name %q", c.envName)
		}
	} else {
		// Retrieve the environment name from the jenv file name.
		base := filepath.Base(c.jenvFile.Path)
		c.envName = base[:len(base)-len(filepath.Ext(base))]
	}

	// No other arguments are expected.
	return cmd.CheckEmpty(args)
}

func (c *JenvCommand) Run(ctx *cmd.Context) error {
	// Read data from the provided jenv file.
	data, err := c.jenvFile.Read(ctx)
	if err != nil {
		if os.IsNotExist(errors.Cause(err)) {
			return errors.NotFoundf("jenv file %q", c.jenvFile.Path)
		}
		return errors.Annotatef(err, "cannot read the provided jenv file %q", c.jenvFile.Path)
	}

	// Open the config store.
	store, err := configstore.Default()
	if err != nil {
		return errors.Annotate(err, "cannot get config store")
	}

	// Create and update the new environment info object.
	info := store.CreateInfo(c.envName)
	if err := updateEnvironmentInfo(info, data); err != nil {
		return errors.Annotatef(err, "invalid jenv file %q", c.jenvFile.Path)
	}

	// Write the environment info to JUJU_HOME.
	if err := info.Write(); err != nil {
		if errors.Cause(err) == configstore.ErrEnvironInfoAlreadyExists {
			descriptiveErr := errors.Errorf("an model named %q already exists: "+
				"you can provide a second parameter to rename the model",
				c.envName)
			return errors.Wrap(err, descriptiveErr)
		}
		return errors.Annotate(err, "cannot write the jenv file")
	}

	// Switch to the new model.
	oldModelName, err := switchEnvironment(c.envName)
	if err != nil {
		return errors.Annotatef(err, "cannot switch to the new model %q", c.envName)
	}
	if oldModelName == "" {
		fmt.Fprintf(ctx.Stdout, "-> %s\n", c.envName)
	} else {
		fmt.Fprintf(ctx.Stdout, "%s -> %s\n", oldModelName, c.envName)
	}
	return nil
}

// updateEnvironmentInfo updates the given environment info with the values
// stored in the provided YAML encoded data.
func updateEnvironmentInfo(info configstore.EnvironInfo, data []byte) error {
	var values configstore.EnvironInfoData
	if err := yaml.Unmarshal(data, &values); err != nil {
		return errors.Annotate(err, "cannot unmarshal jenv data")
	}

	// Ensure the required values are present.
	if missing := getMissingEnvironmentInfoFields(values); len(missing) != 0 {
		return errors.Errorf("missing required fields in jenv data: %s", strings.Join(missing, ", "))
	}

	// Update the environment info.
	info.SetAPICredentials(configstore.APICredentials{
		User:     values.User,
		Password: values.Password,
	})
	info.SetAPIEndpoint(configstore.APIEndpoint{
		Addresses: values.Controllers,
		Hostnames: values.ServerHostnames,
		CACert:    values.CACert,
		ModelUUID: values.ModelUUID,
	})
	info.SetBootstrapConfig(values.Config)
	return nil
}

// getMissingEnvironmentInfoFields returns a list of field names missing in the
// given environment info values. The only fields taken into consideration here
// are the ones explicitly set by the "juju add-user" command.
func getMissingEnvironmentInfoFields(values configstore.EnvironInfoData) (missing []string) {
	if values.User == "" {
		missing = append(missing, "User")
	}
	if values.Password == "" {
		missing = append(missing, "Password")
	}
	if values.ModelUUID == "" {
		missing = append(missing, "ModelUUID")
	}
	if len(values.Controllers) == 0 {
		missing = append(missing, "Controllers")
	}
	if values.CACert == "" {
		missing = append(missing, "CACert")
	}
	return missing
}

// switchEnvironment changes the default environment to the given name and
// return, if set, the current default environment name.
func switchEnvironment(envName string) (string, error) {
	if defaultEnv := os.Getenv(osenv.JujuModelEnvKey); defaultEnv != "" {
		return "", errors.Errorf("cannot switch when %s is overriding the model (set to %q)", osenv.JujuModelEnvKey, defaultEnv)
	}
	currentEnv, err := modelcmd.GetDefaultModel()
	if err != nil {
		return "", errors.Annotate(err, "cannot get the default model")
	}
	if err := modelcmd.WriteCurrentModel(envName); err != nil {
		return "", errors.Trace(err)
	}
	return currentEnv, nil
}
