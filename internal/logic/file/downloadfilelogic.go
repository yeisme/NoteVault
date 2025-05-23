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
// This function differs from other API functions as it writes the file content directly to the HTTP response stream
// instead of returning JSON data. The error return value only indicates if an error occurred during the process.
func (l *DownloadFileLogic) DownloadFile(req *types.FileDownloadRequest) error {

	// Initialize the query using gorm gen
	query := dao.Use(l.svcCtx.DB)

	fileQueryBuilder := query.File.WithContext(l.ctx).Where(query.File.DeletedAt.Eq(0))
	fileVersionQueryBuilder := query.FileVersion.WithContext(l.ctx).Where(query.FileVersion.DeletedAt.Eq(0))

	// Get file metadata with DeletedAt.IsNull() condition
	fileRecord, err := fileQueryBuilder.Where(
		query.File.FileID.Eq(req.FileID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.Error(l.w, fmt.Errorf("file not found: %s", req.FileID))
			return fmt.Errorf("file not found: %s", req.FileID)
		}
		httpx.Error(l.w, fmt.Errorf("failed to query file: %w", err))
		return fmt.Errorf("failed to query file: %w", err)
	}

	// Determine which file version to download
	var filePath string
	var contentType string

	if req.VersionNumber != nil && *req.VersionNumber > 0 {
		// Get specific file version using Gen API with DeletedAt.IsNull() condition
		versionNumber := int32(*req.VersionNumber)
		fileVersion, err := fileVersionQueryBuilder.Where(
			query.FileVersion.FileID.Eq(req.FileID),
			query.FileVersion.VersionNumber.Eq(versionNumber),
		).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				httpx.Error(l.w, fmt.Errorf("file version not found: %s, version: %d", req.FileID, *req.VersionNumber))
				return fmt.Errorf("file version not found: %s, version: %d", req.FileID, *req.VersionNumber)
			}
			httpx.Error(l.w, fmt.Errorf("failed to query file version: %w", err))
			return fmt.Errorf("failed to query file version: %w", err)
		}
		filePath = fileVersion.Path
		contentType = fileVersion.ContentType
	} else {
		// Use latest version of the file
		filePath = fileRecord.Path
		contentType = fileRecord.ContentType
	}

	// Get file from object storage
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

	// Set response headers
	l.w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileRecord.FileName))
	l.w.Header().Set("Content-Type", contentType)

	// Stream file to response
	if _, err := io.Copy(l.w, object); err != nil {
		httpx.Error(l.w, fmt.Errorf("failed to write file to response: %w", err))
		return fmt.Errorf("failed to write file to response: %w", err)
	}

	return nil
}
