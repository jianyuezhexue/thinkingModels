package user

// ==================== 请求DTO ====================

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
	Captcha  string `json:"captcha"`                     // 验证码（可选）
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`  // 用户名
	Password string `json:"password" binding:"required,min=8"`         // 密码
	Nickname string `json:"nickname" binding:"max=50"`                 // 昵称
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`  // 手机号
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	Nickname string `json:"nickname" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11,numeric"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ==================== 响应DTO ====================

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`  // Access Token
	RefreshToken string   `json:"refreshToken"` // Refresh Token
	ExpiresIn    int64    `json:"expiresIn"`    // 过期时间（秒）
	UserInfo     UserInfo `json:"userInfo"`     // 用户信息
}

// UserInfo 用户信息DTO（脱敏，不含密码）
type UserInfo struct {
	ID            uint64 `json:"id"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	LastLoginTime string `json:"lastLoginTime"`
	CreatedAt     string `json:"createdAt"`
}

// ==================== 搜索条件 ====================

// SearchUser 用户搜索条件
type SearchUser struct {
	Page       int64  `json:"page" form:"page" search:"page"`                               // 分页
	PageSize   int64  `json:"pageSize" form:"pageSize" search:"pageSize"`                   // 分页大小
	ID         uint64 `json:"id" form:"id" search:"type:eq;column:id;table:users"`          // ID精确匹配
	Username   string `json:"username" form:"username" search:"type:like;column:username;table:users"` // 用户名模糊查询
	Status     int    `json:"status" form:"status" search:"type:eq;column:status;table:users"`         // 状态
	Email      string `json:"email" form:"email" search:"type:eq;column:email;table:users"`            // 邮箱
	Phone      string `json:"phone" form:"phone" search:"type:eq;column:phone;table:users"`            // 手机号
	CreatedAt  []string `json:"createdAt" form:"createdAt" search:"type:between;column:created_at;table:users"` // 创建时间范围
}

// ListUserResponse 用户列表响应
type ListUserResponse struct {
	Page     int64      `json:"page" comment:"页数"`
	PageSize int64      `json:"pageSize" comment:"每页数量"`
	Total    int64      `json:"total" comment:"总条数"`
	List     []*UserInfo `json:"list" comment:"数据"`
}

// DelUser 删除用户请求
type DelUser struct {
	Ids []uint64 `json:"ids" binding:"required"`
}
