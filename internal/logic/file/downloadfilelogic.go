package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/yeisme/notevault/pkg/storage/repository/dao"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type DownloadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

// 根据文件ID下载文件。可选下载特定版本。
func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request, w http.ResponseWriter) *DownloadFileLogic {
	return &DownloadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
		w:      w,
	}
}

// DownloadFile handles the file download request.
// 下载文件函数与其他API函数不同，它不需要返回JSON数据，而是直接将文件内容写入HTTP响应流。函数返回error只是用来表示过程中是否发生错误。
// 在当前实现中，文件内容已经通过io.Copy(w, object)写入HTTP响应，客户端会直接接收到文件内容而非JSON响应。
func (l *DownloadFileLogic) DownloadFile(req *types.FileDownloadRequest) error {
	// 初始化数据库查询
	query := dao.Use(l.svcCtx.DB)

	// 获取文件元数据
	fileRecord, err := query.File.Where(query.File.FileID.Eq(req.FileID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.Error(l.w, fmt.Errorf("file not found: %s", req.FileID))
			return fmt.Errorf("file not found: %s", req.FileID)
		}
		httpx.Error(l.w, fmt.Errorf("failed to query file: %w", err))
		return fmt.Errorf("failed to query file: %w", err)
	}

	// 确定要下载的文件版本
	var filePath string
	var contentType string

	if req.VersionID != nil {
		// 获取特定版本的文件
		versionID := fmt.Sprintf("%s_%d", req.FileID, *req.VersionID)
		fileVersion, err := query.FileVersion.Where(query.FileVersion.VersionID.Eq(versionID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				httpx.Error(l.w, fmt.Errorf("file version not found: %s, version: %d", req.FileID, *req.VersionID))
				return fmt.Errorf("file version not found: %s, version: %d", req.FileID, *req.VersionID)
			}
			httpx.Error(l.w, fmt.Errorf("failed to query file version: %w", err))
			return fmt.Errorf("failed to query file version: %w", err)
		}
		filePath = fileVersion.Path
		contentType = fileVersion.ContentType
	} else {
		// 获取最新版本的文件
		filePath = fileRecord.Path
		contentType = fileRecord.ContentType
	}

	// 从对象存储中获取文件
	object, err := l.svcCtx.OSS.GetObject(
		l.ctx,
		l.svcCtx.Config.Storage.Oss.BucketName,
		filePath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		httpx.Error(l.w, fmt.Errorf("failed to get file from storage: %w", err))
		return fmt.Errorf("failed to get file from storage: %w", err)
	}
	defer object.Close()

	// 设置响应头
	l.w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRecord.FileName))
	l.w.Header().Set("Content-Type", contentType)

	// 将文件流写入响应
	if _, err := io.Copy(l.w, object); err != nil {
		httpx.Error(l.w, fmt.Errorf("failed to write file to response: %w", err))
		return fmt.Errorf("failed to write file to response: %w", err)
	}

	return nil
}
