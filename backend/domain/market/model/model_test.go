package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"thinkingModels/component/db"
)

// TestInsertSampleModels 插入5条思维模型测试数据
func TestInsertSampleModels(t *testing.T) {
	// 初始化数据库
	db.InitDb()

	// 设置测试上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set("currUserId", "1")
	ctx.Set("currUserName", "test")

	// 5条测试数据
	models := []struct {
		Name          string
		Code          string
		Description   string
		CategoryId    uint64
		Price         float64
		Content       string
		Difficulty    int
		EstimatedTime int
		Status        int
		AuthorName    string
	}{
		{
			Name:          "SWOT分析模型",
			Code:          "swot-analysis",
			Description:   "系统分析优势、劣势、机会与威胁的经典战略框架",
			CategoryId:    1,
			Price:         0,
			Content:       `{"fields":[{"key":"strengths","label":"优势","type":"textarea"},{"key":"weaknesses","label":"劣势","type":"textarea"},{"key":"opportunities","label":"机会","type":"textarea"},{"key":"threats","label":"威胁","type":"textarea"}]}`,
			Difficulty:    1,
			EstimatedTime: 30,
			Status:        1,
			AuthorName:    "思维大师",
		},
		{
			Name:          "5W1H分析法",
			Code:          "5w1h-analysis",
			Description:   "通过六个基本问题全面梳理事物要素的分析方法",
			CategoryId:    1,
			Price:         0,
			Content:       `{"fields":[{"key":"what","label":"What: 是什么？","type":"textarea"},{"key":"why","label":"Why: 为什么？","type":"textarea"},{"key":"who","label":"Who: 谁？","type":"textarea"},{"key":"when","label":"When: 何时？","type":"textarea"},{"key":"where","label":"Where: 何地？","type":"textarea"},{"key":"how","label":"How: 如何做？","type":"textarea"}]}`,
			Difficulty:    1,
			EstimatedTime: 25,
			Status:        1,
			AuthorName:    "系统思考者",
		},
		{
			Name:          "金字塔原理",
			Code:          "pyramid-principle",
			Description:   "结构化思考和表达的经典方法论",
			CategoryId:    2,
			Price:         29.9,
			Content:       `{"fields":[{"key":"conclusion","label":"核心结论","type":"text"},{"key":"arguments","label":"支撑论点","type":"array"},{"key":"evidence","label":"具体证据","type":"array"}]}`,
			Difficulty:    2,
			EstimatedTime: 45,
			Status:        1,
			AuthorName:    "麦肯锡专家",
		},
		{
			Name:          "决策矩阵",
			Code:          "decision-matrix",
			Description:   "多维度评估选项的量化决策工具",
			CategoryId:    2,
			Price:         19.9,
			Content:       `{"fields":[{"key":"options","label":"可选方案","type":"array"},{"key":"criteria","label":"评估标准","type":"array"},{"key":"scores","label":"评分矩阵","type":"matrix"},{"key":"weights","label":"权重分配","type":"array"}]}`,
			Difficulty:    2,
			EstimatedTime: 40,
			Status:        1,
			AuthorName:    "数据分析师",
		},
		{
			Name:          "第一性原理",
			Code:          "first-principles",
			Description:   "马斯克推崇的创新思维方法，回归本质思考问题",
			CategoryId:    3,
			Price:         49.9,
			Content:       `{"fields":[{"key":"assumptions","label":"既有假设","type":"textarea"},{"key":"breakdown","label":"拆解要素","type":"textarea"},{"key":"fundamentals","label":"基本原理","type":"textarea"},{"key":"rebuild","label":"重构方案","type":"textarea"}]}`,
			Difficulty:    3,
			EstimatedTime: 60,
			Status:        1,
			AuthorName:    "创新导师",
		},
	}

	// 插入数据
	for _, m := range models {
		// 检查是否已存在
		var count int64
		db.InitDb().Model(&ModelEntity{}).Where("code = ?", m.Code).Count(&count)
		if count > 0 {
			t.Logf("模型 %s 已存在，跳过", m.Code)
			continue
		}

		// 创建新实体
		newEntity := NewModelEntity(ctx)
		if concrete, ok := newEntity.(*ModelEntity); ok {
			concrete.Name = m.Name
			concrete.Code = m.Code
			concrete.Description = m.Description
			concrete.CategoryId = m.CategoryId
			concrete.Price = m.Price
			concrete.Content = m.Content
			concrete.Difficulty = m.Difficulty
			concrete.EstimatedTime = m.EstimatedTime
			concrete.Status = m.Status
			concrete.AuthorName = m.AuthorName
			concrete.AuthorId = 1
			concrete.Version = "1.0.0"
		}

		_, err := newEntity.Create()
		if err != nil {
			t.Errorf("创建模型 %s 失败: %v", m.Code, err)
		} else {
			t.Logf("成功创建模型: %s", m.Name)
		}
	}

	t.Log("测试数据插入完成")
}
