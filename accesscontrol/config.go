package accesscontrol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type accessControlConfig struct {
	Functions map[string]functionConfig `json:"functions"`
	Policies  map[string]policyConfig   `json:"policies"`
}

type functionConfig struct {
	Permissions             []permission `json:"permissions"`
	AbsoluteDependencies    []string     `json:"absoluteDependencies"`
	ConditionalDependencies []string     `json:"conditionalDependencies"`
}

type policyConfig struct {
	Dependencies []string     `json:"dependencies"`
	Permissions  []permission `json:"permissions"`
}

type permission struct {
	DataType  string `json:"dataType"`
	Operation string `json:"operation"`
}

func getConfig() (accessControlConfig, error) {
	configPath, ok := os.LookupEnv("ACCESS_CONTROL_CONFIG")
	if !ok {
		configPath = "config.json"
	}

	return parseConfig(configPath)
}

func parseConfig(path string) (accessControlConfig, error) {
	config := accessControlConfig{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func (c *accessControlConfig) getImmediatePermissionsForFunction(function string) ([]permission, error) {
	functionConfig, ok := c.Functions[function]
	if !ok {
		return nil, fmt.Errorf("invalid function: %s", function)
	}

	return functionConfig.Permissions, nil
}

func (c *accessControlConfig) getAbsolutePermissionsForFunction(function string) ([]permission, error) {
	absolutePermissionsSet := map[permission]bool{}
	err := c.getAbsolutePermissionsSetForFunction(function, absolutePermissionsSet)
	if err != nil {
		return nil, err
	}

	absolutePermissions := []permission{}
	for absolutePermission := range absolutePermissionsSet {
		absolutePermissions = append(absolutePermissions, absolutePermission)
	}

	return absolutePermissions, nil
}

func (c *accessControlConfig) getAbsolutePermissionsSetForFunction(function string, absolutePermissions map[permission]bool) error {
	functionConfig, ok := c.Functions[function]
	if !ok {
		return fmt.Errorf("invalid function: %s", function)
	}

	for _, absolutePermission := range functionConfig.Permissions {
		absolutePermissions[absolutePermission] = true
	}

	for _, absoluteDependency := range functionConfig.AbsoluteDependencies {
		err := c.getAbsolutePermissionsSetForFunction(absoluteDependency, absolutePermissions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *accessControlConfig) getPermissionsForPolicy(policy string) ([]permission, error) {
	policyPermissionsSet := map[permission]bool{}
	err := c.getPermisionsSetForPolicy(policy, policyPermissionsSet)
	if err != nil {
		return nil, err
	}

	policyPermissions := []permission{}
	for policyPermission := range policyPermissionsSet {
		policyPermissions = append(policyPermissions, policyPermission)
	}

	return policyPermissions, nil
}

func (c *accessControlConfig) getPermisionsSetForPolicy(policy string, policyPermissions map[permission]bool) error {
	policyConfig, ok := c.Policies[policy]
	if !ok {
		return fmt.Errorf("invalid policy: %s", policy)
	}

	for _, policyPermission := range policyConfig.Permissions {
		policyPermissions[policyPermission] = true
	}

	for _, policyDependency := range policyConfig.Dependencies {
		err := c.getPermisionsSetForPolicy(policyDependency, policyPermissions)
		if err != nil {
			return err
		}
	}

	return nil
}
