package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds user configuration from config.yaml
type Config struct {
	YourChannelID   string `yaml:"YOUR_CHANNEL_ID"`
	ChannelSize     int64  `yaml:"channel_size"`
	LogDir          string `yaml:"log_dir"`
	CredentialsFile string `yaml:"credentials_file"`
	BannedWordsFile string `yaml:"banned_words_file"` // optional, default fallback
}

// LoadConfig reads config.yaml into Config struct
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Default banned words path if not set
	if cfg.BannedWordsFile == "" {
		cfg.BannedWordsFile = "configs/banned_words.txt"
	}

	// Setup logging
	if err := setupLogging(cfg.LogDir); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setupLogging initializes log output into a file inside logDir
func setupLogging(logDir string) error {
	if logDir == "" {
		// fallback to stdout only
		log.Println("⚠️  No log_dir provided, using stdout only")
		return nil
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "tubeguardian.log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Write logs to both stdout and file
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("📂 Logging initialized: %s\n", logFile)

	return nil
}
