### 1. "List files with support for pagination, filtering, and sorting."

1. route definition

- Url: /api/v1/files
- Method: GET
- Request: `ListFilesRequest`
- Response: `ListFilesResponse`

2. request definition



```golang
type ListFilesRequest struct {
	UserID string `form:"userId,optional"` // Filter by user ID (may be used by administrators)
	FileName string `form:"fileName,optional"` // Fuzzy match by file name
	FileType string `form:"fileType,optional"` // Exact match by file type
	Tag string `form:"tag,optional"` // Exact match by a single tag (multiple tags may be supported in the future)
	CreatedAtStart int64 `form:"createdAtStart,optional"` // Creation time range start (Unix timestamp)
	CreatedAtEnd int64 `form:"createdAtEnd,optional"` // Creation time range end (Unix timestamp)
	UpdatedAtStart int64 `form:"updatedAtStart,optional"` // Update time range start (Unix timestamp)
	UpdatedAtEnd int64 `form:"updatedAtEnd,optional"` // Update time range end (Unix timestamp)
	Page int `form:"page,default=1"` // Page number
	PageSize int `form:"pageSize,default=10"` // Page size
	SortBy string `form:"sortBy,optional,options=name|date|size|type"` // Sort field: name, date (updatedAt), size, type
	Order string `form:"order,optional,options=asc|desc"` // Sort order
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

### 2. "Delete a file by file ID. Optionally delete a specific version."

1. route definition

- Url: /api/v1/files/:fileId
- Method: DELETE
- Request: `FileDeleteRequest`
- Response: `FileDeleteResponse`

2. request definition



```golang
type FileDeleteRequest struct {
	FileID string `path:"fileId"`
	VersionNumber *int `json:"versionNumber,optional"` // Optional, specify to delete a specific version of the file
}
```


3. response definition



```golang
type FileDeleteResponse struct {
	Message string `json:"message"`
}
```

### 3. "Get version history for a file."

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

### 4. "(Advanced) Get differences between two versions of a file (mainly for text files)."

1. route definition

- Url: /api/v1/files/:fileId/versions/diff
- Method: GET
- Request: `FileVersionDiffRequest`
- Response: `FileVersionDiffResponse`

2. request definition



```golang
type FileVersionDiffRequest struct {
	FileID string `path:"fileId"`
	BaseVersion int `form:"baseVersion"` // Base version number
	TargetVersion int `form:"targetVersion"` // Target version number
}
```


3. response definition



```golang
type FileVersionDiffResponse struct {
	FileID string `json:"fileId"`
	BaseVersion int `json:"baseVersion"`
	TargetVersion int `json:"targetVersion"`
	DiffContent string `json:"diffContent"` // Difference content (e.g., unified diff format)
	Message string `json:"message,optional"`
}
```

### 5. "Revert a file to a specific version."

1. route definition

- Url: /api/v1/files/:fileId/versions/revert
- Method: POST
- Request: `RevertFileVersionRequest`
- Response: `RevertFileVersionResponse`

2. request definition



```golang
type RevertFileVersionRequest struct {
	FileID string `path:"fileId"`
	Version int `json:"version"` // Version number to revert to
	CommitMessage string `json:"commitMessage,optional"` // Commit message for the revert operation
}
```


3. response definition



```golang
type RevertFileVersionResponse struct {
	Metadata FileMetadata `json:"metadata"` // Current file metadata after reverting (version is updated)
	Message string `json:"message"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // Unique file ID
	UserID string `json:"userId"` // ID of the user who owns the file
	FileName string `json:"fileName"` // File name
	FileType string `json:"fileType"` // File type, e.g., "document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME type, e.g., "application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // File size in bytes
	Path string `json:"path"` // Storage path or key
	CreatedAt int64 `json:"createdAt"` // Creation time (Unix timestamp)
	UpdatedAt int64 `json:"updatedAt"` // Update time (Unix timestamp)
	Version int `json:"version"` // Current file version number
	Status int16 `json:"status"` // File status: 0=normal, 1=archived, 2=trashed, 3=pending deletion
	TrashedAt int64 `json:"trashedAt,optional"` // When the file was moved to trash (Unix timestamp)
	Tags []string `json:"tags,optional"` // Tags
	Description string `json:"description,optional"` // Description
	CommitMessage string `json:"commitMessage,optional"` // Version commit message
}
```

### 6. "Batch delete files."

1. route definition

- Url: /api/v1/files/batch/delete
- Method: POST
- Request: `BatchDeleteFilesRequest`
- Response: `BatchDeleteFilesResponse`

2. request definition



```golang
type BatchDeleteFilesRequest struct {
	FileIDs []string `json:"fileIds"`
	VersionNumber *int `json:"versionNumber,optional"` // Optional, specify to delete a specific version of the files
}
```


3. response definition



```golang
type BatchDeleteFilesResponse struct {
	Succeeded []string `json:"succeeded"` // List of file IDs that were successfully deleted
	Failed []string `json:"failed"` // List of file IDs that failed to delete (and reasons, optional)
	Message string `json:"message"`
}
```

### 7. "Download a file by file ID. Optionally download a specific version."

1. route definition

- Url: /api/v1/files/download/:fileId
- Method: GET
- Request: `FileDownloadRequest`
- Response: `-`

2. request definition



```golang
type FileDownloadRequest struct {
	FileID string `path:"fileId"`
	VersionNumber *int `form:"versionNumber,optional"` // Optional, specify to download a specific version of the file
}
```


3. response definition


### 8. "Get metadata for a specific file. Optionally get metadata for a specific version."

1. route definition

- Url: /api/v1/files/metadata/:fileId
- Method: GET
- Request: `GetFileMetadataRequest`
- Response: `GetFileMetadataResponse`

2. request definition



```golang
type GetFileMetadataRequest struct {
	FileID string `path:"fileId"`
	VersionNumber *int `form:"versionNumber,optional"` // Optional, get metadata for a specific version of the file
}
```


3. response definition



```golang
type GetFileMetadataResponse struct {
	Metadata FileMetadata `json:"metadata"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // Unique file ID
	UserID string `json:"userId"` // ID of the user who owns the file
	FileName string `json:"fileName"` // File name
	FileType string `json:"fileType"` // File type, e.g., "document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME type, e.g., "application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // File size in bytes
	Path string `json:"path"` // Storage path or key
	CreatedAt int64 `json:"createdAt"` // Creation time (Unix timestamp)
	UpdatedAt int64 `json:"updatedAt"` // Update time (Unix timestamp)
	Version int `json:"version"` // Current file version number
	Status int16 `json:"status"` // File status: 0=normal, 1=archived, 2=trashed, 3=pending deletion
	TrashedAt int64 `json:"trashedAt,optional"` // When the file was moved to trash (Unix timestamp)
	Tags []string `json:"tags,optional"` // Tags
	Description string `json:"description,optional"` // Description
	CommitMessage string `json:"commitMessage,optional"` // Version commit message
}
```

### 9. "Update metadata for a specific file. This typically creates a new version."

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
	CommitMessage string `json:"commitMessage,optional"` // Version commit message
}
```


3. response definition



```golang
type UpdateFileMetadataResponse struct {
	Metadata FileMetadata `json:"metadata"` // Updated metadata, including the new version number
	Message string `json:"message"`
}

type FileMetadata struct {
	FileID string `json:"fileId"` // Unique file ID
	UserID string `json:"userId"` // ID of the user who owns the file
	FileName string `json:"fileName"` // File name
	FileType string `json:"fileType"` // File type, e.g., "document", "image", "video", "text"
	ContentType string `json:"contentType"` // MIME type, e.g., "application/pdf", "image/jpeg", "text/plain"
	Size int64 `json:"size"` // File size in bytes
	Path string `json:"path"` // Storage path or key
	CreatedAt int64 `json:"createdAt"` // Creation time (Unix timestamp)
	UpdatedAt int64 `json:"updatedAt"` // Update time (Unix timestamp)
	Version int `json:"version"` // Current file version number
	Status int16 `json:"status"` // File status: 0=normal, 1=archived, 2=trashed, 3=pending deletion
	TrashedAt int64 `json:"trashedAt,optional"` // When the file was moved to trash (Unix timestamp)
	Tags []string `json:"tags,optional"` // Tags
	Description string `json:"description,optional"` // Description
	CommitMessage string `json:"commitMessage,optional"` // Version commit message
}
```

### 10. "Upload a new file. The actual file is sent as multipart/form-data."

1. route definition

- Url: /api/v1/files/upload
- Method: POST
- Request: `FileUploadRequest`
- Response: `FileUploadResponse`

2. request definition



```golang
type FileUploadRequest struct {
	FileName string `form:"fileName,optional"` // Optional: If not provided, the name of the uploaded file will be used
	FileType string `form:"fileType,optional"` // Optional: Can be inferred or specified
	Description string `form:"description,optional"` // Description
	Tags string `form:"tags,optional"` // Comma-separated tags
	CommitMessage string `form:"commitMessage,optional"` // Version commit message
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
	Version int `json:"version"` // File version number after upload
}
```

