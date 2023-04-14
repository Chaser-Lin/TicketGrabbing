package services

import (
	"Project/MyProject/cache"
	"Project/MyProject/dao"
	"Project/MyProject/models"
	"Project/MyProject/response"
	"Project/MyProject/utils"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 用户登录服务参数
type UserLoginService struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

// 用户注册服务参数
type UserRegisterService struct {
	Email      string `json:"email" form:"email" binding:"required,email"`
	Username   string `json:"username" form:"username" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required,min=6"`
	VerifyCode string `json:"verify_code" form:"verify_code" binding:"required,number,len=6"`
}

// 修改用户信息服务参数
type UpdateUsernameService struct {
	UserID int `json:"user_id" form:"user_id" binding:"number"`
	// 修改密码需要旧密码，以及重复输入新密码
	Username string `json:"username" form:"username" binding:"required"`
}

// 修改用户信息服务参数
type UpdatePasswordService struct {
	UserID int `json:"user_id" form:"user_id" binding:"number"`
	// 修改密码需要旧密码，以及重复输入新密码
	OldPassword      string `json:"old_password" form:"old_password" binding:"min=6,required"`
	NewPassword      string `json:"new_password" form:"new_password" binding:"min=6,required"`
	NewPasswordAgain string `json:"new_password_again" form:"new_password_again" binding:"min=6,required"`
}

// 修改邮箱服务参数
type UpdateEmailService struct {
	UserID int `json:"user_id" form:"user_id" binding:"number"`
	// 修改邮箱需要给新邮箱发送验证码
	Email      string `json:"email" form:"email" binding:"email,required"`
	VerifyCode string `json:"verify_code" form:"verify_code" binding:"number,len=6,required"`
}

// 用户相关服务：登录、注册、获取用户信息
type UserServiceImplement interface {
	Login(*UserLoginService) (accessToken string, refreshToken string, isAdmin bool, err error)
	Register(*UserRegisterService) error
	GetUserInfo(userID int) (*models.User, error)
	UpdateUsername(*UpdateUsernameService) error
	UpdatePassword(*UpdatePasswordService) error
	UpdateEmail(*UpdateEmailService) error
	CheckUserExist(email string) (bool, error)
}

// 实现用户服务接口的实例
type UserService struct {
	UserDal dao.UserDaoImplement
}

func NewUserServices(userDal dao.UserDaoImplement) UserServiceImplement {
	return &UserService{userDal}
}

/*
管理员账户admin

	id = 1
	username = admin
	password = admin123456
	email = admin@admin.com
*/
const (
	AdminID       = 1
	AdminName     = "admin"
	AdminEmail    = "admin@admin.com"
	AdminPassword = "admin123456"
)

func (u *UserService) Login(service *UserLoginService) (accessToken string, refreshToken string, isAdmin bool, err error) {
	var userInfo *models.User
	// 管理员账户，特殊判断
	if service.Email == AdminEmail {
		isAdmin = true
	}
	userInfo, err = u.UserDal.GetUserByEmail(service.Email)
	if err != nil {
		err = response.ErrWrongPassword
		return
	}
	if err = utils.CheckPassword(userInfo.HashedPassword, service.Password); err != nil {
		err = response.ErrWrongPassword
		return
	}
	accessToken, _, err = utils.TokenMaker.CreateToken(userInfo.UserID, userInfo.Email, utils.AccessTokenDurationTime)
	if err != nil {
		err = response.ErrCreateToken
		return
	}
	refreshToken, _, err = utils.TokenMaker.CreateToken(userInfo.UserID, userInfo.Email, utils.RefreshTokenDurationTime)
	if err != nil {
		err = response.ErrCreateToken
		return
	}
	err = cache.SetSession(userInfo.UserID, refreshToken)
	if err != nil {
		err = response.ErrRedisOperation
		return
	}
	//}

	return
}

func (u *UserService) Register(service *UserRegisterService) error {
	_, err := u.UserDal.GetUserByEmail(service.Email)
	if err == nil {
		return response.ErrEmailExist
	}

	hashedPassword, err := utils.HashPassword(service.Password)
	if err != nil {
		return response.ErrEncrypt
	}
	user := &models.User{
		Email:          service.Email,
		HashedPassword: hashedPassword,
		Username:       service.Username,
	}

	if err = u.UserDal.AddUser(user); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrUserExist
			}
		}
		return response.ErrDbOperation
	}
	return nil
}

func (u *UserService) GetUserInfo(userID int) (*models.User, error) {
	user, err := u.UserDal.GetUserByID(userID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrUserNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return user, nil
}

func (u *UserService) UpdatePassword(service *UpdatePasswordService) error {
	hashedPassword, err := utils.HashPassword(service.NewPassword)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:         service.UserID,
		HashedPassword: hashedPassword,
	}

	err = u.UserDal.UpdateUser(user)
	if err == gorm.ErrRecordNotFound {
		return response.ErrUserNotExist
	} else if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (u *UserService) UpdateUsername(service *UpdateUsernameService) error {
	user := &models.User{
		UserID:   service.UserID,
		Username: service.Username,
	}

	err := u.UserDal.UpdateUser(user)
	if err == gorm.ErrRecordNotFound {
		return response.ErrUserNotExist
	} else if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (u *UserService) UpdateEmail(service *UpdateEmailService) error {
	user := &models.User{
		UserID: service.UserID,
		Email:  service.Email,
	}

	err := u.UserDal.UpdateUser(user)
	if err == gorm.ErrRecordNotFound {
		return response.ErrUserNotExist
	} else if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (u *UserService) CheckUserExist(email string) (exist bool, err error) {
	_, err = u.UserDal.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, response.ErrDbOperation
	}
	return true, response.ErrEmailExist
}

func CreateAdmin() {
	userDal := dao.NewUserDao()
	if _, err := userDal.GetUserByEmail(AdminEmail); err == gorm.ErrRecordNotFound {
		hashedPassword, err := utils.HashPassword(AdminPassword)
		if err != nil {
			log.Fatal("CreateAdmin err: ", err)
		}
		admin := &models.User{
			Username:       AdminName,
			HashedPassword: hashedPassword,
			Email:          AdminEmail,
		}
		if err = userDal.AddUser(admin); err != nil {
			log.Fatal("创建管理员账户失败")
		}
	}
}
