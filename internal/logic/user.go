package logic

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pickcampus-backend/internal/bootstrap"
	"pickcampus-backend/internal/common"
	"pickcampus-backend/internal/types"
	"pickcampus-backend/models"
	"pickcampus-backend/models/factory"
	"pickcampus-backend/models/repo"
)

// UserLogic 用户业务逻辑。
type UserLogic struct {
	Ctx context.Context
	DB  *gorm.DB
}

// NewUserLogic 构造用户逻辑，直接从单例拿 db。
func NewUserLogic(ctx context.Context) *UserLogic {
	return &UserLogic{Ctx: ctx, DB: bootstrap.Cli(ctx)}
}

// Register 注册：校验邮箱唯一 → bcrypt 哈希 → 落库。
func (l *UserLogic) Register(req *types.RegisterRequest) (*types.UserInfo, error) {
	userRepo := factory.UserRepo(l.DB)

	exists, err := userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	if exists {
		return nil, NewBizError(common.ErrCodeEmailExists, "该邮箱已注册")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewBizError(common.ErrCodeCryptoError, "密码加密失败")
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Username:     req.Username,
		Role:         common.UserRoleUser,
		Status:       common.UserStatusActive,
	}
	if err := userRepo.Create(user); err != nil {
		// 并发注册竞态：ExistsByEmail 放行后被 DB 唯一索引兜底拦下，回报正确提示
		if errors.Is(err, repo.ErrDuplicateEmail) {
			return nil, NewBizError(common.ErrCodeEmailExists, "该邮箱已注册")
		}
		return nil, NewBizError(common.ErrCodeDatabaseError, "创建用户失败")
	}

	return toUserInfo(user), nil
}

// Login：查用户 → bcrypt 比对 → 签 JWT → 写 Redis 会话。
// 用户不存在与密码错误统一返回「邮箱或密码错误」，防账号枚举。
func (l *UserLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	userRepo := factory.UserRepo(l.DB)

	user, err := userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, NewBizError(common.ErrCodeInvalidCredentials, "邮箱或密码错误")
		}
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	if user.Status == common.UserStatusDisabled {
		return nil, NewBizError(common.ErrCodeForbidden, "账户已被禁用")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, NewBizError(common.ErrCodeInvalidCredentials, "邮箱或密码错误")
	}

	token, expiresIn, err := common.GenerateToken(user.ID, user.Email, user.Role, user.Status)
	if err != nil {
		return nil, NewBizError(common.ErrCodeTokenError, "签发 token 失败")
	}

	// 写 Redis 会话（token:{uid}），过期时间与 token 一致
	crud := bootstrap.NewCRUD(l.Ctx, bootstrap.GetCli())
	if err := crud.Set(common.GetUserTokenRedisKey(user.ID), token, time.Duration(expiresIn)*time.Second); err != nil {
		return nil, NewBizError(common.ErrCodeInternalError, "会话写入失败")
	}

	return &types.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		UserInfo:  toUserInfo(user),
	}, nil
}

// Logout：删除 Redis 会话，使该用户所有旧 token 失效。
func (l *UserLogic) Logout(userID int64) error {
	crud := bootstrap.NewCRUD(l.Ctx, bootstrap.GetCli())
	if err := crud.Delete(common.GetUserTokenRedisKey(userID)); err != nil {
		return NewBizError(common.ErrCodeInternalError, "登出失败")
	}
	return nil
}

// GetUserInfo：拿当前用户信息。
func (l *UserLogic) GetUserInfo(userID int64) (*types.UserInfo, error) {
	userRepo := factory.UserRepo(l.DB)
	user, err := userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, NewBizError(common.ErrCodeUserNotFound, "用户不存在")
		}
		return nil, NewBizError(common.ErrCodeDatabaseError, "数据库错误")
	}
	return toUserInfo(user), nil
}

// toUserInfo 领域模型转对外 DTO。
func toUserInfo(u *models.User) *types.UserInfo {
	return &types.UserInfo{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		Role:       u.Role,
		Status:     u.Status,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}
