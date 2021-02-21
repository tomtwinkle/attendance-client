package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const configFile = "attendance.yaml"

type ConfigYaml struct {
	TenantCode string       `yaml:"tenantcode"`
	OBCiD      string       `yaml:"obcid"`
	Password   string       `yaml:"password"`
	Slack      *ConfigSlack `yaml:"slack"`
}

type ConfigSlack struct {
	Token    string               `yaml:"token"`
	Channels []ConfigSlackChannel `yaml:"channels"`
}

type ConfigSlackChannel struct {
	Name     string
	ClockIn  *ConfigSlackAction `yaml:"clockin"`
	ClockOut *ConfigSlackAction `yaml:"clockout"`
	GoOut    *ConfigSlackAction `yaml:"goout"`
	Returned *ConfigSlackAction `yaml:"returned"`
}

type ConfigSlackAction struct {
	Message string `yaml:"message"`
}

type config struct {
	ConfigPath string
}

type Config interface {
	Init() (*ConfigYaml, error)
}

func NewConfig() Config {
	return &config{ConfigPath: getConfigPath()}
}

func getConfigPath() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDirPath := filepath.Dir(execPath)
	return filepath.Join(execDirPath, configFile)
}

func (c *config) Init() (*ConfigYaml, error) {
	if _, err := os.Stat(c.ConfigPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return c.writeConfig()
	}
	return c.readConfig()
}

func (c *config) readConfig() (*ConfigYaml, error) {
	file, err := os.Open(c.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	var cfg ConfigYaml
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.TenantCode == "" {
		return nil, fmt.Errorf("tenant_code is Required [%s]", c.ConfigPath)
	}
	if cfg.OBCiD == "" {
		return nil, fmt.Errorf("obc_id is Required [%s]", c.ConfigPath)
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("password is Required [%s]", c.ConfigPath)
	}
	return &cfg, nil
}

func (c *config) writeConfig() (*ConfigYaml, error) {
	tenantCode, err := c.inputTenant()
	if err != nil {
		return nil, err
	}
	obcId, err := c.inputOBCiD()
	if err != nil {
		return nil, err
	}
	password, err := c.inputPassword()
	if err != nil {
		return nil, err
	}
	var configSlack *ConfigSlack
	if ok, err := c.inputUseSlack(); err != nil {
		return nil, err
	} else if ok {
		token, err := c.inputSlackToken()
		if err != nil {
			return nil, err
		}
		channelName, err := c.inputSlackChannel()
		if err != nil {
			return nil, err
		}
		configSlack = &ConfigSlack{
			Token: token,
			Channels: []ConfigSlackChannel{
				{
					Name:     channelName,
					ClockIn:  &ConfigSlackAction{Message: "出勤しました"},
					ClockOut: &ConfigSlackAction{Message: "退出します"},
					GoOut:    &ConfigSlackAction{Message: "外出します"},
					Returned: &ConfigSlackAction{Message: "再入しました"},
				},
			},
		}
	}

	file, err := os.Create(c.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	e := yaml.NewEncoder(file)
	defer e.Close()
	cfg := ConfigYaml{
		TenantCode: tenantCode,
		OBCiD:      obcId,
		Password:   password,
		Slack:      configSlack,
	}
	if err := e.Encode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
