package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileMetadataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取特定文件的元数据。可选获取特定版本的元数据。
func NewGetFileMetadataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileMetadataLogic {
	return &GetFileMetadataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileMetadataLogic) GetFileMetadata(req *types.GetFileMetadataRequest) (resp *types.GetFileMetadataResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
