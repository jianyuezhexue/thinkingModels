package category

// 新增实体
type CreateCategory struct {
	Name        string `json:"name" uri:"name" binding:"required"`     // 分类名称
	Icon        string `json:"icon" uri:"icon"`                      // 分类图标URL
	Description string `json:"description" uri:"description"`        // 分类描述
	Heat        int    `json:"heat" uri:"heat"`                      // 热度值
}

// 更新实体
type UpdateCategory struct {
	Id          uint64 `json:"id" binding:"required"`                // 主键ID
	Name        string `json:"name" uri:"name"`                      // 分类名称
	Icon        string `json:"icon" uri:"icon"`                      // 分类图标URL
	Description string `json:"description" uri:"description"`        // 分类描述
	Heat        int    `json:"heat" uri:"heat"`                      // 热度值
}

// 搜索实体
type SearchCategory struct {
	Page        int64  `json:"page" uri:"page" search:"page"`                                                                // 分页
	PageSize    int64  `json:"pageSize" uri:"pageSize" search:"pageSize"`                                                    // 分页大小
	Id          uint64 `json:"id" uri:"id" search:"type:eq;column:id;table:category"`                                          // ID
	Ids         []uint64 `json:"ids" uri:"ids" search:"type:in;column:id;table:category"`                                      // IDs
	CreatedAt   []string `json:"createdAt" uri:"createdAt" search:"type:between;column:created_at;table:category"`             // 创建时间
	UpdatedAt   []string `json:"updatedAt" uri:"updatedAt" search:"type:between;column:updated_at;table:category"`             // 创建时间
	Name        string   `json:"name" uri:"name" search:"type:like;column:name;table:category"`                                // 分类名称（模糊查询）
	Heat        int      `json:"heat" uri:"heat" search:"type:eq;column:heat;table:category"`                                  // 热度值
}

// 删除实体
type DelCategory struct {
	Ids []uint64 `json:"ids" binding:"required"`
}

// 分页返回
type ListReap struct {
	Page     int64             `json:"page" comment:"页数"`
	PageSize int64             `json:"pageSize" comment:"每页数量"`
	Total    int64             `json:"total" comment:"总条数"`
	List     []*CategoryEntity `json:"list" comment:"数据"`
}

// 全量列表返回（用于筛选条件）
type AllReap struct {
	List []*CategoryEntity `json:"list" comment:"全部分类数据"`
}

// 增加热度请求
type IncreaseHeatRequest struct {
	Id    uint64 `json:"id" binding:"required"`    // 分类ID
	Delta int    `json:"delta" binding:"required,min=1"` // 增加的热度值，必须大于0
}
