package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mattmattox/go-smtp-relay/pkg/version"
)

// AppConfig structure for environment-based configurations.
type AppConfig struct {
	Debug           bool   `json:"debug"`
	MetricsPort     int    `json:"metricsPort"`
	SmtpPort        int    `json:"smtpPort"`
	Version         bool   `json:"version"`
	EmailServerHost string `json:"emailServerHost"`
	EmailServerPort int    `json:"emailServerPort"`
	EmailServerUser string `json:"emailServerUser"`
	EmailServerPass string `json:"emailServerPass"`
	FromAddress     string `json:"fromAddress"`
}

// CFG is the global configuration instance populated by LoadConfiguration.
var CFG AppConfig

// LoadConfiguration loads the configuration from the environment variables and command line flags.
func LoadConfiguration() {
	debug := flag.Bool("debug", parseEnvBool("DEBUG", false), "Enable debug mode")
	metricsPort := flag.Int("metricsPort", parseEnvInt("METRICS_PORT", 2112), "Port for the metrics server")
	smtpPort := flag.Int("serverPort", parseEnvInt("SMTP_PORT", 25), "Port for the server")
	showVersion := flag.Bool("version", false, "Show version and exit")
	emailServerHost := getEnvOrDefault("EMAIL_SERVER_HOST", "")
	emailServerPort := parseEnvInt("EMAIL_SERVER_PORT", 587)
	emailServerUser := getEnvOrDefault("EMAIL_SERVER_USER", "")
	emailServerPass := getEnvOrDefault("EMAIL_SERVER_PASS", "")
	fromAddress := getEnvOrDefault("FROM_ADDRESS", "")

	flag.Parse()

	CFG.Debug = *debug
	CFG.MetricsPort = *metricsPort
	CFG.SmtpPort = *smtpPort
	CFG.Version = *showVersion
	CFG.EmailServerHost = emailServerHost
	CFG.EmailServerPort = emailServerPort
	CFG.EmailServerUser = emailServerUser
	CFG.EmailServerPass = emailServerPass
	CFG.FromAddress = fromAddress

	if CFG.Version {
		fmt.Printf("Version: %s\nGit Commit: %s\nBuild Time: %s\n", version.Version, version.GitCommit, version.BuildTime)
		os.Exit(0)
	}
}

// getEnvOrDefault returns the value of the environment variable with the given key or the default value if the key is not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseEnvInt parses the environment variable with the given key and returns its integer representation or the default value if the key is not set.
func parseEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

// parseEnvBool parses the environment variable with the given key, checks case-insensitively,
// and returns its boolean representation or the default value if the key is not set or invalid.
func parseEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s not set. Using default value: %t", key, defaultValue)
		return defaultValue
	}

	// Normalize the value to lowercase and trim any whitespace.
	normalizedValue := strings.TrimSpace(strings.ToLower(value))

	// Handle explicit boolean strings ("true", "false", etc.).
	switch normalizedValue {
	case "true", "t", "1", "yes", "y":
		return true
	case "false", "f", "0", "no", "n":
		return false
	default:
		// Try to parse the value using strconv for additional safety.
		boolValue, err := strconv.ParseBool(normalizedValue)
		if err != nil {
			log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
			return defaultValue
		}
		return boolValue
	}
}
