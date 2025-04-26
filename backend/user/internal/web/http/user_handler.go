package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/service"
)

// UserHandler HTTP请求处理器
type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

// NewUserHandler 创建用户HTTP处理器
func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// RegisterRoutes 注册HTTP路由
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	// 公开API，不需要认证
	public := router.Group("/api/v1")
	{
		public.POST("/users/register", h.Register)
		public.POST("/users/login", h.Login)
		public.POST("/users/refresh-token", h.RefreshToken)
	}

	// 需要认证的API
	authorized := router.Group("/api/v1")
	authorized.Use(h.AuthMiddleware())
	{
		authorized.POST("/users/logout", h.Logout)
		authorized.GET("/users/:id", h.GetUserInfo)
		authorized.PUT("/users/:id", h.UpdateUserInfo)
		authorized.PUT("/users/:id/password", h.ChangePassword)
	}

	// 管理员API
	admin := router.Group("/api/v1/admin")
	admin.Use(h.AdminMiddleware())
	{
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
		admin.PUT("/users/:id/status", h.UpdateUserStatus)
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Password string `json:"password" binding:"required,min=6,max=32"`
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建用户实体
	user := &entity.User{
		Username: req.Username,
		Nickname: req.Username, // 默认昵称与用户名相同
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 调用注册服务
	createdUser, err := h.authService.RegisterUser(c, user, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"data":    toUserResponse(createdUser),
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用登录服务
	user, credential, err := h.authService.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"user":          toUserResponse(user),
			"access_token":  credential.AccessToken,
			"refresh_token": credential.RefreshToken,
			"expires_in":    int64(credential.ExpiresIn.Seconds()),
		},
	})
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供令牌"})
		return
	}

	// 调用登出服务
	err := h.authService.Logout(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新令牌
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用刷新令牌服务
	credential, err := h.authService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "刷新令牌成功",
		"data": gin.H{
			"access_token":  credential.AccessToken,
			"refresh_token": credential.RefreshToken,
			"expires_in":    int64(credential.ExpiresIn.Seconds()),
		},
	})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": toUserResponse(user),
	})
}

// UpdateUserInfo 更新用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取当前用户
	currentUserID := c.GetInt64("user_id")
	if currentUserID != id && c.GetInt("user_role") != 2 { // 不是本人且不是管理员
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	var req struct {
		Nickname string     `json:"nickname"`
		Avatar   string     `json:"avatar"`
		Email    string     `json:"email" binding:"omitempty,email"`
		Phone    string     `json:"phone"`
		Gender   string     `json:"gender" binding:"omitempty,oneof=male female unknown"`
		Birthday *time.Time `json:"birthday"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	user, err := h.userService.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthday != nil {
		user.Birthday = *req.Birthday
	}

	// 调用更新服务
	updatedUser, err := h.userService.UpdateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    toUserResponse(updatedUser),
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取当前用户
	currentUserID := c.GetInt64("user_id")
	if currentUserID != id && c.GetInt("user_role") != 2 { // 不是本人且不是管理员
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用修改密码服务
	err = h.userService.ChangePassword(c, id, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// ListUsers 获取用户列表（管理员）
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取用户列表
	users, total, err := h.userService.ListUsers(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应对象
	var userResponses []gin.H
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users":      userResponses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateUser 创建用户（管理员）
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Password string `json:"password" binding:"required,min=6,max=32"`
		Nickname string `json:"nickname"`
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone"`
		Role     int    `json:"role" binding:"omitempty,oneof=1 2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建用户实体
	user := &entity.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	if user.Role == 0 {
		user.Role = 1 // 默认为普通用户
	}

	// 调用创建用户服务
	createdUser, err := h.userService.CreateUser(c, user, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "创建成功",
		"data":    toUserResponse(createdUser),
	})
}

// DeleteUser 删除用户（管理员）
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取当前用户ID
	currentUserID := c.GetInt64("user_id")
	if currentUserID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账号"})
		return
	}

	// 调用删除用户服务
	err = h.userService.DeleteUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdateUserStatus 更新用户状态（管理员）
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取当前用户ID
	currentUserID := c.GetInt64("user_id")
	if currentUserID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能修改自己的账号状态"})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2 3"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用更新状态服务
	err = h.userService.UpdateUserStatus(c, id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// AuthMiddleware 认证中间件
func (h *UserHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := h.authService.ValidateToken(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func (h *UserHandler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行认证中间件
		h.AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// 检查是否是管理员
		role := c.GetInt("user_role")
		if role != 2 { // 2表示管理员
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 将用户实体转换为HTTP响应对象
func toUserResponse(user *entity.User) gin.H {
	return gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"email":      user.Email,
		"phone":      user.Phone,
		"gender":     user.Gender,
		"birthday":   user.Birthday,
		"status":     user.Status,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}
