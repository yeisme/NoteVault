### 1. "列出文件，支持分页、筛选和排序。"

1. route definition

- Url: /api/v1/files
- Method: GET
- Request: `ListFilesRequest`
- Response: `ListFilesResponse`

2. request definition



```golang
type ListFilesRequest struct {
	UserID string `form:"userId,optional"` // 按用户ID筛选（管理员可能会使用）
	FileName string `form:"fileName,optional"` // 按文件名模糊匹配
	FileType string `form:"fileType,optional"` // 文件类型精确匹配
	Tag string `form:"tag,optional"` // 按单个标签精确匹配 (未来可支持多标签)
	CreatedAtStart int64 `form:"createdAtStart,optional"` // 创建时间范围开始 (Unix timestamp)
	CreatedAtEnd int64 `form:"createdAtEnd,optional"` // 创建时间范围结束 (Unix timestamp)
	UpdatedAtStart int64 `form:"updatedAtStart,optional"` // 更新时间范围开始 (Unix timestamp)
	UpdatedAtEnd int64 `form:"updatedAtEnd,optional"` // 更新时间范围结束 (Unix timestamp)
	Page int `form:"page,default=1"` // 页码
	PageSize int `form:"pageSize,default=10"` // 每页大小
	SortBy string `form:"sortBy,optional,options=name|date|size|type"` // 排序字段: name, date (updatedAt), size, type
	Order string `form:"order,optional,options=asc|desc"` // 排序顺序
}
```


3. response definition



```golang
type ListFilesResponse struct {
	Files []FileMetadata `json:"files"`
	TotalCount int64 `json:"totalCount"`
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
}
```

### 2. "根据文件ID删除文件。"

1. route definition

- Url: /api/v1/files/:fileId
- Method: DELETE
- Request: `FileDeleteRequest`
- Response: `FileDeleteResponse`

2. request definition



```golang
type FileDeleteRequest struct {
	FileID string `path:"fileId"`
}
```


3. response definition



```golang
type FileDeleteResponse struct {
	Message string `json:"message"`
}
```

### 3. "获取文件的版本历史。"

1. route definition

- Url: /api/v1/files/:fileId/versions
- Method: GET
- Request: `GetFileVersionsRequest`
- Response: `GetFileVersionsResponse`

2. request definition



```golang
type GetFileVersionsRequest struct {
	FileID string `path:"fileId"`
}
```


3. response definition



```golang
type GetFileVersionsResponse struct {
	FileID string `json:"fileId"`
	Versions []FileVersionInfo `json:"versions"`
}
```

### 4. "(高级) 获取文件两个版本之间的差异信息 (主要针对文本文件)。"

1. route definition

- Url: /api/v1/files/:fileId/versions/diff
- Method: GET
- Request: `FileVersionDiffRequest`
- Response: `FileVersionDiffResponse`

2. request definition



```golang
type FileVersionDiffRequest struct {
	FileID string `path:"fileId"`
	BaseVersion int `form:"baseVersion"` // 基础版本号
	TargetVersion int `form:"targetVersion"` // 目标版本号
}
```


3. response definition



```golang
type FileVersionDiffResponse struct {
	FileID string `json:"fileId"`
	BaseVersion int `json:"baseVersion"`
	TargetVersion int `json:"targetVersion"`
	DiffContent string `json:"diffContent"` // 差异内容 (例如 unified diff 格式)
	Message string `json:"message,optional"`
}
```

### 5. "将文件恢复到特定版本。"

1. route definition

- Url: /api/v1/files/:fileId/versions/revert
- Method: POST
- Request: `RevertFileVersionRequest`
- Response: `RevertFileVersionResponse`

2. request definition



```golang
type RevertFileVersionRequest struct {
	FileID string `path:"fileId"`
	Version int `json:"version"` // 要恢复到的版本号
	CommitMessage string `json:"commitMessage,optional"` // 恢复操作的提交信息
}
```


3. response definition



```golang
type RevertFileVersionResponse struct {
	Metadata FileMetadata `json:"metadata"` // 恢复后，文件当前的元数据（版本已更新）
	Message string `json:"message"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // 文件唯一ID
	UserID string `json:"userId"` // 文件所属用户ID
	FileName string `json:"fileName"` // 文件名
	FileType string `json:"fileType"` // 文件类型，例如："document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME类型，例如："application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // 文件大小（字节）
	Path string `json:"path"` // 存储路径或键
	CreatedAt int64 `json:"createdAt"` // 创建时间（Unix时间戳）
	UpdatedAt int64 `json:"updatedAt"` // 更新时间（Unix时间戳）
	Version int `json:"version"` // 文件当前版本号
	Tags []string `json:"tags,optional"` // 标签
	Description string `json:"description,optional"` // 描述
}
```

### 6. "批量删除文件。"

1. route definition

- Url: /api/v1/files/batch/delete
- Method: POST
- Request: `BatchDeleteFilesRequest`
- Response: `BatchDeleteFilesResponse`

2. request definition



```golang
type BatchDeleteFilesRequest struct {
	FileIDs []string `json:"fileIds"`
}
```


3. response definition



```golang
type BatchDeleteFilesResponse struct {
	Succeeded []string `json:"succeeded"` // 成功删除的文件ID列表
	Failed []string `json:"failed"` // 删除失败的文件ID列表 (及原因，可选)
	Message string `json:"message"`
}
```

### 7. "根据文件ID下载文件。可选下载特定版本。"

1. route definition

- Url: /api/v1/files/download/:fileId
- Method: GET
- Request: `FileDownloadRequest`
- Response: `-`

2. request definition



```golang
type FileDownloadRequest struct {
	FileID string `path:"fileId"`
	VersionID *int `form:"versionId,optional"` // 可选，指定下载特定版本的文件
}
```


3. response definition


### 8. "获取特定文件的元数据。可选获取特定版本的元数据。"

1. route definition

- Url: /api/v1/files/metadata/:fileId
- Method: GET
- Request: `GetFileMetadataRequest`
- Response: `GetFileMetadataResponse`

2. request definition



```golang
type GetFileMetadataRequest struct {
	FileID string `path:"fileId"`
	VersionID *int `form:"versionId,optional"` // 可选，获取特定版本文件的元数据
}
```


3. response definition



```golang
type GetFileMetadataResponse struct {
	Metadata FileMetadata `json:"metadata"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // 文件唯一ID
	UserID string `json:"userId"` // 文件所属用户ID
	FileName string `json:"fileName"` // 文件名
	FileType string `json:"fileType"` // 文件类型，例如："document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME类型，例如："application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // 文件大小（字节）
	Path string `json:"path"` // 存储路径或键
	CreatedAt int64 `json:"createdAt"` // 创建时间（Unix时间戳）
	UpdatedAt int64 `json:"updatedAt"` // 更新时间（Unix时间戳）
	Version int `json:"version"` // 文件当前版本号
	Tags []string `json:"tags,optional"` // 标签
	Description string `json:"description,optional"` // 描述
}
```

### 9. "更新特定文件的元数据。这通常会创建一个新版本。"

1. route definition

- Url: /api/v1/files/metadata/:fileId
- Method: PUT
- Request: `UpdateFileMetadataRequest`
- Response: `UpdateFileMetadataResponse`

2. request definition



```golang
type UpdateFileMetadataRequest struct {
	FileID string `path:"fileId"`
	FileName string `json:"fileName,optional"`
	Description string `json:"description,optional"`
	Tags []string `json:"tags,optional"`
	CommitMessage string `json:"commitMessage,optional"` // 版本提交信息
}
```


3. response definition



```golang
type UpdateFileMetadataResponse struct {
	Metadata FileMetadata `json:"metadata"` // 更新后的元数据，包含新的版本号
	Message string `json:"message"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // 文件唯一ID
	UserID string `json:"userId"` // 文件所属用户ID
	FileName string `json:"fileName"` // 文件名
	FileType string `json:"fileType"` // 文件类型，例如："document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME类型，例如："application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // 文件大小（字节）
	Path string `json:"path"` // 存储路径或键
	CreatedAt int64 `json:"createdAt"` // 创建时间（Unix时间戳）
	UpdatedAt int64 `json:"updatedAt"` // 更新时间（Unix时间戳）
	Version int `json:"version"` // 文件当前版本号
	Tags []string `json:"tags,optional"` // 标签
	Description string `json:"description,optional"` // 描述
}
```

### 10. "上传一个新文件。实际文件以 multipart/form-data 形式发送。"

1. route definition

- Url: /api/v1/files/upload
- Method: POST
- Request: `FileUploadRequest`
- Response: `FileUploadResponse`

2. request definition



```golang
type FileUploadRequest struct {
	FileName string `form:"fileName,optional"` // 可选：如果未提供，则使用上传文件的名称
	FileType string `form:"fileType,optional"` // 可选：可以推断或指定
	Description string `form:"description,optional"` // 描述
	Tags string `form:"tags,optional"` // 逗号分隔的标签
	CommitMessage string `form:"commitMessage,optional"` // 版本提交信息
}
```


3. response definition



```golang
type FileUploadResponse struct {
	FileID string `json:"fileId"`
	FileName string `json:"fileName"`
	ContentType string `json:"contentType"`
	Size int64 `json:"size"`
	Message string `json:"message"`
	Version int `json:"version"` // 上传后的文件版本号
}
```

