package service

import (
	"Project/MyProject/response"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"sync"
	"time"
)

type SendVerifyCodeService struct {
	Email string `json:"email" form:"email" binding:"email"`
}

// 邮箱相关服务：发送验证码
type EmailServiceImplement interface {
	SendVerifyCode(*SendVerifyCodeService) error
	CheckVerifyCode(email, verifyCode string) (bool, error)
}

// 实现邮箱服务接口的实例
type EmailService struct {
	manager *VerifyCodeManager
}

func NewEmailServices(manager *VerifyCodeManager) EmailServiceImplement {
	return &EmailService{manager: manager}
}

// 发送邮箱验证码存入map中（5分钟有效时间）
func (e *EmailService) SendVerifyCode(service *SendVerifyCodeService) error {
	//在handler层判断邮箱是否已被注册
	return e.manager.SendVerifyCodeToEmail(service.Email)
}

// 检验邮箱验证码是否正确
func (e *EmailService) CheckVerifyCode(email, verifyCode string) (bool, error) {
	e.manager.mutex.Lock()
	defer e.manager.mutex.Unlock()
	if message, ok := e.manager.m[email]; ok {
		if time.Now().After(message.ExpiredTime) {
			delete(e.manager.m, email) // 删除过期验证码
			return false, response.ErrVerifyCodeExpired
		}
		if message.VerifyCode == verifyCode {
			delete(e.manager.m, email) // 验证通过，删除验证码信息
			return true, nil
		} else {
			message.ErrorTry++
			if message.ErrorTry == MaxErrorTry {
				delete(e.manager.m, email) // 错误尝试次数太多，直接删除该验证码
				return false, response.ErrVerifyCodeMaxTry
			}
			return false, response.ErrVerifyCode
		}
	}
	return false, response.ErrVerifyCodeNotExist
}

const (
	RegisterEmailSubject = "火车票抢票系统用户注册验证码"
	GoMailHost           = "smtp.qq.com"
	GoMailPort           = 465
	GoMailSenderAddr     = "ticket_grabbing@qq.com"
	GoMailSenderName     = "火车票抢票系统管理员"
	GoMailPassword       = "scsoxsdoemggdahg"
	MaxErrorTry          = 5
)

// 验证码信息，包含过期时间
type VerifyCodeMessage struct {
	VerifyCode  string
	ErrorTry    int
	ExpiredTime time.Time
}

// 验证码管理模块，使用map简易实现，mutex保证并发安全
type VerifyCodeManager struct {
	mutex sync.Mutex
	m     map[string]*VerifyCodeMessage
}

func NewVerifyCodeManager() *VerifyCodeManager {
	return &VerifyCodeManager{
		m: make(map[string]*VerifyCodeMessage),
	}
}

func (manager *VerifyCodeManager) AddVerifyCodeMessage(email string, message *VerifyCodeMessage) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.m[email] = message
}

func sendEmail(email string, subject string, content string) error {
	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", GoMailSenderAddr, GoMailSenderName)
	// 收件人
	m.SetHeader("To", m.FormatAddress(email, ""))
	// 主题
	m.SetHeader("Subject", subject)
	// 正文
	m.SetBody("text/html", content)

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewDialer(GoMailHost, GoMailPort, GoMailSenderAddr, GoMailPassword)
	return d.DialAndSend(m)
}

// 生成6位随机验证码并发送邮件
func (manager *VerifyCodeManager) SendVerifyCodeToEmail(email string) error {
	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))
	verifyCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// 设置验证码过期时间为5分钟
	expiredTime := time.Now().Add(5 * time.Minute)

	message := &VerifyCodeMessage{
		VerifyCode:  verifyCode,
		ExpiredTime: expiredTime,
	}

	log.Printf("email:(%s), verifyCode:(%s)", email, verifyCode)

	manager.AddVerifyCodeMessage(email, message)

	content := fmt.Sprintf(`
	<div>
		<div>
			尊敬的%s，您好！
		</div>
		<div style="padding: 8px 40px 8px 50px;">
			<p>您于 %s 提交的邮箱验证，本次验证码为<u><strong>%s</strong></u>，为了保证账号安全，验证码有效期为5分钟。<br>
				如非本人操作，请忽略此邮件，感谢您的理解与使用。</p>
		</div>
		<div>
			<p>此邮箱为系统邮箱，请勿回复。</p>
		</div>
	</div>
	`, email, currentTime, verifyCode)

	err := sendEmail(email, RegisterEmailSubject, content)
	if err != nil {
		return response.ErrInvalidEmail
	}
	return nil
}
