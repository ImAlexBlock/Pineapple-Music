package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     int    `mapstructure:"port"`
	DataDir  string `mapstructure:"data_dir"`
	LogLevel string `mapstructure:"log_level"`

	// Turnstile (optional)
	TurnstileSiteKey string `mapstructure:"turnstile_site_key"`
	TurnstileSecret  string `mapstructure:"turnstile_secret"`

	// Rate limiting
	RateLimitRPS   float64 `mapstructure:"rate_limit_rps"`
	RateLimitBurst int     `mapstructure:"rate_limit_burst"`

	// Session
	SessionMaxAge int  `mapstructure:"session_max_age"` // seconds
	SecureCookie  bool `mapstructure:"secure_cookie"`    // set Secure flag on cookies (enable behind HTTPS)

	// Upload
	MaxUploadSize int64 `mapstructure:"max_upload_size"` // bytes

	// Trusted proxies (CIDR list, e.g. "127.0.0.1/8,::1/128")
	TrustedProxies string `mapstructure:"trusted_proxies"`
}

func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "pineapple.db")
}

func (c *Config) MusicDir() string {
	return filepath.Join(c.DataDir, "music")
}

func (c *Config) TrustedProxiesList() []string {
	if c.TrustedProxies == "" {
		return nil
	}
	parts := strings.Split(c.TrustedProxies, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("port", 3880)
	v.SetDefault("data_dir", "./data")
	v.SetDefault("log_level", "info")
	v.SetDefault("rate_limit_rps", 10.0)
	v.SetDefault("rate_limit_burst", 20)
	v.SetDefault("session_max_age", 86400) // 24h
	v.SetDefault("secure_cookie", false)
	v.SetDefault("max_upload_size", 52428800) // 50MB
	v.SetDefault("trusted_proxies", "")

	// Environment variables with PM_ prefix
	v.SetEnvPrefix("PM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Optional config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./data")
	_ = v.ReadInConfig() // ignore if not found

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
