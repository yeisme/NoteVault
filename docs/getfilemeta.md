---
mode: 'ask'
---

# 获取文件元数据处理流程详解

下面我将详细介绍获取文件元数据的处理流程，并使用多个mermaid图表来可视化整个过程。

## 整体流程概览

```mermaid
flowchart TD
    A[开始请求] --> B[解析请求参数]
    B --> C[查询文件基本信息]
    C --> D{文件是否存在?}
    D -->|否| E[返回文件未找到错误]
    D -->|是| F{是否指定版本号?}
    F -->|是| G[查询指定版本信息]
    F -->|否| H[查询文件标签]
    G --> I{版本是否存在?}
    I -->|否| J[返回版本未找到错误]
    I -->|是| H
    H --> K[组装元数据响应]
    K --> L[返回响应]
```

## 详细步骤分析

### 1. 请求参数解析

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Logic
    Client->>Handler: GET /api/v1/files/metadata/{fileID}?version_number=X
    Handler->>Handler: 解析路径参数 FileID
    Handler->>Handler: 解析查询参数 VersionNumber
    Handler->>Logic: 传递 GetFileMetadataRequest
    Note over Logic: req.FileID: 必需参数<br/>req.VersionNumber: 可选参数
```

### 2. 文件基本信息查询

```mermaid
flowchart TD
    A[初始化查询构建器] --> B[构建文件查询条件]
    B --> C[deleted_at = 0 AND file_id = ?]
    C --> D[执行查询]
    D --> E{查询结果}
    E -->|找到记录| F[获取文件信息]
    E -->|未找到| G[返回 gorm.ErrRecordNotFound]
    F --> H[继续下一步]
    G --> I[转换为业务错误: file not found]
```

### 3. 版本信息查询（可选）

```mermaid
sequenceDiagram
    participant Logic
    participant Database
    participant FileVersionQuery
    alt 指定了版本号且版本号 > 0
        Logic->>FileVersionQuery: 构建版本查询条件
        Note over FileVersionQuery: deleted_at = 0<br/>AND file_id = ?<br/>AND version_number = ?
        FileVersionQuery->>Database: 执行查询
        alt 找到版本记录
            Database-->>Logic: 返回版本信息
            Logic-->>Logic: 准备使用版本数据
        else 未找到版本记录
            Database-->>Logic: 返回 gorm.ErrRecordNotFound
            Logic-->>Client: 返回 "file version not found"
        end
    else 未指定版本号或版本号 <= 0
        Logic-->>Logic: 跳过版本查询，使用当前版本
    end
```

### 4. 标签信息查询

```mermaid
flowchart TD
    A[构建标签查询] --> B[左连接 tags 和 file_tags 表]
    B --> C[SELECT tags.name]
    C --> D[WHERE file_tags.file_id = ?]
    D --> E[执行查询并扫描到字符串数组]
    E --> F{查询成功?}
    F -->|是| G[获取标签列表]
    F -->|否| H[记录错误日志但继续流程]
    G --> I[准备组装响应]
    H --> I
```

### 5. 元数据组装过程

```mermaid
flowchart TD
    A[开始组装元数据] --> B[使用文件基本信息填充]
    B --> C[设置基础字段]
    C --> D{是否有版本信息?}
    D -->|是| E[用版本信息覆盖部分字段]
    D -->|否| F[保持文件原始信息]
    E --> G[设置版本特定字段]
    F --> H[添加标签信息]
    G --> H
    H --> I[完成元数据组装]
```

## 数据库查询关系图

```mermaid
erDiagram
    FILES ||--o{ FILE_VERSIONS : "has versions"
    FILES ||--o{ FILE_TAGS : "has tags"
    TAGS ||--o{ FILE_TAGS : "tagged to files"
    
    FILES {
        string file_id PK
        string user_id
        string file_name
        string file_type
        string content_type
        bigint size
        string path
        bigint created_at
        bigint updated_at
        bigint deleted_at
        smallint status
        bigint trashed_at
        int current_version
        text description
    }
    
    FILE_VERSIONS {
        string version_id PK
        string file_id FK
        int version_number
        bigint size
        string path
        string content_type
        bigint created_at
        bigint deleted_at
        smallint status
        text commit_message
    }
    
    TAGS {
        string tag_id PK
        string name
    }
    
    FILE_TAGS {
        string file_id FK
        string tag_id FK
    }
```

## 响应数据结构

```mermaid
classDiagram
    class GetFileMetadataResponse {
        +FileMetadata metadata
    }
    
    class FileMetadata {
        +string FileID
        +string UserID
        +string FileName
        +string FileType
        +string ContentType
        +int64 Size
        +string Path
        +int64 CreatedAt
        +int64 UpdatedAt
        +int Version
        +int16 Status
        +int64 TrashedAt
        +string Description
        +string CommitMessage
        +[]string Tags
    }
    
    GetFileMetadataResponse --> FileMetadata
```

## 错误处理流程

```mermaid
flowchart TD
    A[捕获错误] --> B{错误类型判断}
    B -->|gorm.ErrRecordNotFound 文件查询| C[返回 file not found]
    B -->|gorm.ErrRecordNotFound 版本查询| D[返回 file version not found]
    B -->|数据库连接错误| E[记录错误日志]
    B -->|其他数据库错误| F[记录错误日志]
    E --> G[返回包装后的错误信息]
    F --> G
    C --> H[返回格式化错误]
    D --> H
    G --> H
```

## 查询优化策略

```mermaid
flowchart TD
    A[查询优化] --> B[使用软删除过滤]
    B --> C[deleted_at = 0]
    C --> D[利用索引优化]
    D --> E[file_id 主键索引]
    E --> F[复合索引查询版本]
    F --> G[file_id + version_number]
    G --> H[标签查询使用 LEFT JOIN]
    H --> I[避免 N+1 查询问题]
```

## 版本信息覆盖逻辑

```mermaid
sequenceDiagram
    participant Logic
    participant FileMetadata
    participant FileVersion
    Logic->>FileMetadata: 设置文件基本信息
    Note over FileMetadata: FileID, UserID, FileName<br/>FileType, Path, CreatedAt<br/>UpdatedAt, Status, etc.
    alt 查询到特定版本
        Logic->>FileVersion: 获取版本特定信息
        FileVersion-->>Logic: Version, Size, ContentType<br/>CommitMessage
        Logic->>FileMetadata: 覆盖版本相关字段
        Note over FileMetadata: metadata.Version = fileVersion.VersionNumber<br/>metadata.Size = fileVersion.Size<br/>metadata.ContentType = fileVersion.ContentType<br/>metadata.CommitMessage = fileVersion.CommitMessage
    else 使用当前版本
        Logic->>FileMetadata: 保持文件表中的信息
        Note over FileMetadata: metadata.Version = file.CurrentVersion<br/>其他字段保持不变
    end
```

## API 调用示例流程

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Logic
    participant Database
    
    Client->>Handler: GET /api/v1/files/metadata/abc123?version_number=2
    Handler->>Logic: GetFileMetadata(FileID: "abc123", VersionNumber: 2)
    
    Logic->>Database: SELECT * FROM files WHERE file_id="abc123" AND deleted_at=0
    Database-->>Logic: 返回文件基本信息
    
    Logic->>Database: SELECT * FROM file_versions WHERE file_id="abc123" AND version_number=2 AND deleted_at=0
    Database-->>Logic: 返回版本信息
    
    Logic->>Database: SELECT tags.name FROM tags LEFT JOIN file_tags ON ... WHERE file_tags.file_id="abc123"
    Database-->>Logic: 返回标签列表 ["document", "important"]
    
    Logic->>Logic: 组装元数据响应
    Logic-->>Handler: 返回完整的文件元数据
    Handler-->>Client: JSON响应包含所有元数据信息
```

## 关键特性说明

1. **软删除支持**：所有查询都过滤 `deleted_at = 0` 的记录
2. **版本控制**：支持查询特定版本的元数据信息
3. **标签系统**：通过关联表查询文件的所有标签
4. **错误处理**：区分文件不存在和版本不存在的错误
5. **数据一致性**：确保返回的元数据信息完整准确
6. **查询优化**：使用适当的索引和查询策略提高性能

整个流程设计注重数据完整性和查询效率，为用户提供准确的文件元数据信息。
