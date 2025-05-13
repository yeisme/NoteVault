package etc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	ProjectName = "notevault"
	configType  = "yaml"
	// DefaultConfigFile is the default name of the config file to search for.
	DefaultConfigFile = ProjectName + "service." + configType
	// DefaultSearchPaths are the default paths to search for config files
	DefaultSearchPaths = []string{".", "./etc/", "/etc/another-mentor/", "/app/etc/"}
)

// getSearchPaths returns a slice of paths to search for configuration files.
func getSearchPaths() []string {
	paths := make([]string, len(DefaultSearchPaths))
	copy(paths, DefaultSearchPaths)

	// Add a custom path if specified in the environment variable
	if envPath := os.Getenv("NOTEVAULT_CONFIG_PATH"); envPath != "" {
		paths = append(paths, envPath)
	}

	// Add a user-specific path if the home directory can be determined
	userHomeDir, err := os.UserHomeDir()
	if err == nil && userHomeDir != "" {
		paths = append(paths,
			filepath.Join(userHomeDir, ".another-mentor"),
		)
	}

	return paths
}

// LoadConfig loads the configuration from the specified file or searches for it in predefined paths.
// If configFileArg is empty, it will look for the default config file in the search paths.
// If configFileArg is provided, it will be used directly as the config file path.
// It panics if the config file is not found or cannot be parsed.
func LoadConfig(configFilePath string) (config.Config, string) {
	c := new(config.Config)
	v := viper.New()

	v.SetConfigType(configType)

	if configFilePath != "" {
		v.SetConfigFile(configFilePath)
	} else {
		v.SetConfigName(ProjectName + "service") // Name of config file (without extension)

		searchPaths := getSearchPaths()
		for _, path := range searchPaths {
			v.AddConfigPath(path) // Add path to search paths
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if configFilePath == "" {
			// Error occurred while searching for the default config file
			panic(fmt.Sprintf("default config file %s (name: %s, type: %s) not found in search paths %v. Error: %s",
				DefaultConfigFile, ProjectName+"service", configType, getSearchPaths(), err))
		} else {
			// Error occurred while loading the specified config file
			panic(fmt.Sprintf("error loading config file %s: %s", configFilePath, err))
		}
	}

	// 使用 gozero 的 conf 包来解析配置
	conf.MustLoad(v.ConfigFileUsed(), c)

	return *c, v.ConfigFileUsed()
}
