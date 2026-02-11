package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"thinkingModels/domain/user/user"
	"thinkingModels/logic"
)

// UserLogic 用户业务逻辑
type UserLogic struct {
	logic.BaseLogic
}

// 初始化UserLogic
func NewUserLogic(ctx *gin.Context) *UserLogic {
	return &UserLogic{BaseLogic: logic.BaseLogic{Ctx: ctx}}
}

// Create 创建用户
func (l *UserLogic) Create(req *user.RegisterRequest) (*user.UserInfo, error) {
	// 实例化模型
	userEntity := user.NewUserEntity(l.Ctx)

	// 数据赋值
	_, err := userEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = userEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 密码加密
	err = userEntity.HashPassword()
	if err != nil {
		return nil, err
	}

	// 保存数据
	res, err := userEntity.Create()
	if err != nil {
		return nil, err
	}

	// 返回用户信息DTO（直接构造，脱敏处理）
	return &user.UserInfo{
		ID:            res.Id,
		Username:      res.Username,
		Nickname:      res.Nickname,
		Email:         res.Email,
		Phone:         res.Phone,
		Avatar:        res.Avatar,
		Status:        res.Status,
		LastLoginTime: res.LastLoginTime.String(),
		CreatedAt:     res.CreatedAt.String(),
	}, nil
}

// Update 更新用户信息
func (l *UserLogic) Update(req *user.UpdateUserRequest) (*user.UserInfo, error) {
	// 实例化模型
	userEntity := user.NewUserEntity(l.Ctx)

	// 先加载旧数据
	_, err := userEntity.LoadById(req.ID)
	if err != nil {
		return nil, err
	}

	// 数据赋值
	_, err = userEntity.SetData(req)
	if err != nil {
		return nil, err
	}

	// 数据校验
	err = userEntity.Validate()
	if err != nil {
		return nil, err
	}

	// 更新数据
	res, err := userEntity.Update()
	if err != nil {
		return nil, err
	}

	// 返回用户信息DTO（直接构造，脱敏处理）
	return &user.UserInfo{
		ID:            res.Id,
		Username:      res.Username,
		Nickname:      res.Nickname,
		Email:         res.Email,
		Phone:         res.Phone,
		Avatar:        res.Avatar,
		Status:        res.Status,
		LastLoginTime: res.LastLoginTime.String(),
		CreatedAt:     res.CreatedAt.String(),
	}, nil
}

// Get 查询用户详情
func (l *UserLogic) Get(id uint64) (*user.UserInfo, error) {
	// 实例化模型
	userEntity := user.NewUserEntity(l.Ctx)

	// 查询数据
	res, err := userEntity.LoadById(id)
	if err != nil {
		return nil, err
	}

	// 返回用户信息DTO（直接构造，脱敏处理）
	return &user.UserInfo{
		ID:            res.Id,
		Username:      res.Username,
		Nickname:      res.Nickname,
		Email:         res.Email,
		Phone:         res.Phone,
		Avatar:        res.Avatar,
		Status:        res.Status,
		LastLoginTime: res.LastLoginTime.String(),
		CreatedAt:     res.CreatedAt.String(),
	}, nil
}

// List 查询用户列表
func (l *UserLogic) List(req *user.SearchUser) (*user.ListUserResponse, error) {
	// 实例化模型
	userEntity := user.NewUserEntity(l.Ctx)

	// 获取搜索条件
	cond := userEntity.MakeConditon(*req)

	// 查询总数
	total, err := userEntity.Count(cond)
	if err != nil {
		return nil, err
	}

	// 获取列表
	list, err := userEntity.List(cond)
	if err != nil {
		return nil, err
	}

	// 转换为DTO列表（直接构造，脱敏处理）
	userInfoList := make([]*user.UserInfo, 0, len(list))
	for _, item := range list {
		userInfoList = append(userInfoList, &user.UserInfo{
			ID:            item.Id,
			Username:      item.Username,
			Nickname:      item.Nickname,
			Email:         item.Email,
			Phone:         item.Phone,
			Avatar:        item.Avatar,
			Status:        item.Status,
			LastLoginTime: item.LastLoginTime.String(),
			CreatedAt:     item.CreatedAt.String(),
		})
	}

	// 返回数据
	res := &user.ListUserResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		List:     userInfoList,
	}
	return res, nil
}

// Del 删除用户
func (l *UserLogic) Del(ids []uint64) (any, error) {
	// 实例化模型
	userEntity := user.NewUserEntity(l.Ctx)

	// 删除数据
	err := userEntity.Del(ids...)
	return nil, err
}

// Login 用户登录
func (l *UserLogic) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	// 1. 根据用户名查询用户
	userEntity := user.NewUserEntity(l.Ctx)
	cond := userEntity.MakeConditon(user.SearchUser{Username: req.Username})
	dbUser, err := userEntity.LoadData(cond)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 检查用户状态
	if dbUser.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}

	// 3. 验证密码（实体方法）
	if !dbUser.VerifyPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 4. 更新登录信息（实体方法）
	dbUser.UpdateLoginInfo(l.Ctx.ClientIP())
	_, err = dbUser.Update()
	if err != nil {
		return nil, err
	}

	// 5. 生成JWT Token（实体方法 - 充血模型）
	tokenPair, err := dbUser.GenerateToken()
	if err != nil {
		return nil, errors.New("Token生成失败")
	}

	return &user.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		UserInfo: user.UserInfo{ // 直接构造，脱敏处理
			ID:            dbUser.Id,
			Username:      dbUser.Username,
			Nickname:      dbUser.Nickname,
			Email:         dbUser.Email,
			Phone:         dbUser.Phone,
			Avatar:        dbUser.Avatar,
			Status:        dbUser.Status,
			LastLoginTime: dbUser.LastLoginTime.String(),
			CreatedAt:     dbUser.CreatedAt.String(),
		},
	}, nil
}

// Logout 用户登出
func (l *UserLogic) Logout(token string) error {
	// TODO: 实现登出逻辑，将Token加入黑名单
	return nil
}
