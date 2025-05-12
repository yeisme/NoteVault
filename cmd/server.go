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
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 解析配置文件
			c, filePathToLoad := etc.LoadConfig(ConfigFile)

			// 初始化数据库连接
			if err := storage.InitStorage(c); err != nil {
				return err
			}

			server := rest.MustNewServer(c.RestConf)
			defer server.Stop()

			ctx := svc.NewServiceContext(c)
			handler.RegisterHandlers(server, ctx)
			logx.Infof("LoadConfig: %s Starting server at %s:%d...\n", filePathToLoad, c.Host, c.Port)
			server.Start()
			return nil
		},
		Example: `notevault server -f ./etc/notevaultservice.yaml`,
	}
)
