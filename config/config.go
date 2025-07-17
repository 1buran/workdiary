package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigPath = ".config/workdiary/config.json"

type EmployeeProfile struct {
	DailyHoursLimit, HourlyRate float32
}

type RedmineClientConf struct {
	Name, Url, Token, UserId string
	Disabled, LogEnabled     bool
	EmployeeProfile          EmployeeProfile `json:"employee"`
}

type ApiClientConf struct {
	Redmine []RedmineClientConf
}

type Infrastructure struct {
	ApiClient ApiClientConf
}

type Theme struct {
	Colors map[string]string
}

type Application struct {
	Theme Theme
}

type Config struct {
	Infra Infrastructure
	App   Application
}

// Color of item, write to stderr a message if requested color is not found in theme.
func (c *Config) Color(s string) string {
	if color, ok := c.App.Theme.Colors[s]; ok {
		return color
	}
	return "11" // fallback color
}

func DefaultConfig() Config {
	return Config{
		App: Application{
			Theme: Theme{
				Colors: map[string]string{
					"dayOff":         "#a958ad",
					"workingDay":     "#0d420d",
					"expectedAmount": "#ff9ff3",
					"infactAmount":   "#4cd137",
					"summary":        "#fd79a8",
				},
			},
		},
		Infra: Infrastructure{
			ApiClient: ApiClientConf{
				Redmine: []RedmineClientConf{
					{
						Name:            "example",
						Url:             "http://example.com",
						UserId:          "100",
						Token:           "xxxxxxxxxxxxxx",
						LogEnabled:      true,
						EmployeeProfile: EmployeeProfile{DailyHoursLimit: 8, HourlyRate: 10.5},
					},
				},
			},
		},
	}
}

func PrintDefaultConfig() {
	cfg := DefaultConfig()
	b, _ := json.Marshal(cfg)

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	fmt.Println(out.String())
}

func ReadConfig(confPath string) (*Config, error) {
	cfg := DefaultConfig()

	userHome, _ := os.UserHomeDir() // ignore rare cases when user home is undefined
	userConfig := filepath.Join(userHome, ConfigPath)

	if len(confPath) > 0 {
		userConfig = confPath
	}

	b, err := os.ReadFile(userConfig)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
