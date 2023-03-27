package service

import (
	"Project/MyProject/dal"
	"Project/MyProject/dal/models"
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
type UpdateUserInfoService struct {
	UserID   int    `json:"user_id" form:"user_id" binding:"required,number"`
	Username string `json:"username" form:"username"`
	Phone    string `json:"phone" form:"phone"`
	// 修改邮箱需要给新邮箱发送验证码
	Email      string `json:"email" form:"email" binding:"email"`
	VerifyCode string `json:"verify_code" form:"verify_code" binding:"number,len=6"`
	// 修改密码需要旧密码，以及重复输入新密码
	OldPassword      string `json:"old_password" form:"old_password" binding:"min=6"`
	NewPassword      string `json:"new_password" form:"new_password" binding:"min=6"`
	NewPasswordAgain string `json:"new_password_again" form:"new_password_again" binding:"min=6"`
}

// 用户相关服务：登录、注册、获取用户信息
type UserServiceImplement interface {
	Login(*UserLoginService) (token string, err error)
	Register(*UserRegisterService) error
	GetUserInfo(userID int) (*models.User, error)
	UpdateUserInfo(*UpdateUserInfoService) error
	CheckUserExist(email string) (bool, error)
}

// 实现用户服务接口的实例
type UserService struct {
	UserDal dal.UserDalImplement
}

func NewUserServices(userDal dal.UserDalImplement) UserServiceImplement {
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

func (u *UserService) Login(service *UserLoginService) (token string, err error) {
	var userInfo *models.User
	// 管理员账户，特殊判断
	if service.Email == AdminEmail {
		if service.Password != AdminPassword {
			return "", response.ErrWrongPassword
		}
		token, _, err = utils.TokenMaker.CreateToken(AdminID, AdminEmail, utils.TokenDurationTime)
		if err != nil {
			return "", response.ErrCreateToken
		}
	} else {
		userInfo, err = u.UserDal.GetUserByEmail(service.Email)
		if err != nil {
			return "", response.ErrWrongPassword
		}
		if err = utils.CheckPassword(userInfo.HashedPassword, service.Password); err != nil {
			return "", response.ErrWrongPassword
		}
		token, _, err = utils.TokenMaker.CreateToken(userInfo.UserID, userInfo.Email, utils.TokenDurationTime)
		if err != nil {
			return "", response.ErrCreateToken
		}
	}

	return token, nil
}

func (u *UserService) Register(service *UserRegisterService) error {
	_, err := u.UserDal.GetUserByEmail(service.Email)
	if err == nil {
		return response.ErrEmailExist
	}
	//_, err = u.UserDal.GetUserByUsername(service.Email)
	//if err == nil {
	//	return response.ErrUsernameExist
	//}
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

func (u *UserService) UpdateUserInfo(service *UpdateUserInfoService) error {
	var hashedPassword string
	var err error
	if service.NewPassword != "" {
		hashedPassword, err = utils.HashPassword(service.NewPassword)
		if err != nil {
			return err
		}
	}
	user := &models.User{
		UserID:         service.UserID,
		Username:       service.Username,
		HashedPassword: hashedPassword,
		Email:          service.Email,
		Phone:          service.Phone,
	}

	err = u.UserDal.UpdateUser(user)
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
	userDal := dal.NewUserDal()
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
