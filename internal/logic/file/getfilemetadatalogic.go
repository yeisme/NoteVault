package file

import (
	"context"
	"errors"
	"fmt"

	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/yeisme/notevault/pkg/storage/repository/dao"
	"github.com/yeisme/notevault/pkg/storage/repository/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetFileMetadataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get file metadata for a specific file. Optionally get metadata for a specific version.
func NewGetFileMetadataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileMetadataLogic {
	return &GetFileMetadataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileMetadataLogic) GetFileMetadata(req *types.GetFileMetadataRequest) (resp *types.GetFileMetadataResponse, err error) {

	// Initialize response
	resp = &types.GetFileMetadataResponse{}

	// Initialize the query using gorm gen
	query := dao.Use(l.svcCtx.DB)

	fileQueryBuilder := query.File.WithContext(l.ctx).Where(query.File.DeletedAt.Eq(0))
	fileVersionQueryBuilder := query.FileVersion.WithContext(l.ctx).Where(query.FileVersion.DeletedAt.Eq(0))

	// 1. Query file basic information by FileID
	file, err := fileQueryBuilder.Where(query.File.FileID.Eq(req.FileID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("file not found")
		}
		l.Error("failed to query file information", logx.Field("error", err), logx.Field("fileID", req.FileID))
		return nil, fmt.Errorf("failed to query file information: %w", err)
	}

	// 2. If version number is specified, query the specific version information
	var fileVersion *model.FileVersion
	if req.VersionNumber != nil && *req.VersionNumber > 0 {
		versionNumber := int32(*req.VersionNumber)
		fileVersion, err = fileVersionQueryBuilder.Where(
			query.FileVersion.FileID.Eq(req.FileID),
			query.FileVersion.VersionNumber.Eq(versionNumber),
		).First()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("file version not found")
			}
			l.Error("failed to query file version information",
				logx.Field("error", err),
				logx.Field("fileID", req.FileID),
				logx.Field("versionNumber", versionNumber),
			)
			return nil, fmt.Errorf("failed to query file version information: %w", err)
		}
	}

	// 3. Query tags associated with the file
	var tagList []string
	err = query.Tag.Select(query.Tag.Name).
		LeftJoin(query.FileTag, query.FileTag.TagID.EqCol(query.Tag.TagID)).
		Where(query.FileTag.FileID.Eq(req.FileID)).
		Scan(&tagList)

	if err != nil {
		l.Error("failed to query file tags", logx.Field("error", err), logx.Field("fileID", req.FileID))
	}

	// 4. Assemble metadata response
	metadata := types.FileMetadata{
		FileID:      file.FileID,
		UserID:      file.UserID,
		FileName:    file.FileName,
		FileType:    file.FileType,
		ContentType: file.ContentType,
		Size:        file.Size,
		Path:        file.Path,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
		Version:     int(file.CurrentVersion),
		Status:      int16(file.Status),
		TrashedAt:   file.TrashedAt,
		Description: file.Description,
		Tags:        tagList,
	}

	// If queried a specific version, override some fields with version information
	if fileVersion != nil {
		metadata.Version = int(fileVersion.VersionNumber)
		metadata.Size = fileVersion.Size
		metadata.ContentType = fileVersion.ContentType
		metadata.CommitMessage = fileVersion.CommitMessage
	}

	resp.Metadata = metadata
	return resp, nil
}
