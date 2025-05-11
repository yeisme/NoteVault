package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileMetadataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新特定文件的元数据。这通常会创建一个新版本。
func NewUpdateFileMetadataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileMetadataLogic {
	return &UpdateFileMetadataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFileMetadataLogic) UpdateFileMetadata(req *types.UpdateFileMetadataRequest) (resp *types.UpdateFileMetadataResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
