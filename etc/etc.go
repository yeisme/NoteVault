package etc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	ProjectName = "notevault"
	configType  = "yaml"
	// DefaultConfigFile is the default name of the config file to search for.
	DefaultConfigFile = ProjectName + "service." + configType
	// DefaultSearchPaths are the default paths to search for config files
	DefaultSearchPaths = []string{".", "./etc/", "/etc/another-mentor/"}
)

// findConfigFile 在给定的搜索路径中查找指定的配置文件
func findConfigFile(filename string, searchPaths []string) (string, bool) {

	if filename == "" {
		return "", false
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		potentialPath := filepath.Join(path, filename)
		if _, err := os.Stat(potentialPath); err == nil {
			return potentialPath, true
		}
	}
	return "", false
}

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
// If configFileArg is provided, it will check if the file exists in the specified path or search for it in the predefined paths.
// It panics if the config file is not found in any of the paths.
func LoadConfig(configFileArg string) (config.Config, string) {
	c := new(config.Config)
	searchPaths := getSearchPaths()

	var filePathToLoad string
	var found bool

	if configFileArg == "" {
		// If configFileArg is empty, look for the default config file
		filePathToLoad, found = findConfigFile(DefaultConfigFile, searchPaths)

	} else {
		// First check if the provided config file exists directly
		_, err := os.Stat(configFileArg)
		if err == nil {
			filePathToLoad = configFileArg
			found = true
		} else if os.IsNotExist(err) {
			// If the file does not exist, try to find it in the search paths
			filePathToLoad, found = findConfigFile(configFileArg, searchPaths)
		}
	}
	if !found {
		panic(fmt.Sprintf("default config file %s not found in search paths: %v", DefaultConfigFile, searchPaths))
	}

	conf.MustLoad(filePathToLoad, c)
	return *c, filePathToLoad
}
