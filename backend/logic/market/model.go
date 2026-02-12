package market

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/market/model"
	"thinkingModels/logic"
)

// ModelLogic 思维模型业务逻辑
type ModelLogic struct {
	logic.BaseLogic
}

// NewModelLogic 初始化ModelLogic
func NewModelLogic(ctx *gin.Context) *ModelLogic {
	return &ModelLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建思维模型
func (l *ModelLogic) Create(req *model.CreateModel) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 数据赋值
	_, err := entity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	if concreteEntity, ok := entity.(*model.ModelEntity); ok {
		concreteEntity.Status = 0 // 默认为草稿状态
		concreteEntity.Version = "1.0.0"
		if concreteEntity.Difficulty == 0 {
			concreteEntity.Difficulty = 1
		}
		if concreteEntity.EstimatedTime == 0 {
			concreteEntity.EstimatedTime = 30
		}
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 保存数据
	res, err := entity.Create()
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToModelInfo(res), nil
}

// Update 更新思维模型
func (l *ModelLogic) Update(req *model.UpdateModel) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 先加载旧数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 数据赋值
	_, err = entity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	if err := entity.Validate(); err != nil {
		return nil, err
	}

	// 更新数据
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToModelInfo(res), nil
}

// Get 查询思维模型详情
func (l *ModelLogic) Get(id uint64) (*model.ModelDetail, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 查询数据
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	// 转换为详情DTO
	return convertToModelDetail(res), nil
}

// GetSimple 查询思维模型简单信息（不含内容）
func (l *ModelLogic) GetSimple(id uint64) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 查询数据
	res, err := entity.LoadById(id)
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToModelInfo(res), nil
}

// List 查询思维模型列表
func (l *ModelLogic) List(req *model.SearchModel) (*model.ListModelResponse, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 处理前端传入的category字符串，映射为CategoryId
	if req.Category != "" && req.CategoryId == 0 {
		req.CategoryId = categoryStringToId(req.Category)
	}

	// 处理isFree参数，映射为价格条件
	if req.IsFree != nil {
		if *req.IsFree {
			// 免费：价格为0
			req.MinPrice = 0
			req.MaxPrice = 0
		}
		// 付费的情况可以在查询后过滤，或者设置MinPrice > 0
	}

	// 获取搜索条件
	cond := entity.MakeConditon(*req)

	// 查询总数
	total, err := entity.Count(cond)
	if err != nil {
		return nil, err
	}

	// 获取列表
	list, err := entity.List(cond)
	if err != nil {
		return nil, err
	}

	// 过滤付费模型（当isFree=false时）
	if req.IsFree != nil && !*req.IsFree {
		filteredList := make([]*model.ModelEntity, 0)
		for _, item := range list {
			if item.Price > 0 {
				filteredList = append(filteredList, item)
			}
		}
		list = filteredList
		total = int64(len(list))
	}

	// 转换为DTO列表
	modelInfoList := make([]*model.ModelInfo, 0, len(list))
	for _, item := range list {
		modelInfoList = append(modelInfoList, convertToModelInfo(item))
	}

	// 返回数据
	res := &model.ListModelResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     modelInfoList,
	}
	return res, nil
}

// categoryStringToId 将分类字符串映射为分类ID
func categoryStringToId(category string) uint64 {
	categoryMap := map[string]uint64{
		"business":  1,
		"strategy":  2,
		"innovation": 3,
		"analysis":  4,
		"decision":  5,
		"creative":  6,
	}
	if id, ok := categoryMap[category]; ok {
		return id
	}
	return 0 // 返回0表示不筛选分类
}

// Del 删除思维模型
func (l *ModelLogic) Del(ids []uint64) (any, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 删除数据
	err := entity.Del(ids...)
	return nil, err
}

// Publish 发布思维模型
func (l *ModelLogic) Publish(req *model.PublishModel) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 加载数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 执行发布
	if err := entity.Publish(); err != nil {
		return nil, err
	}

	// 保存更新
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// Unpublish 下架思维模型
func (l *ModelLogic) Unpublish(req *model.UnpublishModel) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 加载数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 执行下架
	if err := entity.Unpublish(); err != nil {
		return nil, err
	}

	// 保存更新
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// UpdateStatus 更新思维模型状态
func (l *ModelLogic) UpdateStatus(req *model.UpdateModelStatus) (*model.ModelInfo, error) {
	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 加载数据
	_, err := entity.LoadById(req.Id)
	if err != nil {
		return nil, err
	}

	// 更新状态
	if concreteEntity, ok := entity.(*model.ModelEntity); ok {
		concreteEntity.Status = req.Status
		// 如果发布，设置发布时间
		if req.Status == 1 && concreteEntity.PublishTime.IsZero() {
			concreteEntity.Publish()
		}
	}

	// 保存更新
	res, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return convertToModelInfo(res), nil
}

// GetByCode 根据编码查询思维模型
func (l *ModelLogic) GetByCode(code string) (*model.ModelInfo, error) {
	if code == "" {
		return nil, errors.New("模型编码不能为空")
	}

	// 实例化模型
	entity := model.NewModelEntity(l.Ctx)

	// 构建查询条件
	cond := entity.MakeConditon(model.SearchModel{Code: code})

	// 查询数据
	res, err := entity.LoadData(cond)
	if err != nil {
		return nil, err
	}

	// 转换为DTO返回
	return convertToModelInfo(res), nil
}

// 辅助函数：转换为ModelInfo
func convertToModelInfo(entity *model.ModelEntity) *model.ModelInfo {
	if entity == nil {
		return nil
	}

	// 将uint64 ID转换为string
	id := strconv.FormatUint(entity.Id, 10)
	categoryId := strconv.FormatUint(entity.CategoryId, 10)
	authorId := strconv.FormatUint(entity.AuthorId, 10)

	// 获取分类名称（从CategoryId映射，实际需要查询分类表）
	categoryName := getCategoryName(entity.CategoryId)

	// 构造返回值 - 匹配前端期望的数据结构
	return &model.ModelInfo{
		Id:          id,
		Title:       entity.Name,
		Description: entity.Description,
		Cover:       entity.CoverImage,
		Author: model.ModelAuthor{
			Id:     authorId,
			Name:   entity.AuthorName,
			Avatar: "", // 暂时为空，实际需要查询用户表获取
		},
		IsFree:     entity.Price == 0,
		Price:      entity.Price,
		Category:   categoryName,
		CategoryId: categoryId,
		Tags:       []string{}, // 暂时为空，需要单独查询标签
		Stats: model.ModelStats{
			Adoptions:   entity.UsageCount, // 使用usageCount作为adoptions
			Practices:   entity.UsageCount, // 临时使用相同值
			Discussions: 0,                 // 需要单独统计
			Forks:       0,                 // 需要单独统计
			Likes:       entity.UsageCount, // 临时使用相同值
		},
		UpdatedAt: entity.UpdatedAt.String(),
		Status:    entity.Status,
	}
}

// 辅助函数：获取分类名称（简单实现，实际需要查询数据库）
func getCategoryName(categoryId uint64) string {
	// 简单的映射，实际应该从数据库查询
	categoryMap := map[uint64]string{
		1: "商业管理",
		2: "战略规划",
		3: "创新思维",
		4: "分析工具",
		5: "决策方法",
		6: "创意构思",
	}
	if name, ok := categoryMap[categoryId]; ok {
		return name
	}
	return "其他"
}

// 辅助函数：转换为ModelDetail
func convertToModelDetail(entity *model.ModelEntity) *model.ModelDetail {
	if entity == nil {
		return nil
	}

	info := convertToModelInfo(entity)
	// 填充content内容
	info.Content = entity.Content
	return &model.ModelDetail{
		ModelInfo: *info,
	}
}
