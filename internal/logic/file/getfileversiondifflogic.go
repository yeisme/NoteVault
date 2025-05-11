package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileVersionDiffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// (高级) 获取文件两个版本之间的差异信息 (主要针对文本文件)。
func NewGetFileVersionDiffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileVersionDiffLogic {
	return &GetFileVersionDiffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileVersionDiffLogic) GetFileVersionDiff(req *types.FileVersionDiffRequest) (resp *types.FileVersionDiffResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
