package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
		userId = "test_user" // 临时测试用户ID
		l.Logger.Info("使用测试用户ID", logx.Field("userId", userId))
	}

	// frontend also check file size, but we need to check it again here
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

	// Create a temporary buffer to store file for hash calculation
	// Note: For large files, consider using a more memory-efficient approach
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %w", err)
	}

	// Calculate file hash (SHA-256)
	hasher := sha256.New()
	hasher.Write(fileBytes)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Use file hash as fileID
	fileID := fileHash

	// Get the timestamp and yearMonth format earlier in the code
	now := time.Now()
	yearMonth := now.Format("200601")
	now_time := now.Unix()

	// Initialize the query using gorm gen
	fileQuery := dao.Use(l.svcCtx.DB)

	// Check if file with same hash already exists in database using gorm gen
	existingFile, err := fileQuery.File.Where(fileQuery.File.FileID.Eq(fileID)).First()
	if err == nil && existingFile != nil {
		// File with same hash already exists
		return nil, fmt.Errorf("File with identical content already exists (ID: %s)", fileID)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		l.Error("Failed to check for existing file", logx.Field("error", err))
		return nil, fmt.Errorf("Failed to check for duplicate file: %w", err)
	}

	// Build storage path: userId/yearMonth/fileHash
	storePath := fmt.Sprintf("%s/%s/%s", userId, yearMonth, fileID)

	// Optional: Check if file exists in OSS storage
	// This is useful if files might be directly uploaded to OSS without database entries
	_, err = l.svcCtx.OSS.StatObject(
		l.ctx,
		l.svcCtx.Config.Storage.Oss.BucketName,
		storePath,
		minio.StatObjectOptions{},
	)
	if err == nil {
		// File already exists in OSS
		return nil, fmt.Errorf("File with identical content already exists in storage")
	} else if !isMinioErrorNotFound(err) {
		// Error other than "not found"
		l.Error("Failed to check OSS for existing file", logx.Field("error", err))
		return nil, fmt.Errorf("Failed to check storage for duplicate file: %w", err)
	}

	// Reset file reader for upload
	fileReader := io.NopCloser(strings.NewReader(string(fileBytes)))

	// Upload file to OSS
	_, err = l.svcCtx.OSS.PutObject(
		l.ctx,
		l.svcCtx.Config.Storage.Oss.BucketName,
		storePath,
		fileReader,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		l.Error("Failed to upload file to OSS", logx.Field("error", err))
		return nil, fmt.Errorf("Failed to upload file to storage: %w", err)
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
		VersionID:     fmt.Sprintf("%d", fileModel.CurrentVersion),
		FileID:        fileID,
		VersionNumber: 1,
		Size:          fileHeader.Size,
		Path:          storePath,
		ContentType:   contentType,
		CreatedAt:     now_time,
		CommitMessage: "Initial upload",
	}

	// Use gorm gen to create the file
	if err := fileQuery.File.Create(&fileModel); err != nil {
		l.Error("Failed to save file metadata", logx.Field("error", err))
		// File uploaded but metadata save failed, should delete file from OSS
		l.cleanupFile(storePath)
		return nil, fmt.Errorf("Failed to save file information: %w", err)
	}

	if err := fileQuery.FileVersion.Create(&fileVersion); err != nil {
		l.Error("Failed to save file version", logx.Field("error", err))
		// 如果版本保存失败，尝试删除之前创建的文件记录
		_, deleteErr := fileQuery.File.Where(fileQuery.File.FileID.Eq(fileID)).Delete(&fileModel)
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

// Helper function to check if the error is a "not found" error from Minio
func isMinioErrorNotFound(err error) bool {
	errResp, ok := err.(minio.ErrorResponse)
	return ok && (errResp.Code == "NoSuchKey" || errResp.Code == "NotFound")
}
