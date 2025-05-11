package etc

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

// Test_findConfigFile tests the findConfigFile function
func Test_findConfigFile(t *testing.T) {
	// Create a temporary directory and file for testing
	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "testconfig.yaml")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	type args struct {
		filename    string
		searchPaths []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{{
		name: "empty filename",
		args: args{
			filename:    "",
			searchPaths: []string{tmpDir},
		},
		want:  "",
		want1: false,
	},
		{
			name: "file exists in path",
			args: args{
				filename:    "testconfig.yaml",
				searchPaths: []string{tmpDir},
			},
			want:  testFilePath,
			want1: true,
		},
		{
			name: "file does not exist",
			args: args{
				filename:    "nonexistent.yaml",
				searchPaths: []string{tmpDir},
			},
			want:  "",
			want1: false,
		},
		{
			name: "empty search path",
			args: args{
				filename:    "testconfig.yaml",
				searchPaths: []string{""},
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findConfigFile(tt.args.filename, tt.args.searchPaths)
			if got != tt.want {
				t.Errorf("findConfigFile() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findConfigFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// Test_getSearchPaths tests the getSearchPaths function
func Test_getSearchPaths(t *testing.T) {
	// Save original environment and restore after test
	origEnvPath := os.Getenv("NOTEVAULT_CONFIG_PATH")
	defer os.Setenv("NOTEVAULT_CONFIG_PATH", origEnvPath)

	// Get user home directory for expected values
	userHomeDir, _ := os.UserHomeDir()
	defaultExpected := make([]string, len(DefaultSearchPaths))
	copy(defaultExpected, DefaultSearchPaths)
	if userHomeDir != "" {
		defaultExpected = append(defaultExpected, filepath.Join(userHomeDir, ".another-mentor"))
	}

	// Create enhanced expected with env var
	withEnvExpected := make([]string, len(defaultExpected))
	copy(withEnvExpected, defaultExpected)
	testEnvPath := "/custom/env/path"
	if len(withEnvExpected) > len(DefaultSearchPaths) {
		// If home directory was added, insert before it
		withEnvExpected = append(withEnvExpected[:len(DefaultSearchPaths)], append([]string{testEnvPath}, withEnvExpected[len(DefaultSearchPaths):]...)...)
	} else {
		withEnvExpected = append(withEnvExpected, testEnvPath)
	}

	tests := []struct {
		name        string
		setupEnv    func()
		want        []string
		skipHomeDir bool
	}{
		{
			name:     "default paths",
			setupEnv: func() { os.Unsetenv("NOTEVAULT_CONFIG_PATH") },
			want:     defaultExpected,
		},
		{
			name:     "with env path",
			setupEnv: func() { os.Setenv("NOTEVAULT_CONFIG_PATH", testEnvPath) },
			want:     withEnvExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			got := getSearchPaths()

			// Compare without exact order if with env variable (as long as it contains the expected values)
			if tt.name == "with env path" {
				// Check that all expected paths are in got
				for _, exp := range tt.want {
					found := slices.Contains(got, exp)
					if !found {
						t.Errorf("getSearchPaths() = %v, missing expected path %v", got, exp)
					}
				}

				// Check lengths match
				if len(got) != len(tt.want) {
					t.Errorf("getSearchPaths() returned %d paths, want %d paths", len(got), len(tt.want))
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSearchPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLoadConfig tests the LoadConfig function
func TestLoadConfig(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	defaultConfigPath := filepath.Join(tmpDir, DefaultConfigFile)
	customConfigPath := filepath.Join(tmpDir, "custom.yaml")

	// Write test content to files
	testContent := `
Name: "TestService"
Host: "localhost"
Port: 8080
Auth:
  AccessSecret: strings
  AccessExpire: 7200
`

	if err := os.WriteFile(defaultConfigPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create default config file: %v", err)
	}
	if err := os.WriteFile(customConfigPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create custom config file: %v", err)
	}

	// Save original and set up test DefaultSearchPaths
	origDefaultSearchPaths := DefaultSearchPaths
	origProjectName := ProjectName
	defer func() {
		DefaultSearchPaths = origDefaultSearchPaths
		ProjectName = origProjectName
	}()

	// Set up the test DefaultSearchPaths to include our temp directory
	DefaultSearchPaths = []string{tmpDir}

	// We can't easily test the actual config loading without mocking,
	// so this test will check if it finds the right files but not the actual config parsing.
	// For a more complete test, you'd need to create a mock for conf.MustLoad.

	type args struct {
		configFileArg string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "load default config",
			args: args{
				configFileArg: "",
			},
			wantPanic: false,
		},
		{
			name: "load explicit config path",
			args: args{
				configFileArg: customConfigPath,
			},
			wantPanic: false,
		},
		{
			name: "nonexistent config",
			args: args{
				configFileArg: "nonexistent.yaml",
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("LoadConfig() should have panicked")
					}
				}()
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("LoadConfig() should not have panicked: %v", r)
					}
				}()
			}

			// Call LoadConfig but don't verify the actual config content
			_, _ = LoadConfig(tt.args.configFileArg)
		})
	}
}
