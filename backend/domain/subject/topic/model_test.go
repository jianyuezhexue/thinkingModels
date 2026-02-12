package topic

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"thinkingModels/component/db"
)

// TestInsertSampleTopics 插入10条课题测试数据（我的课题）
func TestInsertSampleTopics(t *testing.T) {
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

	// 10条课题测试数据
	topics := []struct {
		Title       string
		Description string
		Status      int
		Priority    int
		Tags        string
		ModelId     uint64
		Deadline    time.Time
	}{
		{
			Title:       "产品战略规划分析",
			Description: "使用SWOT分析法对新产品进行全面的市场定位和竞争策略分析，明确产品的核心优势和差异化竞争点。",
			Status:      0,
			Priority:    3,
			Tags:        "战略,产品,SWOT",
			ModelId:     1,
			Deadline:    time.Now().AddDate(0, 0, 7),
		},
		{
			Title:       "用户增长瓶颈突破方案",
			Description: "针对当前用户增长放缓的问题，运用第一性原理思考，从本质出发重新设计用户获取和留存策略。",
			Status:      0,
			Priority:    2,
			Tags:        "增长,用户,第一性原理",
			ModelId:     5,
			Deadline:    time.Now().AddDate(0, 0, 14),
		},
		{
			Title:       "团队协作流程优化",
			Description: "分析现有跨部门协作中的痛点和阻塞点，设计更高效的沟通机制和协作流程。",
			Status:      0,
			Priority:    1,
			Tags:        "协作,流程优化,团队",
			ModelId:     2,
			Deadline:    time.Now().AddDate(0, 0, 21),
		},
		{
			Title:       "技术架构升级决策",
			Description: "评估现有技术架构的性能瓶颈，使用决策矩阵对比不同技术方案的可行性、成本和收益。",
			Status:      0,
			Priority:    3,
			Tags:        "技术,架构,决策矩阵",
			ModelId:     4,
			Deadline:    time.Now().AddDate(0, 1, 0),
		},
		{
			Title:       "Q4营销方案汇报材料准备",
			Description: "使用金字塔原理结构化整理Q4季度营销策略，确保汇报逻辑清晰、重点突出。",
			Status:      1,
			Priority:    2,
			Tags:        "汇报,金字塔原理,营销",
			ModelId:     3,
			Deadline:    time.Now().AddDate(0, 0, -5),
		},
		{
			Title:       "客户投诉处理机制改进",
			Description: "运用5W1H分析法全面梳理客户投诉处理流程，找出关键改进点，提升客户满意度。",
			Status:      0,
			Priority:    2,
			Tags:        "客服,5W1H,流程改进",
			ModelId:     2,
			Deadline:    time.Now().AddDate(0, 0, 10),
		},
		{
			Title:       "新市场进入策略研究",
			Description: "对东南亚市场进行深入调研，分析进入策略、竞争态势和本地化需求。",
			Status:      0,
			Priority:    3,
			Tags:        "市场,战略,国际化",
			ModelId:     1,
			Deadline:    time.Now().AddDate(0, 2, 0),
		},
		{
			Title:       "年度绩效评估体系优化",
			Description: "重新设计绩效考核指标体系，使其更加公平、透明，并与公司战略目标对齐。",
			Status:      2,
			Priority:    1,
			Tags:        "绩效,HR,体系设计",
			ModelId:     4,
			Deadline:    time.Now().AddDate(0, -10, 0),
		},
		{
			Title:       "产品定价策略调整",
			Description: "分析竞品定价、成本结构和用户价格敏感度，制定新的产品定价方案。",
			Status:      1,
			Priority:    2,
			Tags:        "定价,商业,策略",
			ModelId:     5,
			Deadline:    time.Now().AddDate(0, 0, -2),
		},
		{
			Title:       "供应链风险管理体系建设",
			Description: "识别供应链中的关键风险点，建立风险预警机制和应急响应预案。",
			Status:      0,
			Priority:    3,
			Tags:        "供应链,风险,管理",
			ModelId:     2,
			Deadline:    time.Now().AddDate(0, 1, 15),
		},
	}

	// 插入数据
	for _, tp := range topics {
		// 检查是否已存在（根据标题判断）
		var count int64
		db.InitDb().Model(&TopicEntity{}).Where("title = ?", tp.Title).Count(&count)
		if count > 0 {
			t.Logf("课题 '%s' 已存在，跳过", tp.Title)
			continue
		}

		// 创建新实体
		entity := NewTopicEntity(ctx)
		if concrete, ok := entity.(*TopicEntity); ok {
			concrete.Title = tp.Title
			concrete.Description = tp.Description
			concrete.Status = tp.Status
			concrete.Priority = tp.Priority
			concrete.Tags = tp.Tags
			concrete.ModelId = tp.ModelId
			concrete.UserId = 1 // 当前用户ID
			concrete.Deadline = db.LocalTime(tp.Deadline)
		}

		_, err := entity.Create()
		if err != nil {
			t.Errorf("创建课题 '%s' 失败: %v", tp.Title, err)
		} else {
			t.Logf("成功创建课题: %s", tp.Title)
		}
	}

	t.Log("10条课题测试数据插入完成")
}
