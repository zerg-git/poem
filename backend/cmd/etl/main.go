package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"poem/backend/models"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RawPoem 原始 JSON 结构
type RawPoem struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Content    []string `json:"content"`  // 备用
	Para       []string `json:"para"`     // 备用
	Rhythmic   string   `json:"rhythmic"` // 宋词
	Dynasty    string   `json:"dynasty"`  // 元曲
	Chapter    string   `json:"chapter"`
	Section    string   `json:"section"`
	Volume     string   `json:"volume"`
	Prologue   string   `json:"prologue"`
	Notes      []string `json:"notes"` // 注释
}

// RawAuthor 原始作者 JSON 结构
type RawAuthor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_description"`
}

// RawSiShuWuJing 四书五经 JSON 结构
type RawSiShuWuJing struct {
	Chapter    string   `json:"chapter"`
	Paragraphs []string `json:"paragraphs"`
}

// MultiStringSlice handles JSON fields that can be either a string or an array of strings
type MultiStringSlice []string

// UnmarshalJSON implements json.Unmarshaler
func (s *MultiStringSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}
	var multi []string
	if err := json.Unmarshal(data, &multi); err == nil {
		*s = multi
		return nil
	}
	return fmt.Errorf("comment should be string or []string")
}

// RawYouMengYing 幽梦影 JSON 结构
type RawYouMengYing struct {
	Content string           `json:"content"`
	Comment MultiStringSlice `json:"comment"`
}

// Cache for IDs
var (
	authorCache = make(map[string]uint)
	catCache    = make(map[string]uint)
	cacheMutex  sync.RWMutex
)

func main() {
	// 1. 初始化数据库
	dbPath := "../../../poems.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 迁移模式 - 重新创建表
	db.Migrator().DropTable(&models.Work{}, &models.Author{}, &models.Category{}, &models.Comment{})
	err = db.AutoMigrate(&models.Category{}, &models.Author{}, &models.Work{}, &models.Comment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 2. 确定数据目录
	possiblePaths := []string{
		"d:\\demo\\poem\\chinese-poetry",
		"d:\\demo\\poem\\chinese-poetry-master",
		"../../../chinese-poetry",
		"../../../chinese-poetry-master",
		"chinese-poetry",
		"chinese-poetry-master",
	}

	rootDir := ""
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			rootDir = p
			break
		}
	}

	if rootDir == "" {
		log.Fatal("Could not find chinese-poetry data directory")
	}

	fmt.Printf("Data root: %s\n", rootDir)

	// 3. 种子分类数据
	seedCategories(db)

	// 4. 处理全唐诗
	processDir(db, filepath.Join(rootDir, "全唐诗"), "quantangshi", "唐", func(filename string) bool {
		return strings.HasPrefix(filename, "poet.tang.")
	})
	processAuthors(db, filepath.Join(rootDir, "全唐诗", "authors.tang.json"), "唐")
	processPoemFile(db, filepath.Join(rootDir, "全唐诗", "唐诗三百首.json"), "quantangshi", "唐")

	// 5. 处理宋词
	processDir(db, filepath.Join(rootDir, "宋词"), "songci", "宋", func(filename string) bool {
		return strings.HasPrefix(filename, "ci.song.")
	})
	processAuthors(db, filepath.Join(rootDir, "宋词", "author.song.json"), "宋")
	processPoemFile(db, filepath.Join(rootDir, "宋词", "宋词三百首.json"), "songci", "宋")

	// 6. 处理元曲
	processPoemFile(db, filepath.Join(rootDir, "元曲", "yuanqu.json"), "yuanqu", "元")

	// 7. 处理四书五经
	siShuDir := filepath.Join(rootDir, "四书五经")
	processSiShuWuJing(db, filepath.Join(siShuDir, "daxue.json"), "曾子", "sishuwujing")
	processSiShuWuJing(db, filepath.Join(siShuDir, "mengzi.json"), "孟子", "sishuwujing")
	processSiShuWuJing(db, filepath.Join(siShuDir, "zhongyong.json"), "子思", "sishuwujing")
	if _, err := os.Stat(filepath.Join(rootDir, "论语", "lunyu.json")); err == nil {
		processSiShuWuJing(db, filepath.Join(rootDir, "论语", "lunyu.json"), "孔子", "lunyu")
	}

	// 8. 处理幽梦影
	processYouMengYing(db, filepath.Join(rootDir, "幽梦影", "youmengying.json"), "youmengying")

	// 9. 处理诗经
	processShiJing(db, filepath.Join(rootDir, "诗经", "shijing.json"), "shijing")

	// 10. 处理楚辞
	processChuCi(db, filepath.Join(rootDir, "楚辞", "chuci.json"), "chuci")

	// 11. 处理曹操诗集
	processPoemFile(db, filepath.Join(rootDir, "曹操诗集", "caocao.json"), "caocao", "汉/魏")

	// 12. 处理纳兰性德
	processPoemFile(db, filepath.Join(rootDir, "纳兰性德", "纳兰性德诗集.json"), "nalan", "清")

	// 13. 处理水墨唐诗
	processPoemFile(db, filepath.Join(rootDir, "水墨唐诗", "shuimotangshi.json"), "shuimotangshi", "唐")

	// 14. 处理五代诗词
	processDir(db, filepath.Join(rootDir, "五代诗词", "huajianji"), "wudai", "五代", func(filename string) bool {
		return strings.HasSuffix(filename, ".json") && filename != "README.md"
	})
	processPoemFile(db, filepath.Join(rootDir, "五代诗词", "nantang", "poetrys.json"), "wudai", "五代")

	fmt.Println("Done!")
}

func seedCategories(db *gorm.DB) {
	categories := []models.Category{
		{Name: "quantangshi", DisplayName: "全唐诗", Description: "全唐诗收录唐诗四万八千九百余首"},
		{Name: "songci", DisplayName: "宋词", Description: "全宋词收录宋词二万余首"},
		{Name: "yuanqu", DisplayName: "元曲", Description: "元代文学形式，包括散曲和杂剧"},
		{Name: "shijing", DisplayName: "诗经", Description: "中国古代诗歌开端"},
		{Name: "chuci", DisplayName: "楚辞", Description: "屈原创作的诗歌总集"},
		{Name: "lunyu", DisplayName: "论语", Description: "儒家经典"},
		{Name: "sishuwujing", DisplayName: "四书五经", Description: "儒家经典著作"},
		{Name: "youmengying", DisplayName: "幽梦影", Description: "清代张潮著"},
		{Name: "caocao", DisplayName: "曹操诗集", Description: "曹操诗歌全集"},
		{Name: "nalan", DisplayName: "纳兰性德", Description: "清代词人纳兰性德诗集"},
		{Name: "shuimotangshi", DisplayName: "水墨唐诗", Description: "水墨风格唐诗精选"},
		{Name: "wudai", DisplayName: "五代诗词", Description: "五代十国时期的诗词作品"},
		{Name: "mengxue", DisplayName: "蒙学", Description: "古代启蒙教材"},
	}

	for _, c := range categories {
		db.FirstOrCreate(&c, models.Category{Name: c.Name})
		cacheMutex.Lock()
		catCache[c.Name] = c.ID
		cacheMutex.Unlock()
	}
}

func getCategoryID(name string) uint {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return catCache[name]
}

func getOrCreateAuthor(db *gorm.DB, name string, dynasty string) uint {
	key := name + "|" + dynasty
	cacheMutex.RLock()
	if id, ok := authorCache[key]; ok {
		cacheMutex.RUnlock()
		return id
	}
	cacheMutex.RUnlock()

	var author models.Author
	// Use db (which could be a transaction) to query
	err := db.Where("name = ? AND dynasty = ?", name, dynasty).First(&author).Error
	if err != nil {
		author = models.Author{Name: name, Dynasty: dynasty}
		db.Create(&author)
	}

	cacheMutex.Lock()
	authorCache[key] = author.ID
	cacheMutex.Unlock()
	return author.ID
}

func processAuthors(db *gorm.DB, filePath string, dynasty string) {
	fmt.Printf("Processing authors %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	var rawAuthors []RawAuthor
	if err := json.Unmarshal(content, &rawAuthors); err != nil {
		log.Printf("Error unmarshal file %s: %v", filePath, err)
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, ra := range rawAuthors {
			desc := ra.Description
			if desc == "" {
				desc = ra.Desc
			}

			var author models.Author
			err := tx.Where("name = ? AND dynasty = ?", ra.Name, dynasty).First(&author).Error
			if err == nil {
				// Update existing author with bio
				if author.Biography == "" && desc != "" {
					author.Biography = desc
					tx.Save(&author)
				}
			} else {
				// Create new
				author = models.Author{
					Name:      ra.Name,
					Dynasty:   dynasty,
					Biography: desc,
				}
				tx.Create(&author)
			}

			key := ra.Name + "|" + dynasty
			cacheMutex.Lock()
			authorCache[key] = author.ID
			cacheMutex.Unlock()
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d authors)\n", len(rawAuthors))
	}
}

func processDir(db *gorm.DB, dirPath string, categoryName string, defaultDynasty string, filter func(string) bool) {
	fmt.Printf("Processing dir %s...\n", dirPath)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error reading dir %s: %v", dirPath, err)
		return
	}

	catID := getCategoryID(categoryName)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		if filter != nil && !filter(file.Name()) {
			continue
		}

		processFile(db, filepath.Join(dirPath, file.Name()), catID, defaultDynasty)
	}
}

func normalizeDynasty(d string) string {
	d = strings.TrimSpace(d)
	switch strings.ToLower(d) {
	case "tang":
		return "唐"
	case "song":
		return "宋"
	case "yuan":
		return "元"
	case "ming":
		return "明"
	case "qing":
		return "清"
	}
	return d
}

func processFile(db *gorm.DB, filePath string, catID uint, defaultDynasty string) {
	fmt.Printf("Processing file %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	var rawPoems []RawPoem
	if err := json.Unmarshal(content, &rawPoems); err != nil {
		log.Printf("Error unmarshal file %s: %v", filePath, err)
		return
	}

	// Use Transaction for bulk insert performance
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, rp := range rawPoems {
			// Field normalization
			paragraphs := rp.Paragraphs
			if len(paragraphs) == 0 && len(rp.Content) > 0 {
				paragraphs = rp.Content
			}
			if len(paragraphs) == 0 && len(rp.Para) > 0 {
				paragraphs = rp.Para
			}

			if len(paragraphs) == 0 {
				continue
			}

			title := rp.Title
			if title == "" {
				title = rp.Rhythmic
			}

			dynasty := rp.Dynasty
			if dynasty == "" {
				dynasty = defaultDynasty
			} else {
				dynasty = normalizeDynasty(dynasty)
			}

			authorName := rp.Author
			if authorName == "" {
				// Special case for Cao Cao if category is caocao
				if catID == getCategoryID("caocao") {
					authorName = "曹操"
				} else {
					authorName = "Unknown"
				}
			}

			// Use tx here to ensure we are inside the transaction
			authorID := getOrCreateAuthor(tx, authorName, dynasty)

			work := models.Work{
				CategoryID: catID,
				AuthorID:   authorID,
				Title:      title,
				Rhythmic:   rp.Rhythmic,
				Content:    models.JSONArr(paragraphs),
				OriginalID: rp.ID,
				Volume:     rp.Volume,
				Section:    rp.Section,
				Prologue:   rp.Prologue,
			}

			if err := tx.Create(&work).Error; err != nil {
				continue
			}

			// Handle comments
			if len(rp.Notes) > 0 {
				var comments []models.Comment
				for _, note := range rp.Notes {
					comments = append(comments, models.Comment{
						WorkID:  work.ID,
						Content: note,
						Type:    "note",
					})
				}
				// Batch insert comments
				if len(comments) > 0 {
					tx.CreateInBatches(comments, 100)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d works)\n", len(rawPoems))
	}
}

func processPoemFile(db *gorm.DB, filePath string, categoryName string, defaultDynasty string) {
	catID := getCategoryID(categoryName)
	processFile(db, filePath, catID, defaultDynasty)
}

func processSiShuWuJing(db *gorm.DB, filePath string, defaultAuthor string, categoryName string) {
	fmt.Printf("Processing %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	// Try unmarshal as array first
	var rawData []RawSiShuWuJing
	if err := json.Unmarshal(content, &rawData); err != nil {
		// If failed, try unmarshal as single object
		var singleObj RawSiShuWuJing
		if err2 := json.Unmarshal(content, &singleObj); err2 == nil {
			rawData = []RawSiShuWuJing{singleObj}
		} else {
			log.Printf("Error unmarshal file %s: %v", filePath, err)
			return
		}
	}

	catID := getCategoryID(categoryName)

	err = db.Transaction(func(tx *gorm.DB) error {
		authorID := getOrCreateAuthor(tx, defaultAuthor, "先秦")

		for _, d := range rawData {
			work := models.Work{
				CategoryID: catID,
				AuthorID:   authorID,
				Title:      d.Chapter,
				Content:    models.JSONArr(d.Paragraphs),
			}
			tx.Create(&work)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d works)\n", len(rawData))
	}
}

func processYouMengYing(db *gorm.DB, filePath string, categoryName string) {
	fmt.Printf("Processing %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	var rawData []RawYouMengYing
	if err := json.Unmarshal(content, &rawData); err != nil {
		log.Printf("Error unmarshal file %s: %v", filePath, err)
		return
	}

	catID := getCategoryID(categoryName)

	err = db.Transaction(func(tx *gorm.DB) error {
		authorID := getOrCreateAuthor(tx, "张潮", "清")

		for i, d := range rawData {
			work := models.Work{
				CategoryID: catID,
				AuthorID:   authorID,
				Title:      fmt.Sprintf("幽梦影-%d", i+1),
				Content:    models.JSONArr([]string{d.Content}),
			}
			tx.Create(&work)

			if len(d.Comment) > 0 {
				var comments []models.Comment
				for _, c := range d.Comment {
					parts := strings.SplitN(c, "曰：", 2)
					commenter := ""
					noteContent := c
					if len(parts) == 2 {
						commenter = parts[0]
						noteContent = parts[1]
					}

					comments = append(comments, models.Comment{
						WorkID:    work.ID,
						Content:   noteContent,
						Commenter: commenter,
						Type:      "comment",
					})
				}
				if len(comments) > 0 {
					tx.CreateInBatches(comments, 100)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d items)\n", len(rawData))
	}
}

func processShiJing(db *gorm.DB, filePath string, categoryName string) {
	fmt.Printf("Processing %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	var rawPoems []RawPoem
	if err := json.Unmarshal(content, &rawPoems); err != nil {
		log.Printf("Error unmarshal file %s: %v", filePath, err)
		return
	}

	catID := getCategoryID(categoryName)

	err = db.Transaction(func(tx *gorm.DB) error {
		// 诗经作者通常认为是佚名，或具体篇目有归属，这里简化处理
		authorID := getOrCreateAuthor(tx, "佚名", "先秦")

		for _, rp := range rawPoems {
			work := models.Work{
				CategoryID: catID,
				AuthorID:   authorID,
				Title:      rp.Title,
				Volume:     rp.Chapter,                 // 国风
				Section:    rp.Section,                 // 周南
				Content:    models.JSONArr(rp.Content), // Shijing uses 'content'
			}
			tx.Create(&work)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d works)\n", len(rawPoems))
	}
}

func processChuCi(db *gorm.DB, filePath string, categoryName string) {
	fmt.Printf("Processing %s... ", filepath.Base(filePath))
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return
	}

	var rawPoems []RawPoem
	if err := json.Unmarshal(content, &rawPoems); err != nil {
		log.Printf("Error unmarshal file %s: %v", filePath, err)
		return
	}

	catID := getCategoryID(categoryName)

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, rp := range rawPoems {
			authorName := rp.Author
			if authorName == "" {
				authorName = "屈原" // Default to Qu Yuan if missing
			}
			authorID := getOrCreateAuthor(tx, authorName, "先秦")

			work := models.Work{
				CategoryID: catID,
				AuthorID:   authorID,
				Title:      rp.Title,
				Section:    rp.Section, // 离骚
				Content:    models.JSONArr(rp.Content),
			}
			tx.Create(&work)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Done (%d works)\n", len(rawPoems))
	}
}
