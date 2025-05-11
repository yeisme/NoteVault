package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileVersionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件的版本历史。
func NewGetFileVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileVersionsLogic {
	return &GetFileVersionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileVersionsLogic) GetFileVersions(req *types.GetFileVersionsRequest) (resp *types.GetFileVersionsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
