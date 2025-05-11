package file

import (
	"context"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除文件。
func NewBatchDeleteFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFilesLogic {
	return &BatchDeleteFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteFilesLogic) BatchDeleteFiles(req *types.BatchDeleteFilesRequest) (resp *types.BatchDeleteFilesResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
