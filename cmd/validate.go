package cmd

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/yeisme/notevault/etc"
)

var (
	withFormat bool

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Search and validate the configuration file, print the loaded file path and content(with -w flag)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// LoadConfig will search for the config file in the specified path or default paths
			// and return the configuration struct and the path to the loaded file.
			// If the config file is not found, it will panic and exit the program.
			c, filePathToLoad := etc.LoadConfig(ConfigFile)

			fmt.Printf("LoadConfig: %s\n", filePathToLoad)
			if withFormat {
				yamlData, err := yaml.Marshal(c)
				if err != nil {
					return fmt.Errorf("failed to marshal config: %w", err)
				}
				fmt.Printf("\n%s\n", yamlData)
			}

			return nil
		},
		Example: `notevault validate -f ./etc/notevaultservice.yaml`,
	}
)

func init() {
	validateCmd.Flags().BoolVarP(&withFormat, "with-format", "w", false, "格式化输出")
}
