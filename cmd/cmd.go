package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFilePath string
	rootCmd    = &cobra.Command{
		Use:   "notevault",
		Short: "notevault cli",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	// 设置全局标志
	rootCmd.PersistentFlags().StringVarP(&ConfigFilePath, "config", "f", "", "config file path")
}

// Execute 用于执行顶层命令
func Execute(version string, environment string, args ...string) {
	rootCmd.Version = version

	// 添加子命令
	rootCmd.AddCommand(
		serverCmd,
		validateCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
