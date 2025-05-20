# 文件列表查询流程详解

本文档详细解释文件列表查询的处理流程。

## 整体流程概览

```mermaid
flowchart TD
    A[开始查询] --> B[解析请求参数]
    B --> C[验证分页参数]
    C --> D[初始化查询构建器]
    D --> E[应用过滤条件]
    E --> F[处理标签过滤]
    F --> G[计算结果总数]
    G --> H[设置排序规则]
    H --> I[应用分页限制]
    I --> J[执行文件查询]
    J --> K[查询关联标签]
    K --> L[组装响应数据]
    L --> M[返回结果]
```

## 详细步骤分析

### 1. 请求参数解析与验证

```mermaid
flowchart TD
    A[接收请求] --> B[解析请求参数]
    B --> C{验证分页参数}
    C -->|Page < 1| D[设置 Page = 1]
    C -->|PageSize 无效| E[设置 PageSize = 10]
    C -->|参数有效| F[保持原值]
    D --> G[继续处理]
    E --> G
    F --> G
```

### 2. 过滤条件应用

```mermaid
flowchart TD
    A[初始化查询构建器] --> B{用户ID过滤}
    B -->|有| C[添加用户ID条件]
    B -->|无| D{文件名过滤}
    C --> D
    D -->|有| E[添加文件名模糊匹配]
    D -->|无| F{文件类型过滤}
    E --> F
    F -->|有| G[添加文件类型条件]
    F -->|无| H{时间范围过滤}
    G --> H
    H -->|有创建时间起始| I[添加创建时间>=条件]
    H -->|无| J{其他时间条件}
    I --> J
    J -->|有创建时间结束| K[添加创建时间<=条件]
    J -->|无| L{更新时间条件}
    K --> L
    L -->|有更新时间| M[添加更新时间条件]
    L -->|无| N[继续处理]
    M --> N
```

### 3. 标签过滤处理

```mermaid
sequenceDiagram
    participant Logic
    participant Database
    Logic->>Logic: 检查是否有标签过滤
    alt 有标签过滤
        Logic->>Database: 查询标签ID
        alt 标签存在
            Database-->>Logic: 返回标签信息
            Logic->>Logic: 关联文件标签表
            Logic->>Logic: 添加标签ID条件
        else 标签不存在
            Database-->>Logic: 未找到标签
            Logic->>Logic: 准备返回空结果
        end
    else 无标签过滤
        Logic->>Logic: 跳过标签处理
    end
```

### 4. 排序和分页处理

```mermaid
flowchart TD
    A[接收排序参数] --> B{排序字段?}
    B -->|name| C[按文件名排序]
    B -->|size| D[按文件大小排序]
    B -->|type| E[按文件类型排序]
    B -->|date或默认| F[按更新时间排序]
    C --> G{排序方向}
    D --> G
    E --> G
    F --> G
    G -->|asc| H[升序排列]
    G -->|desc| I[降序排列]
    H --> J[计算偏移量]
    I --> J
    J --> K[应用LIMIT和OFFSET]
```

### 5. 查询执行和标签关联

```mermaid
sequenceDiagram
    participant Logic
    participant DB
    Logic->>DB: 执行文件查询
    DB-->>Logic: 返回文件列表
    loop 对每个文件
        Logic->>DB: 查询关联的标签
        DB-->>Logic: 返回文件的标签关系
        loop 对每个标签关系
            Logic->>DB: 查询标签详情
            DB-->>Logic: 返回标签信息
            Logic->>Logic: 添加到标签列表
        end
        Logic->>Logic: 构建文件元数据
    end
    Logic->>Logic: 组装完整响应
```

## 数据库查询分析

```mermaid
flowchart TD
    A[构建基础查询] --> B[添加WHERE条件]
    B --> C[添加JOIN条件]
    C --> D[COUNT查询获取总数]
    D --> E[添加ORDER BY]
    E --> F[添加LIMIT和OFFSET]
    F --> G[执行最终查询]
    G --> H[标签关联查询]
```

## 关键结构示意图

```mermaid
erDiagram
    FILE ||--o{ FILE_TAG : has
    TAG ||--o{ FILE_TAG : has
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
    TAG {
        string TagID PK
        string Name
    }
    FILE_TAG {
        string FileID FK
        string TagID FK
    }
```

## 过滤和排序处理流程

```mermaid
flowchart LR
    A[原始查询] --> B[用户过滤]
    B --> C[文件名过滤]
    C --> D[文件类型过滤]
    D --> E[时间范围过滤]
    E --> F[标签过滤]
    F --> G[排序处理]
    G --> H[分页处理]
    H --> I[执行查询]
```

## 响应组装流程

```mermaid
flowchart TD
    A[获取查询结果] --> B[初始化返回数组]
    B --> C[循环处理每个文件]
    C --> D[查询文件标签]
    D --> E[组装文件元数据]
    E --> F[添加到响应数组]
    F --> G{还有更多文件?}
    G -->|是| C
    G -->|否| H[设置分页信息]
    H --> I[设置总数信息]
    I --> J[返回完整响应]
```

## 错误处理流程

```mermaid
flowchart TD
    A[执行操作] --> B{发生错误?}
    B -->|是| C[记录日志]
    B -->|否| D[继续处理]
    C --> E{错误类型?}
    E -->|数据库错误| F[记录详细错误]
    E -->|查询构建错误| G[记录构建错误]
    F --> H[返回错误响应]
    G --> H
    D --> I[成功返回]
```

## 关键点说明

1. **高效过滤**：支持多种过滤条件，包括用户ID、文件名、文件类型、时间范围和标签
2. **灵活排序**：支持按文件名、大小、类型和日期的升序/降序排序
3. **标准分页**：实现标准的分页机制，控制返回结果数量
4. **标签关联**：通过关联查询获取每个文件的标签信息
5. **性能考虑**：先计算总数再获取分页数据，避免不必要的数据传输

整个查询流程设计考虑了灵活性、性能和数据完整性，确保文件列表查询结果准确且高效。
