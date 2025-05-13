package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yeisme/notevault/etc"
	"github.com/yeisme/notevault/internal/handler"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/pkg/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var (
	dryrun *bool

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 解析配置文件
			c, filePathToLoad := etc.LoadConfig(ConfigFile)

			logx.Infof("LoadConfig: %s", filePathToLoad)

			// 初始化数据库连接
			if err := storage.InitStorage(c.Storage, c.Log); err != nil {
				if *dryrun {
					logx.Errorf("Database connection check failed: %v", err)
					logx.Info("Running in dry run mode, continuing despite database errors")
				} else {
					return err
				}
			}

			// 如果是干启动模式，输出配置检查成功并退出
			if *dryrun {
				logx.Info("Configuration check completed, Dry run completed, exiting...")
				return nil
			}

			server := rest.MustNewServer(c.RestConf)
			defer server.Stop()

			ctx := svc.NewServiceContext(c)
			handler.RegisterHandlers(server, ctx)
			logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
			server.Start()
			return nil
		},
		Example: `notevault server -f ./etc/notevaultservice.yaml`,
	}
)

func init() {
	dryrun = serverCmd.Flags().BoolP("dryrun", "d", false, "dryrun mode, only check config and db connection")
}
