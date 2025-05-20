# 文件下载处理流程详解

本文档详细介绍文件下载的处理流程，并通过多个mermaid图表进行可视化说明。

## 整体流程概览

```mermaid
flowchart TD
    A[开始下载] --> B[解析请求参数]
    B --> C[验证用户身份]
    C --> D[查询文件元数据]
    D --> E{文件存在?}
    E -->|否| F[返回文件不存在错误]
    E -->|是| G{请求特定版本?}
    G -->|是| H[查询版本元数据]
    H --> I{版本存在?}
    I -->|否| J[返回版本不存在错误]
    I -->|是| K[获取版本文件路径]
    G -->|否| L[获取最新版本路径]
    K --> M[从对象存储获取文件]
    L --> M
    M --> N{获取成功?}
    N -->|否| O[返回存储访问错误]
    N -->|是| P[设置响应头]
    P --> Q[将文件流写入响应]
    Q --> R[完成下载]
```

## 详细步骤分析

### 1. 请求处理与参数解析

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Logic
    
    Client->>Handler: 发送下载请求(FileID, [VersionID])
    Handler->>Handler: 解析请求参数
    alt 参数解析成功
        Handler->>Logic: 创建下载逻辑并传递请求
        Logic->>Logic: 处理下载请求
    else 参数解析失败
        Handler-->>Client: 返回参数错误
    end
```

### 2. 文件元数据查询

```mermaid
flowchart TD
    A[初始化数据库查询] --> B[根据FileID查询文件元数据]
    B --> C{查询成功?}
    C -->|否| D{ErrRecordNotFound?}
    D -->|是| E[返回文件不存在错误]
    D -->|否| F[返回数据库查询错误]
    C -->|是| G[继续处理流程]
```

### 3. 版本处理流程

```mermaid
sequenceDiagram
    participant Logic
    participant Database
    
    Logic->>Logic: 检查是否指定版本
    alt 指定版本
        Logic->>Logic: 构建版本ID(FileID_VersionNumber)
        Logic->>Database: 查询版本元数据
        alt 查询成功
            Database-->>Logic: 返回版本信息
            Logic->>Logic: 使用版本路径和内容类型
        else 查询失败
            Database-->>Logic: 返回错误
            Logic-->>Client: 返回版本不存在错误
        end
    else 未指定版本
        Logic->>Logic: 使用文件最新版本路径和内容类型
    end
```

### 4. 文件获取与响应处理

```mermaid
flowchart TD
    A[确定文件路径和内容类型] --> B[从对象存储获取文件]
    B --> C{获取成功?}
    C -->|否| D[返回存储访问错误]
    C -->|是| E[设置Content-Disposition响应头]
    E --> F[设置Content-Type响应头]
    F --> G[将文件流写入HTTP响应]
    G --> H{写入成功?}
    H -->|否| I[返回写入错误]
    H -->|是| J[完成下载]
```

### 5. 对象存储交互流程

```mermaid
sequenceDiagram
    participant Logic
    participant ObjectStorage
    participant ResponseWriter
    
    Logic->>ObjectStorage: GetObject(bucket, path, options)
    alt 获取成功
        ObjectStorage-->>Logic: 返回文件对象流
        Logic->>ResponseWriter: 设置响应头
        Logic->>Logic: defer object.Close()
        Logic->>ResponseWriter: io.Copy(w, object)
        alt 复制成功
            ResponseWriter-->>Logic: 返回复制字节数
            Logic-->>Logic: 返回nil(成功)
        else 复制失败
            ResponseWriter-->>Logic: 返回错误
            Logic-->>Logic: 返回写入错误
        end
    else 获取失败
        ObjectStorage-->>Logic: 返回错误
        Logic-->>Logic: 返回存储访问错误
    end
```

## 数据库查询分析

```mermaid
flowchart LR
    A[构建基础查询] --> B[查询文件记录]
    B --> C{需要版本信息?}
    C -->|是| D[查询文件版本]
    C -->|否| E[使用文件记录]
    D --> F[确定文件路径]
    E --> F
```

## 数据库模型关系

```mermaid
erDiagram
    FILE ||--o{ FILE_VERSION : has
    FILE {
        string FileID PK
        string UserID
        string FileName
        string FileType
        string ContentType
        int64 Size
        string Path
        int64 CreatedAt
        int64 UpdatedAt
        int32 CurrentVersion
        string Description
    }
    FILE_VERSION {
        string VersionID PK
        string FileID FK
        int32 VersionNumber
        int64 Size
        string Path
        string ContentType
        int64 CreatedAt
        string CommitMessage
    }
```

## 错误处理流程

```mermaid
flowchart TD
    A[处理过程发生错误] --> B{错误类型?}
    B -->|文件不存在| C[返回404错误]
    B -->|版本不存在| D[返回404错误并指明版本]
    B -->|数据库错误| E[返回500错误并记录日志]
    B -->|存储访问错误| F[返回500错误并记录日志]
    B -->|写入响应错误| G[返回500错误并记录日志]
    C --> H[返回给客户端]
    D --> H
    E --> H
    F --> H
    G --> H
```

## 成功响应流程

```mermaid
flowchart LR
    A[获取文件内容] --> B[设置Content-Disposition]
    B --> C[设置Content-Type]
    C --> D[将文件流写入响应]
    D --> E[客户端开始下载]
```

## HTTP处理器与业务逻辑分离

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Logic
    participant DB
    participant OSS
    
    Client->>Handler: 请求下载文件
    Handler->>Handler: 解析请求参数
    Handler->>Logic: 创建Logic实例并传递请求
    Logic->>DB: 查询文件元数据
    DB-->>Logic: 返回元数据
    alt 需要特定版本
        Logic->>DB: 查询版本信息
        DB-->>Logic: 返回版本信息
    end
    Logic->>OSS: 获取文件内容
    OSS-->>Logic: 返回文件流
    Logic->>Client: 直接写入文件内容到响应
```

## 关键点说明

1. **直接流式响应**: 下载逻辑直接将文件流写入HTTP响应，而不是返回JSON数据
2. **版本支持**: 支持下载指定版本的文件，实现文件历史版本访问
3. **合适的响应头**: 设置正确的Content-Type和Content-Disposition头，确保浏览器正确处理下载
4. **流式处理**: 使用io.Copy进行流式传输，避免将整个文件加载到内存
5. **错误处理**: 详细的错误检查和处理，确保用户获得清晰的错误提示

整个下载流程设计考虑了效率和用户体验，通过直接流式响应减少内存占用，支持版本化下载增强了系统功能性。
