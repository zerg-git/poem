package handlers

import (
	"net/http"
	"poem/backend/models"
	"poem/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PoetryHandler 诗词处理器
type PoetryHandler struct {
	service *services.PoetryService
}

// NewPoetryHandler 创建诗词处理器
func NewPoetryHandler(service *services.PoetryService) *PoetryHandler {
	return &PoetryHandler{service: service}
}

// GetPoems 获取诗词列表
// @Summary 获取诗词列表
// @Tags 诗词
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Success 200 {object} models.APIResponse
// @Router /poems [get]
func (h *PoetryHandler) GetPoems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	result, err := h.service.GetPoems(page, pageSize, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetPoemByID 获取单首诗词
// @Summary 获取单首诗词
// @Tags 诗词
// @Accept json
// @Produce json
// @Param id path string true "诗词ID"
// @Success 200 {object} models.APIResponse
// @Router /poems/{id} [get]
func (h *PoetryHandler) GetPoemByID(c *gin.Context) {
	id := c.Param("id")

	poem, err := h.service.GetPoemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "诗词不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    poem,
	})
}

// GetRandomPoem 获取随机诗词
// @Summary 获取随机诗词
// @Tags 诗词
// @Accept json
// @Produce json
// @Param count query int false "数量" default(1)
// @Param category query string false "分类"
// @Success 200 {object} models.APIResponse
// @Router /poems/random [get]
func (h *PoetryHandler) GetRandomPoem(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "1"))
	category := c.Query("category")

	poems, err := h.service.GetRandomPoems(count, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]interface{}{"poems": poems},
	})
}

// Search 搜索诗词
// @Summary 搜索诗词
// @Tags 搜索
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse
// @Router /search [get]
func (h *PoetryHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "请输入搜索关键词",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.service.Search(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetDynasties 获取朝代列表
// @Summary 获取朝代列表
// @Tags 目录
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /dynasties [get]
func (h *PoetryHandler) GetDynasties(c *gin.Context) {
	dynasties := h.service.GetDynasties()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    dynasties,
	})
}

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Tags 目录
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /categories [get]
func (h *PoetryHandler) GetCategories(c *gin.Context) {
	categories := h.service.GetCategories()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    categories,
	})
}

// GetAuthors 获取作者列表
// @Summary 获取作者列表
// @Tags 作者
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param dynasty query string false "朝代"
// @Success 200 {object} models.APIResponse
// @Router /authors [get]
func (h *PoetryHandler) GetAuthors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	dynasty := c.Query("dynasty")

	result, err := h.service.GetAuthors(page, pageSize, dynasty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetAuthorByName 获取作者详情
// @Summary 获取作者详情
// @Tags 作者
// @Accept json
// @Produce json
// @Param name path string true "作者名称"
// @Success 200 {object} models.APIResponse
// @Router /authors/{name} [get]
func (h *PoetryHandler) GetAuthorByName(c *gin.Context) {
	name := c.Param("name")

	author, err := h.service.GetAuthorByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "作者不存在",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    author,
	})
}

// GetAuthorPoems 获取作者的诗词
// @Summary 获取作者的诗词
// @Tags 作者
// @Accept json
// @Produce json
// @Param name path string true "作者名称"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse
// @Router /authors/{name}/poems [get]
func (h *PoetryHandler) GetAuthorPoems(c *gin.Context) {
	name := c.Param("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.service.GetPoemsByAuthor(name, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}
