package config

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	errInvalidKV = errors.New("Invalid key/value in config file")
	errTmdbKey   = errors.New("TMDB_API_KEY Not set")
)

type Config struct {
	Ip         string
	Port       int
	TmdbApiKey string

	ConfigDir    string
	ImportTmpDir string

	ImdbDataUrl string
	ImdbTmpDir  string

	EnableImdb bool

	TmdbMaxRetries int

	DateTimeLayout string
}

func setConfigDir(c *Config) {
	conf, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config dir: %v\n", err)
	}
	c.ConfigDir = filepath.Join(conf, "riptvtime")

	c.ImportTmpDir = filepath.Join(os.TempDir(), "riptvtime_import_tmp")

	c.ImdbDataUrl = "https://datasets.imdbws.com"
	c.ImdbTmpDir = filepath.Join(os.TempDir(), "riptvtime_imdb_tmp")

	c.DateTimeLayout = "2006-01-02 15:04:05"

	c.EnableImdb = true

	err = os.MkdirAll(c.ConfigDir, 0755)
	if err != nil {
		log.Fatalf("Error creating user config dir: %v - %v\n", conf, err)
	}
	err = os.MkdirAll(c.ImportTmpDir, 0755)
	if err != nil {
		log.Fatalf("Error creating user config dir: %v - %v\n", conf, err)
	}
	err = os.MkdirAll(c.ImdbTmpDir, 0755)
	if err != nil {
		log.Fatalf("Error creating user config dir: %v - %v\n", conf, err)
	}

}

func LoadFromFile(path string) (*Config, error) {
	c := &Config{}

	setConfigDir(c)

	cfgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(cfgFile)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, err
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		kv := strings.Split(line, "=")

		if len(kv) != 2 || len(kv[0]) == 0 || len(kv[1]) == 0 {
			return nil, errInvalidKV
		}

		switch strings.ToLower(kv[0]) {
		case "ip":
			c.Ip = kv[1]
		case "port":
			c.Port, err = strconv.Atoi(kv[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid Port: %v", err.Error())
			}
		case "tmdbapikey":
			c.TmdbApiKey = kv[1]

		case "tmdbmaxretries":
			c.TmdbMaxRetries, err = strconv.Atoi(kv[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid TmdbMaxRetries: %v", err.Error())
			}
		default:
			return nil, errInvalidKV
		}
	}

	return c, nil
}

func LoadFromEnv() (*Config, error) {
	c := &Config{}

	setConfigDir(c)

	c.Ip = os.Getenv("RTT_IP")
	if len(c.Ip) == 0 {
		c.Ip = "127.0.0.1"
		slog.Warn("`RTT_IP` Not set in environment, setting default value", "ip", c.Ip)
	}

	portStr := os.Getenv("RTT_PORT")
	var err error
	c.Port, err = strconv.Atoi(portStr)
	if err != nil {
		c.Port = 5667
		slog.Warn("invalid `RTT_PORT` in environment, setting default value", "port", c.Port)
	}

	c.TmdbApiKey = os.Getenv("TMDB_API_KEY")

	if len(c.TmdbApiKey) == 0 {
		return nil, errTmdbKey
	}

	c.TmdbMaxRetries, err = strconv.Atoi(os.Getenv("TMDB_MAX_RETRIES"))
	if err != nil {
		c.TmdbMaxRetries = 10
		slog.Warn("invalid `TMDB_MAX_RETRIES` in environment, setting default value", "retries", c.TmdbMaxRetries)
	}

	return c, nil
}

func LoadFromUISetup(c *Config) (*Config, error) {
	setConfigDir(c)
	return c, nil
}
