package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevertFileVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将文件恢复到特定版本。
func NewRevertFileVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevertFileVersionLogic {
	return &RevertFileVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevertFileVersionLogic) RevertFileVersion(req *types.RevertFileVersionRequest) (resp *types.RevertFileVersionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
