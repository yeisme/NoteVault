package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yeisme/notevault/etc"
	"github.com/yeisme/notevault/internal/handler"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/pkg/mq"
	"github.com/yeisme/notevault/pkg/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var (
	dryrun *bool

	ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 解析配置文件
			c, filePathToLoad := etc.LoadConfig(ConfigFilePath)

			logx.MustSetup(c.Log)

			// 现在可以安全地使用日志系统
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
			if err := mq.InitMQ(c.MQ); err != nil {
				if *dryrun {
					logx.Errorf("MQ connection check failed: %v", err)
					logx.Info("Running in dry run mode, continuing despite MQ errors")
				} else {
					return err
				}
			}

			// 测试ServiceContext初始化
			logx.Info("Testing ServiceContext initialization...")
			ctx := svc.NewServiceContext(c)
			if ctx == nil {
				if *dryrun {
					logx.Error("ServiceContext initialization failed")
					logx.Info("Running in dry run mode, continuing despite ServiceContext errors")
				} else {
					logx.Error("Failed to create ServiceContext")
					return nil
				}
			} else {
				logx.Debugf("MQ Client Type: %s, Available: %v", ctx.MQ.Client.Type(), ctx.MQ.Client.IsAvailable())
			}
			// 如果是干启动模式，输出配置检查成功并退出
			if *dryrun {
				logx.Info("Configuration check completed, Dry run completed, exiting...")
				return nil
			}

			server := rest.MustNewServer(c.RestConf)
			defer server.Stop()

			// 添加全局中间件

			handler.RegisterHandlers(server, ctx)
			logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
			server.Start()
			return nil
		},
		Example: `notevault server -f ./etc/notevaultservice.yaml`,
	}
)

func init() {
	dryrun = ServerCmd.Flags().BoolP("dryrun", "d", false, "dryrun mode, only check config and db connection")
}
