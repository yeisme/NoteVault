package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFile string
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
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "f", "", "config 文件路径")
}

// Execute 用于执行顶层命令
func Execute(version string) {
	rootCmd.Version = version

	// 添加子命令
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(validateCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
