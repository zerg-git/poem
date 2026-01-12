# API接口文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **数据格式**: JSON
- **字符编码**: UTF-8

## 统一响应格式

### 成功响应

```json
{
  "success": true,
  "data": { ... }
}
```

### 错误响应

```json
{
  "success": false,
  "error": "错误信息"
}
```

## 接口列表

### 1. 朝代接口

#### 1.1 获取所有朝代

**请求**
```
GET /dynasties
```

**响应**
```json
{
  "success": true,
  "data": [
    {
      "id": "tang",
      "name": "唐代",
      "name_en": "Tang Dynasty",
      "description": "唐朝（618-907年），中国历史上最辉煌的朝代之一",
      "period": "618-907",
      "sort_order": 1
    },
    {
      "id": "song",
      "name": "宋代",
      "name_en": "Song Dynasty",
      "description": "宋朝（960-1279年）",
      "period": "960-1279",
      "sort_order": 2
    }
  ]
}
```

#### 1.2 获取朝代详情

**请求**
```
GET /dynasties/:id
```

**路径参数**
- `id`: 朝代ID（如：tang, song）

**响应**
```json
{
  "success": true,
  "data": {
    "id": "tang",
    "name": "唐代",
    "name_en": "Tang Dynasty",
    "description": "唐朝（618-907年）...",
    "period": "618-907",
    "poem_count": 48900,
    "author_count": 3500
  }
}
```

---

### 2. 分类接口

#### 2.1 获取所有分类

**请求**
```
GET /categories
```

**响应**
```json
{
  "success": true,
  "data": [
    {
      "id": "tang-shi",
      "name": "唐诗",
      "name_en": "Tang Poetry",
      "dynasty_id": "tang",
      "description": "全唐诗收录唐诗四万八千九百余首",
      "path": "全唐诗"
    },
    {
      "id": "song-ci",
      "name": "宋词",
      "name_en": "Song Lyrics",
      "dynasty_id": "song",
      "description": "全宋词收录宋词二万余首",
      "path": "宋词"
    }
  ]
}
```

---

### 3. 诗词接口

#### 3.1 获取诗词列表

**请求**
```
GET /poems
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量（最大100） |
| category | string | 否 | - | 分类ID（如：tang-shi） |
| dynasty | string | 否 | - | 朝代ID（如：tang） |
| author | string | 否 | - | 作者名称 |

**响应**
```json
{
  "success": true,
  "data": {
    "poems": [
      {
        "id": "181c920a-cef9-4ffc-a3ba-febc3f6cf910",
        "title": "静夜思",
        "author": "李白",
        "paragraphs": [
          "床前明月光",
          "疑是地上霜",
          "举头望明月",
          "低头思故乡"
        ],
        "category_id": "tang-shi",
        "dynasty_id": "tang"
      }
    ],
    "total": 48900,
    "page": 1,
    "page_size": 20,
    "total_pages": 2445
  }
}
```

#### 3.2 获取单首诗词

**请求**
```
GET /poems/:id
```

**路径参数**
- `id`: 诗词ID

**响应**
```json
{
  "success": true,
  "data": {
    "id": "181c920a-cef9-4ffc-a3ba-febc3f6cf910",
    "title": "静夜思",
    "author": "李白",
    "paragraphs": [
      "床前明月光",
      "疑是地上霜",
      "举头望明月",
      "低头思故乡"
    ],
    "category_id": "tang-shi",
    "dynasty_id": "tang",
    "author_info": {
      "name": "李白",
      "dynasty": "唐",
      "description": "李白（701年－762年）..."
    }
  }
}
```

#### 3.3 获取随机诗词

**请求**
```
GET /poems/random
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| count | int | 否 | 1 | 数量（最大10） |
| category | string | 否 | - | 分类ID |
| dynasty | string | 否 | - | 朝代ID |

**响应**
```json
{
  "success": true,
  "data": {
    "poems": [
      {
        "id": "...",
        "title": "春晓",
        "author": "孟浩然",
        "paragraphs": [...]
      }
    ]
  }
}
```

---

### 4. 作者接口

#### 4.1 获取作者列表

**请求**
```
GET /authors
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| dynasty | string | 否 | - | 朝代ID |

**响应**
```json
{
  "success": true,
  "data": {
    "authors": [
      {
        "id": "f78aa699-e012-4059-9e29-5d30e16cc1d8",
        "name": "李白",
        "dynasty": "唐",
        "description": "李白（701年－762年），字太白...",
        "poem_count": 1200
      }
    ],
    "total": 3500,
    "page": 1,
    "page_size": 20,
    "total_pages": 175
  }
}
```

#### 4.2 获取作者详情

**请求**
```
GET /authors/:id
```

**响应**
```json
{
  "success": true,
  "data": {
    "id": "...",
    "name": "李白",
    "dynasty": "唐",
    "description": "李白（701年－762年），字太白...",
    "poem_count": 1200
  }
}
```

#### 4.3 获取作者的诗词

**请求**
```
GET /authors/:id/poems
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**响应**
```json
{
  "success": true,
  "data": {
    "poems": [...],
    "total": 1200,
    "page": 1,
    "page_size": 20,
    "total_pages": 60
  }
}
```

---

### 5. 搜索接口

#### 5.1 全文搜索

**请求**
```
GET /search
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| q | string | 是 | - | 搜索关键词 |
| type | string | 否 | all | 搜索类型：all/poems/authors |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |

**响应**
```json
{
  "success": true,
  "data": {
    "poems": [...],
    "authors": [...],
    "total": 150,
    "page": 1,
    "page_size": 20,
    "total_pages": 8,
    "query": "明月",
    "duration_ms": 45
  }
}
```

#### 5.2 搜索建议

**请求**
```
GET /search/suggest
```

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|------|------|------|--------|------|
| q | string | 是 | - | 部分关键词 |

**响应**
```json
{
  "success": true,
  "data": {
    "suggestions": [
      "明月几时有",
      "明月光",
      "李白 - 明月..."
    ]
  }
}
```

---

## 错误码

| 错误码 | 描述 |
|--------|------|
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 请求示例

### cURL

```bash
# 获取诗词列表
curl "http://localhost:8080/api/v1/poems?page=1&page_size=10"

# 获取随机诗词
curl "http://localhost:8080/api/v1/poems/random?count=3"

# 搜索
curl "http://localhost:8080/api/v1/search?q=明月"
```

### JavaScript/Axios

```javascript
// 获取诗词列表
const response = await axios.get('http://localhost:8080/api/v1/poems', {
  params: {
    page: 1,
    page_size: 20,
    category: 'tang-shi'
  }
})

// 获取单首诗词
const poem = await axios.get(`http://localhost:8080/api/v1/poems/${id}`)
```
