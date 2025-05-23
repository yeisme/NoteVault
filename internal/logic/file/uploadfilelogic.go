package file

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v7"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/yeisme/notevault/pkg/storage/repository/dao"
	"github.com/yeisme/notevault/pkg/storage/repository/model"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

// Upload a new file. The actual file is sent as multipart/form-data.
func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadFileLogic {

	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

// UploadFile uploads file and metadata
func (l *UploadFileLogic) UploadFile(req *types.FileUploadRequest) (resp *types.FileUploadResponse, err error) {

	// TODO: Decoder jwt token to get userId
	// Get user ID from context
	userId, ok := l.ctx.Value("userId").(string)
	if !ok || userId == "" {
		userId = "notevault"
	}

	// frontend also check file size, but we need to check it again here
	// TODO: use multipart instead of FromFile
	file, fileHeader, err := l.r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve uploaded file: %w", err)
	}
	defer file.Close()

	// Check file size
	maxUploadSize := 16 * 1024 * 1024 // 16MB

	if fileHeader.Size > int64(maxUploadSize) {
		return nil, fmt.Errorf("File size exceeds limit (maximum %d MB)", maxUploadSize/(1024*1024))
	}

	// File name processing
	fileName := req.FileName
	if fileName == "" {
		fileName = fileHeader.Filename
	}

	// File type processing
	contentType := fileHeader.Header.Get("Content-Type")

	fileType := req.FileType
	if fileType == "" {
		fileHeaderBytes := make([]byte, 261)
		if _, err := file.Read(fileHeaderBytes); err != nil {
			return nil, fmt.Errorf("Failed to read file header: %w", err)
		}
		// 重置文件指针以便后续操作
		if _, err := file.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("Failed to reset file pointer: %w", err)
		}
		kind, err := filetype.Match(fileHeaderBytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to determine file type: %w", err)
		}
		if kind == filetype.Unknown {
			// 如果无法识别文件类型，检查文件扩展名
			extension := strings.ToLower(filepath.Ext(fileName))

			// 为常见文本文件类型设置正确的 MIME 类型
			textExtensions := map[string]string{
				".md":   "text/markdown",
				".txt":  "text/plain",
				".csv":  "text/csv",
				".json": "application/json",
				".xml":  "application/xml",
				".html": "text/html",
				".css":  "text/css",
				".js":   "application/javascript",
				".yml":  "application/x-yaml",
				".yaml": "application/x-yaml",
				".toml": "application/toml",
				".ini":  "text/plain",
				".conf": "text/plain",
				".log":  "text/plain",
				".sql":  "application/sql",
			}

			if mimeType, ok := textExtensions[extension]; ok {
				fileType = mimeType
			} else if contentType != "" {
				fileType = contentType
			} else {
				fileType = "application/octet-stream"
			}
		} else {
			fileType = kind.MIME.Value
		}
	}

	// Calculate file hash (sha256) before uploading
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("Failed to calculate file hash: %w", err)
	}
	fileID := fmt.Sprintf("%x", hash.Sum(nil))

	// Reset file pointer to beginning for upload
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("Failed to reset file for upload: %w", err)
	}

	// Get the timestamp and yearMonth format
	now := time.Now()
	yearMonth := now.Format("200601")
	now_time := now.Unix()

	// Build final storage path with fileID
	storePath := fmt.Sprintf("%s/%s/%s", userId, yearMonth, fileID)

	// Upload file directly to final location
	_, err = l.svcCtx.OSS.PutObject(
		l.ctx,
		l.svcCtx.Config.Storage.Oss.BucketName,
		storePath,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		l.Error("Failed to upload file to OSS", logx.Field("error", err))
		return nil, fmt.Errorf("Failed to upload file to storage: %w", err)
	}

	// Initialize the query using gorm gen
	query := dao.Use(l.svcCtx.DB)

	// Check if file with same hash already exists in database using gorm gen
	existingFile, err := query.File.Where(query.File.FileID.Eq(fileID)).First()
	if err == nil && existingFile != nil {
		// File with same hash already exists, clean up the uploaded file
		l.cleanupFile(storePath)
		return nil, fmt.Errorf("File with identical content already exists (ID: %s)", fileID)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		l.Error("Failed to check for existing file", logx.Field("error", err))
		l.cleanupFile(storePath)
		return nil, fmt.Errorf("Failed to check for duplicate file: %w", err)
	}

	// Process file tags, splitting by comma
	// TODO: I hope tags auto generate, use some LLM or NLP to generate tags
	var tags []string
	if req.Tags != "" {
		// Split and clean tags
		rawTags := strings.Split(req.Tags, ",")
		for _, tag := range rawTags {
			trimmedTag := strings.TrimSpace(tag)
			if trimmedTag != "" {
				tags = append(tags, trimmedTag)
			}
		}
	}

	// Create metadata
	metadata := &types.FileMetadata{
		FileID:      fileID,
		UserID:      userId,
		FileName:    fileName,
		FileType:    fileType,
		ContentType: contentType,
		Size:        fileHeader.Size,
		Path:        storePath,
		CreatedAt:   now_time,
		UpdatedAt:   now_time,
		Version:     1,
		Tags:        tags,
		Description: req.Description,
	}

	// Save metadata to database using the model directly
	fileModel := model.File{
		FileID:         metadata.FileID,
		UserID:         metadata.UserID,
		FileName:       metadata.FileName,
		FileType:       metadata.FileType,
		ContentType:    metadata.ContentType,
		Size:           metadata.Size,
		Path:           metadata.Path,
		CreatedAt:      metadata.CreatedAt,
		UpdatedAt:      metadata.UpdatedAt,
		CurrentVersion: int32(metadata.Version),
		Description:    metadata.Description,
	}

	// Save file_version to database
	fileVersion := model.FileVersion{
		VersionID:     fmt.Sprintf("%s_%d", fileID, fileModel.CurrentVersion),
		FileID:        fileID,
		VersionNumber: fileModel.CurrentVersion,
		Size:          fileHeader.Size,
		Path:          storePath,
		ContentType:   contentType,
		CreatedAt:     now_time,
		CommitMessage: "Initial upload",
	}

	// Use gorm gen to create the file
	if err := query.File.Create(&fileModel); err != nil {
		l.Error("Failed to save file metadata", logx.Field("error", err))
		// File uploaded but metadata save failed, should delete file from OSS
		l.cleanupFile(storePath)
		return nil, fmt.Errorf("Failed to save file information: %w", err)
	}

	if err := query.FileVersion.Create(&fileVersion); err != nil {
		// 数据库可能已经存在相同文件
		l.Error("Failed to save file version", logx.Field("error", err))
		// 如果版本保存失败，尝试删除之前创建的文件记录
		_, deleteErr := query.File.Where(query.File.FileID.Eq(fileID)).Delete(&fileModel)
		if deleteErr != nil {
			l.Error("Failed to clean up file record after version creation failed",
				logx.Field("fileID", fileID), logx.Field("error", deleteErr))
			// 继续处理，不要因为清理失败而中断
		}
		l.cleanupFile(storePath)
		return nil, fmt.Errorf("Failed to save file version information: %w", err)
	}

	// Save file tags to database
	if len(tags) > 0 {
		// 使用事务来保证所有的标签操作一致性
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			// 使用事务的查询构建器
			txQuery := dao.Use(tx)

			// 为每个标签创建记录（如果不存在）并关联到文件
			for _, tagName := range tags {
				// 查找或创建标签
				var tagModel *model.Tag
				tagModel, err := txQuery.Tag.Where(txQuery.Tag.Name.Eq(tagName)).First()

				if err != nil {
					// 标签不存在，创建新标签
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tagID := fmt.Sprintf("tag_%x", sha256.Sum256([]byte(tagName)))[:36]
						tagModel = &model.Tag{
							TagID: tagID,
							Name:  tagName,
						}
						if err := txQuery.Tag.Create(tagModel); err != nil {
							return fmt.Errorf("failed to create tag: %w", err)
						}
					} else {
						return fmt.Errorf("failed to query tag: %w", err)
					}
				}

				// 创建文件和标签的关联
				fileTagRelation := model.FileTag{
					FileID: fileID,
					TagID:  tagModel.TagID,
				}

				if err := txQuery.FileTag.Create(&fileTagRelation); err != nil {
					return fmt.Errorf("failed to create file-tag relation: %w", err)
				}
			}

			return nil
		})

		if err != nil {
			l.Error("Failed to save file tags", logx.Field("error", err))
			// 记录错误但不影响文件上传的整体结果
			// 注意：这里我们选择继续而不是失败整个上传过程
			l.Logger.Error("Failed to save tags but file upload is successful", logx.Field("error", err))
		}
	}

	// Create response
	resp = &types.FileUploadResponse{
		FileID:      fileID,
		FileName:    fileName,
		ContentType: contentType,
		Size:        fileHeader.Size,
		Message:     "File upload successful",
		Version:     1,
	}

	logx.Infof("File uploaded successfully: %s", fileName)

	return resp, nil
}

// Clean up the uploaded file (called when an error occurs)
func (l *UploadFileLogic) cleanupFile(path string) {
	err := l.svcCtx.OSS.RemoveObject(
		l.ctx,
		l.svcCtx.Config.Storage.Oss.BucketName,
		path,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		l.Error("Failed to clean up OSS file", logx.Field("path", path), logx.Field("error", err))
	}
}
