package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Config struct {
	Envs map[string]*EnvConfig
}

type EnvConfig struct {
	URL             string
	RevisionsSchema string
}

type atlasHCL struct {
	Envs []envBlock `hcl:"env,block"`
}

type envBlock struct {
	Name      string          `hcl:"name,label"`
	URL       string          `hcl:"url,optional"`
	Migration *migrationBlock `hcl:"migration,block"`
	Remain    hcl.Body        `hcl:",remain"` // 不要な項目の吸収
}

type migrationBlock struct {
	RevisionsSchema string   `hcl:"revisions_schema,optional"`
	Remain          hcl.Body `hcl:",remain"` // 不要な項目の吸収
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = findAtlasConfig()
		if configPath == "" {
			return nil, fmt.Errorf("atlas.hcl not found in current directory")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(configPath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	var atlasConfig atlasHCL
	diags = gohcl.DecodeBody(file.Body, nil, &atlasConfig)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode HCL: %s", diags.Error())
	}

	config := &Config{
		Envs: make(map[string]*EnvConfig),
	}

	for _, env := range atlasConfig.Envs {
		envConfig := &EnvConfig{
			URL: env.URL,
		}

		if env.Migration != nil {
			envConfig.RevisionsSchema = env.Migration.RevisionsSchema
		}

		config.Envs[env.Name] = envConfig
	}

	return config, nil
}

func (c *Config) GetEnv(envName string) (*EnvConfig, error) {
	env, ok := c.Envs[envName]
	if !ok {
		available := make([]string, 0, len(c.Envs))
		for name := range c.Envs {
			available = append(available, name)
		}
		return nil, fmt.Errorf("environment '%s' not found. Available: %v", envName, available)
	}

	return env, nil
}

func findAtlasConfig() string {
	if _, err := os.Stat("atlas.hcl"); err == nil {
		return "atlas.hcl"
	}

	return ""
}
