# Chinese Poetry 数据结构文档

## 概述

`chinese-poetry-master` 是一个包含中国古诗词的 JSON 数据集，涵盖从先秦到清代的各类文学作品。

## 目录结构

```
chinese-poetry-master/
├── 全唐诗/           # 唐诗数据
├── 宋词/             # 宋词数据
├── 元曲/             # 元曲数据
├── 诗经/             # 诗经数据
├── 楚辞/             # 楚辞数据
├── 五代诗词/         # 五代十国诗词
├── 纳兰性德/         # 纳兰性德诗集
├── 蒙学/             # 启蒙读物
├── 四书五经/         # 四书五经
├── 论语/             # 论语
├── 曹操诗集/         # 曹操诗集
├── 幽梦影/           # 幽梦影
├── 水墨唐诗/         # 水墨唐诗
└── 御定全唐詩/       # 御定全唐诗
```

---

## 数据结构详解

### 1. 全唐诗 (全唐诗/)

**文件命名规则：**
- `poet.tang.{编号}.json` - 唐诗数据
- `poet.song.{编号}.json` - 宋诗数据（包含在全唐诗目录中）
- `authors.tang.json` - 唐代作者信息
- `authors.song.json` - 宋代作者信息
- `唐诗三百首.json` - 唐诗三百首精选
- `唐诗补录.json` - 唐诗补录
- `error/` - 错误数据目录（应忽略）

**JSON 结构：**

```json
[
  {
    "id": "181c920a-cef9-4ffc-a3ba-febc3f6cf910",
    "title": "经殺子谷",
    "author": "陶翰",
    "paragraphs": [
      "扶蘇秦帝子，舉代稱其賢。",
      "百萬猶在握，可爭天下權。",
      "束身就一劒，壯志皆棄捐。",
      "塞下有遺跡，千齡人共傳。",
      "疎蕪盡荒草，寂歷空寒煙。",
      "到此盡垂淚，非我獨潸然。"
    ],
    "tags": ["五言古诗"]
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 否 | 唯一标识符（UUID） |
| title | string | 是 | 诗词标题 |
| author | string | 是 | 作者名称 |
| paragraphs | string[] | 是 | 诗词正文（每句/联一个元素） |
| tags | string[] | 否 | 标签（如：五言古诗、七言绝句等） |

---

### 2. 宋词 (宋词/)

**文件命名规则：**
- `ci.song.{编号}.json` - 宋词数据
- `author.song.json` - 宋词作者信息
- `宋词三百首.json` - 宋词三百首精选

**JSON 结构：**

```json
[
  {
    "author": "和岘",
    "paragraphs": [
      "气和玉烛，睿化著鸿明。",
      "缇管一阳生。",
      "郊禋盛礼燔柴毕，旋轸凤凰城。",
      "森罗仪卫振华缨。",
      "载路溢欢声。"
    ],
    "rhythmic": "导引"
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| author | string | 是 | 作者名称 |
| paragraphs | string[] | 是 | 词牌正文 |
| rhythmic | string | 是 | 词牌名（注：宋词中 rhythmic 用作标题） |
| title | string | 否 | 诗词标题（部分数据有此字段） |

**注意：** 宋词数据中很多条目没有 `title` 字段，使用 `rhythmic`（词牌名）作为标题。

---

### 3. 诗经 (诗经/)

**文件：**
- `shijing.json` - 诗经全文

**JSON 结构：**

```json
[
  {
    "title": "关雎",
    "chapter": "国风",
    "section": "周南",
    "content": [
      "关关雎鸠，在河之洲。窈窕淑女，君子好逑。",
      "参差荇菜，左右流之。窈窕淑女，寤寐求之。",
      "求之不得，寤寐思服。悠哉悠哉，辗转反侧。",
      "参差荇菜，左右采之。窈窕淑女，琴瑟友之。",
      "参差荇菜，左右芼之。窈窕淑女，钟鼓乐之。"
    ]
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 诗篇标题 |
| chapter | string | 是 | 章节（如：国风、小雅、大雅、颂） |
| section | string | 是 | 分节（如：周南、召南等） |
| content | string[] | 是 | 诗句内容（注：诗经使用 content 而非 paragraphs） |

---

### 4. 元曲 (元曲/)

**文件：**
- `yuanqu.json` - 元曲全集

**JSON 结构：**

```json
[
  {
    "dynasty": "yuan",
    "author": "关汉卿",
    "paragraphs": [
      "半世为人，不曾教大人心困。",
      "虽是搽胭粉，只争不裹头巾，将那等不做人的婆娘恨。"
    ],
    "title": "诈妮子调风月・仙吕/点绛唇"
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| dynasty | string | 是 | 朝代标识（"yuan"） |
| author | string | 是 | 作者名称 |
| paragraphs | string[] | 是 | 曲文内容 |
| title | string | 是 | 曲牌名/标题 |

---

### 5. 楚辞 (楚辞/)

**文件：**
- `chuci.json` - 楚辞全文

**JSON 结构：**

```json
[
  {
    "title": "离骚",
    "section": "离骚",
    "author": "屈原",
    "content": [
      "帝高阳之苗裔兮，朕皇考曰伯庸。",
      "摄提贞于孟陬兮，惟庚寅吾以降。",
      "皇览揆余初度兮，肇锡余以嘉名：",
      "名余曰正则兮，字余曰灵均。"
    ]
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 篇章标题 |
| section | string | 是 | 所属部分 |
| author | string | 是 | 作者名称 |
| content | string[] | 是 | 诗句内容 |

---

### 6. 五代诗词 (五代诗词/)

**子目录：**
- `huajianji/` - 花间集
- `nantang/` - 南唐诗词

**文件命名规则：**
- `huajianji-{卷数}-{类型}.json` - 花间集按卷分类

**JSON 结构（与唐诗宋词类似）：**

```json
[
  {
    "author": "韦庄",
    "paragraphs": ["..."],
    "rhythmic": "菩萨蛮",
    "title": "..."
  }
]
```

---

### 7. 纳兰性德 (纳兰性德/)

**文件：**
- `纳兰性德诗集.json` - 纳兰性德全集

**JSON 结构：**

```json
[
  {
    "title": "长相思·山一程",
    "para": [
      "山一程，水一程，身向榆关那畔行，夜深千帐灯",
      "风一更，雪一更，聒碎乡心梦不成，故园无此声"
    ],
    "author": "纳兰性德"
  }
]
```

**字段说明：**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 诗词标题 |
| para | string[] | 是 | 诗词内容（注：纳兰性德使用 para 而非 paragraphs） |
| author | string | 是 | 作者名称 |

**注意：** 纳兰性德数据使用 `para` 字段存储诗句。

---

### 8. 蒙学 (蒙学/)

**文件列表：**
- `sanzijing-new.json` - 三字经（新版）
- `sanzijing-traditional.json` - 三字经（繁体版）
- `qianziwen.json` - 千字文
- `baijiaxing.json` - 百家姓
- `dizigui.json` - 弟子规
- `tangshisanbaishou.json` - 唐诗三百首（蒙学版）
- `qianjiashi.json` - 千家诗
- `wenzimengqiu.json` - 文字蒙求
- `youxueqionglin.json` - 幼学琼林
- `zhuzijiaxun.json` - 朱子家训
- `zengguangxianwen.json` - 增广贤文
- `shenglvqimeng.json` - 声律启蒙
- `guwenguanzhi.json` - 古文观止

---

### 9. 四书五经 (四书五经/)

**文件列表：**
- `daxue.json` - 大学
- `zhongyong.json` - 中庸
- `mengzi.json` - 孟子

**JSON 结构（孟子示例）：**

```json
[
  {
    "chapter": "梁惠王上",
    "paragraphs": [
      "孟子见梁惠王。王曰：'叟不远千里而来，亦将有以利吾国乎？'",
      "孟子对曰：'王何必曰利？亦有仁义而已矣。'"
    ]
  }
]
```

---

### 10. 水墨唐诗 (水墨唐诗/)

**文件：**
- `shuimotangshi.json` - 水墨唐诗精选

---

### 11. 其他目录

| 目录 | 说明 |
|------|------|
| 曹操诗集/ | 曹操诗集 |
| 幽梦影/ | 幽梦影 |
| 御定全唐詩/ | 御定全唐诗（官方版本） |

---

## 字段映射总结

不同数据集使用不同的字段名存储诗句内容：

| 数据集 | 内容字段 | 备注 |
|--------|----------|------|
| 全唐诗 | paragraphs | 标准字段 |
| 宋词 | paragraphs | 标准字段 |
| 元曲 | paragraphs | 标准字段 |
| 诗经 | content | 需转换为 paragraphs |
| 楚辞 | content | 需转换为 paragraphs |
| 纳兰性德 | para | 需转换为 paragraphs |
| 五代诗词 | paragraphs | 标准字段 |

---

## 数据加载建议

### 1. 文件过滤规则

**全唐诗：**
- 仅加载 `poet.tang.*.json` 和 `poet.song.*.json` 开头的文件
- 忽略 `authors.*.json`、`唐诗三百首.json` 等非诗词数据文件
- 忽略 `error/` 子目录

**宋词：**
- 仅加载 `ci.song.*.json` 开头的文件
- 忽略 `author.song.json`、`宋词三百首.json`、`ci.db` 等文件

### 2. 字段标准化处理

```go
// 统一处理为 paragraphs 字段
if len(poem.Paragraphs) == 0 {
    if len(poem.Content) > 0 {
        poem.Paragraphs = poem.Content  // 诗经、楚辞
    } else if len(poem.Para) > 0 {
        poem.Paragraphs = poem.Para  // 纳兰性德
    }
}
```

### 3. 标题处理（宋词特殊处理）

```go
// 宋词没有 title 时使用 rhythmic（词牌名）作为标题
if poem.Title == "" && poem.Rhythmic != "" {
    poem.Title = poem.Rhythmic
}
```

### 4. 数据验证

必须满足以下条件才视为有效数据：
- `title` 不为空（或处理后不为空）
- `paragraphs` 不为空且长度 > 0

---

## 统计信息

| 分类 | 文件数（估算） | 数据量（估算） |
|------|----------------|----------------|
| 全唐诗 | 200+ JSON 文件 | 5万+ 首诗 |
| 宋词 | 30+ JSON 文件 | 2万+ 首词 |
| 诗经 | 1 个文件 | 305 篇 |
| 楚辞 | 1 个文件 | 数十篇 |
| 元曲 | 1 个文件 | 数千首 |
| 五代诗词 | 数个文件 | 数千首 |
| 纳兰性德 | 1 个文件 | 数百首 |

---

## 附录：完整字段清单

### 通用字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | string | 唯一标识符（仅部分数据有） |
| title | string | 标题 |
| author | string | 作者 |
| paragraphs | string[] | 诗句内容（标准字段） |
| content | string[] | 诗句内容（诗经、楚辞使用） |
| para | string[] | 诗句内容（纳兰性德使用） |

### 特殊字段

| 字段名 | 类型 | 适用范围 | 说明 |
|--------|------|----------|------|
| rhythmic | string | 宋词、五代诗词 | 词牌名 |
| dynasty | string | 元曲 | 朝代标识 |
| chapter | string | 诗经、四书五经 | 章节 |
| section | string | 诗经、楚辞 | 分节 |
| tags | string[] | 全唐诗 | 标签分类 |

### 元数据字段（程序添加）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| category_id | string | 分类标识（如：tang-shi、song-ci） |
| dynasty_id | string | 朝代标识（如：tang、song、yuan） |
| author_id | string | 作者标识 |
