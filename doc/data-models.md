# 数据模型说明

## 数据集概述

本项目使用 [chinese-poetry](https://github.com/chinese-poetry/chinese-poetry) 数据集，包含超过55万首中国古诗词。

## 数据集结构

```
chinese-poetry-master/
├── 全唐诗/           # 约5.5万首唐诗
│   ├── poet.tang.0.json
│   ├── poet.tang.1000.json
│   └── ...
├── 宋词/             # 约21万首宋词
│   ├── ci.song.0.json
│   ├── ci.song.1000.json
│   └── ...
├── 诗经/             # 诗经305篇
│   └── shijing.json
├── 楚辞/             # 楚辞
│   └── ...
├── 元曲/             # 元曲
│   └── ...
├── 五代诗词/         # 五代诗词
├── 蒙学/             # 蒙学经典
├── 纳兰性德/         # 纳兰性德词集
├── 曹操诗集/         # 曹操诗集
├── 论语/             # 论语
├── 四书五经/         # 四书五经
└── 幽梦影/           # 幽梦影
```

## 数据模型定义

### 1. 统一诗词模型 (Poem)

```go
type Poem struct {
    // 基础字段
    ID        string   `json:"id"`
    Title     string   `json:"title"`
    Author    string   `json:"author"`
    Paragraphs []string `json:"paragraphs"`

    // 可选字段（不同分类可能有不同字段）
    Content   []string `json:"content,omitempty"`   // 诗经使用
    Para      []string `json:"para,omitempty"`      // 纳兰性德使用
    Rhythmic  string   `json:"rhythmic,omitempty"`  // 宋词词牌名
    Dynasty   string   `json:"dynasty,omitempty"`   // 元曲朝代字段
    Tags      []string `json:"tags,omitempty"`      // 标签（唐诗三百首）
    Chapter   string   `json:"chapter,omitempty"`   // 章节（诗经）
    Section   string   `json:"section,omitempty"`   # 章节（诗经）

    // 元数据
    CategoryID string `json:"category_id,omitempty"`
    DynastyID  string `json:"dynasty_id,omitempty"`
    AuthorID   string `json:"author_id,omitempty"`
}
```

### 2. 作者模型 (Author)

```go
type Author struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Dynasty     string `json:"dynasty"`
    Description string `json:"desc"`
    ShortDesc   string `json:"short_description,omitempty"`
    PoemCount   int    `json:"poem_count"`
}
```

### 3. 朝代模型 (Dynasty)

```go
type Dynasty struct {
    ID          string `json:"id"`
    Name        string `json:"name"`        // 中文名：唐代
    NameEn      string `json:"name_en"`     // 英文名：Tang
    Description string `json:"description"`
    Period      string `json:"period"`      // 时间段：618-907
    SortOrder   int    `json:"sort_order"`
}
```

### 4. 分类模型 (Category)

```go
type Category struct {
    ID          string `json:"id"`
    Name        string `json:"name"`        // 唐诗、宋词
    NameEn      string `json:"name_en"`     // Tang Poetry
    Description string `json:"description"`
    DynastyID   string `json:"dynasty_id"`
    Path        string `json:"path"`        // 文件路径
}
```

## 原始JSON格式

### 唐诗格式 (poet.tang.*.json)

```json
[
  {
    "author": "陶翰",
    "paragraphs": [
      "乘君手中白马，客本从龙离沙碛。",
      "天子不异玉珂，丈夫辞官从赤骥。",
      ...
    ],
    "title": "出塞曲",
    "id": "181c920a-cef9-4ffc-a3ba-febc3f6cf910"
  }
]
```

### 宋词格式 (ci.song.*.json)

```json
[
  {
    "author": "丁谓",
    "paragraphs": [
      "酷暑愁深，天低不散耕云。",
      "谁解昭阳日暮，更修篁、依约佳人。",
      ...
    ],
    "rhythmic": "威氏（酷暑愁深）",
    "title": "威氏"
  }
]
```

### 诗经格式 (shijing.json)

```json
[
  {
    "chapter": "小雅",
    "content": [
      "文王有声，遹骏有声。",
      "遹求厥宁，遹观厥成。",
      ...
    ],
    "section": "大雅",
    "title": "文王有声"
  }
]
```

### 元曲格式

```json
[
  {
    "author": "关汉卿",
    "paragraphs": [
      "恰离了绿水青山那搭，早来到竹篱茅舍人家。",
      ...
    ],
    "title": "一枝花·离恨",
    "dynasty": "元",
    "id": "xxx"
  }
]
```



## 数据统计

| 分类 | 诗词数量 | 作者数量 | 文件数量 |
|------|----------|----------|----------|
| 全唐诗 | ~48,900 | ~2,200 | ~5 |
| 全宋词 | ~21,000 | ~1,300 | ~200 |
| 诗经 | 305 | - | 1 |
| 楚辞 | ~50 | ~3 | ~1 |
| 元曲 | ~3,800 | ~200 | ~3 |
| 五代诗词 | ~1,100 | ~100 | ~1 |
| 唐诗三百首 | 313 | 77 | 1 |
| 宋词三百首 | 300+ | ~80 | 1 |

## 数据索引策略

### 1. 作者索引

```go
type AuthorIndex struct {
    Authors map[string]*models.Author // name -> Author
    ByDynasty map[string][]string     // dynasty -> []authorName
}
```

### 2. 诗词标题索引

```go
type TitleIndex struct {
    Titles map[string][]string       // title -> []poemID
    Pinyin map[string][]string       // pinyin -> []poemID
}
```

### 3. 全文搜索索引

```go
type SearchIndex struct {
    Words map[string][]string        // word -> []poemID
    Authors map[string][]string      // author -> []poemID
}
```

## 数据验证

### JSON格式验证

```python
# Python验证脚本示例
import json

def validate_poem(data):
    required_fields = ['title', 'author', 'paragraphs']
    for field in required_fields:
        if field not in data:
            return False, f"Missing field: {field}"
    return True, "Valid"

with open('poet.tang.0.json', 'r', encoding='utf-8') as f:
    poems = json.load(f)
    for poem in poems:
        valid, msg = validate_poem(poem)
        if not valid:
            print(f"Invalid: {msg}")
```

## 数据更新

### 从源更新

```bash
# 更新chinese-poetry数据集
cd chinese-poetry-master
git pull origin master
```

### 数据同步

```go
// 检测数据文件变化
func (jl *JSONLoader) WatchChanges() {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(jl.basePath)

    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                // 重新加载文件
                jl.reloadFile(event.Name)
            }
        }
    }
}
```
