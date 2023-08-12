package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/service"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	birthdayExp *regexp.Regexp
	nickNameExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		birthdayRegexPattern = `^(?:(?:(?:19|20)\d\d)-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12][0-9]|3[01]))$`
		nickNameRegexPattern = `^[\p{L}\d_-]{3,8}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	birthdayExp := regexp.MustCompile(birthdayRegexPattern, regexp.None)
	nickNameExp := regexp.MustCompile(nickNameRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		birthdayExp: birthdayExp,
		nickNameExp: nickNameExp,
	}
}

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	fmt.Printf("now here")
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	// 这边就是数据库操作
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "注册失败")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		print("22222222222")
		print("err:{}", err)
		fmt.Print(err)
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不正确")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	// 登录成功,设置 session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Id       int64  `json:"id"`
		Birthday string `json:"birthday"`
		NickName string `json:"nickName"`
		Profile  string `json:"profile"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 生日校验
	ok, err := u.birthdayExp.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "生日格式有误")
		return
	}
	// 昵称校验
	ok, err = u.nickNameExp.MatchString(req.NickName)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "昵称必须在3-8位之间，包含中英文、数字、特殊字符等")
		return
	}
	// 个人简介校验
	profile := req.Profile
	if utf8.RuneCountInString(profile) >= 20 {
		ctx.String(http.StatusOK, "个人简介最大长度不能超过20个字符")
		return
	}

	err = u.svc.Edit(ctx, domain.User{
		Id:       req.Id,
		Birthday: req.Birthday,
		NickName: req.NickName,
		Profile:  req.Profile,
	})
	if err != nil {
		ctx.String(http.StatusOK, "修改个人信息失败")
		return
	}
	ctx.String(http.StatusOK, "修改个人信息成功")
	return
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	//type EditReq struct {
	//	Id int64 `json:"id"`
	//}
	//var req EditReq
	//if err := ctx.Bind(&req); err != nil {
	//	return
	//}
	id := ctx.Query("id")
	newid, err := strconv.ParseInt(id, 10, 64)
	if newid == 0 {
		ctx.String(http.StatusOK, "请求缺失必要参数")
		return
	}
	user, err := u.svc.Profile(ctx, newid)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
	}
	ctx.JSON(http.StatusOK, user)
	return
}
