package service

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"mysite/logger"
	"github.com/gorilla/schema"
	"mysite/model"
	"mysite/db"
	"github.com/Sirupsen/logrus"
	"errors"
	"github.com/polaris1119/goutils.git"
	"time"
)

type UserService struct {}

var DefaultUserService = new(UserService)

//注册账号
func (self UserService) RegisterUser(ctx *gin.Context, form url.Values) (errMsg string, err error)  {
	logger := logger.GetLogger()
	logger.Info("start register the user")

	schema := schema.NewDecoder()
	user := &model.User{}
	err = schema.Decode(user, form)
	if err != nil {
		logger.WithFields(logrus.Fields{"err1": err, "data": form}).Info("parse form error")
		errMsg = err.Error()
		return
	}

	//检查是否有重名
	if self.UserExist(ctx, "username", form.Get("username")){
		logger.WithFields(logrus.Fields{"username": form.Get("username")}).Info("用户名称存在重名")
		err = errors.New("用户名称重名")
		errMsg = err.Error()
		return
	}

	//检查邮箱是否有重复的
	if self.UserExist(ctx, "email", form.Get("email")){
		logger.WithFields(logrus.Fields{"email": form.Get("email")}).Info("邮箱名称重名")
		err = errors.New("邮箱重复")
		errMsg = err.Error()
		return
	}

	//写入到数据库
	session := db.MasterDB.NewSession()
	defer session.Close()

	session.Begin()
	//存储基本用户信息
	user.Open = 1
	user.Status = model.UserStatusNoAudit
	if _, err = session.Insert(user); err != nil {
		logger.Errorln("User insert Error:", err)
		session.Rollback()
		errMsg = "内部服务器错误"
		return
	}

	//存储用户登陆信息
	userLogin := &model.UserLogin{}
	form.Set("passwd", ctx.PostForm("passwd"))
	err = schema.Decode(userLogin, form)
	if err != nil {
		logger.WithFields(logrus.Fields{"form": form}).Info("存储登陆信息失败")
		errMsg = errors.New("存储登陆信息失败").Error()
		return
	}
	userLogin.Uid = user.Uid;
	err = userLogin.GenMd5Passwd()
	if err != nil {
		logger.Errorln("生成密码失败", err)
		session.Rollback()
		errMsg = "生成密码失败"
		return
	}
	if _, err = session.Insert(userLogin); err != nil {
		logger.Errorln("UserLogin insert Error:", err)
		session.Rollback()
		errMsg = "内部服务器错误"
		return
	}

	//存用户活跃信息，初始活跃+2
	userActive := &model.UserActive{}
	userActive.Avatar = user.Avatar
	userActive.Email = user.Email
	userActive.Username = user.Username
	userActive.Uid = user.Uid
	userActive.Weight = 2;
	if _, err = session.Insert(userActive); err != nil {
		logger.WithFields(logrus.Fields{"form": form}).Info("存用户活跃信息")
		session.Rollback()
		errMsg = errors.New("存用户活跃信息失败").Error()
		return
	}

	session.Commit()

	return "", nil
}

//判断邮箱或用户名称是否现在重复
func (self UserService) UserExist(ctx *gin.Context, field string, val string) bool {
	logger := logger.GetLogger()
	userLogin := &model.UserLogin{}
	_, err := db.MasterDB.Where(field + "=?", val).Get(userLogin)
	if err != nil || userLogin.Uid == 0 {
		logger.WithFields(logrus.Fields{"err": err, "field": field, "value": val, "user": userLogin}).Info("error")
		return false
	}
	return true;
}

//用户登录
func (self UserService) Login(ctx *gin.Context, userName, passwd string) (userLogin *model.UserLogin, err error) {
	logger := logger.GetLogger()
	userLogin = &model.UserLogin{}
	_, err = db.MasterDB.Where("username=? or email=?", userName, userName).Get(userLogin)
	if err != nil {
		logger.WithFields(logrus.Fields{"username":userName, "passwd": passwd}).Info("用户不存在")
		return
	}

	if userLogin.Uid == 0 {
		logger.WithFields(logrus.Fields{"username":userName, "passwd": passwd}).Info("用户不存在")
		err = errors.New("用户不存在")
		return
	}

	//检查用户状态是否正确
	user := &model.User{}
	user.Uid = userLogin.Uid
	_, err = db.MasterDB.Get(user)
	if err != nil {
		logger.WithFields(logrus.Fields{"username": userName}).Info("获取用户信息失败")
		return
	}
	if user.Status > model.UserStatusAudit {
		errMap := map[int]error {
			model.UserStatusFreeze: errors.New("账号被冻结"),
			model.UserStatusRefuse: errors.New("账号审核拒绝"),
			model.UserStatusOutage: errors.New("账号被停用"),
		}
		logger.WithFields(logrus.Fields{"user": user}).Info("账号异常")
		err = errMap[user.Status]
		return
	}

	//判断密码是否一致
	passWdMd5 := goutils.Md5(passwd + userLogin.Passcode)
	if passWdMd5 != userLogin.Passwd {
		logger.WithFields(logrus.Fields{"username": userName}).Info("密码错误")
		err = errors.New("密码错误")
		return
	}

	return
}

func (self UserService) FindCurrentUser(ctx *gin.Context, userName interface{}) *model.Me {
	logger := logger.GetLogger()
	user := &model.User{}
	_, err := db.MasterDB.Where("username=? and status=?", userName, model.UserStatusAudit).Get(user)
	if err != nil {
		logger.WithFields(logrus.Fields{"user": user, "userName": userName}).Info("获取用户信息失败")
		return &model.Me{}
	}

	if user.Uid == 0 {
		logger.WithFields(logrus.Fields{"user": user, "username": userName}).Info("没有这个用户信息")
		return &model.Me{}
	}

	me := &model.Me{
		Uid: user.Uid,
		Username: user.Username,
		Email: user.Email,
		Avatar: user.Avatar,
		Status: user.Status,
		MsgNum: 0,
		IsRoot: false,
		IsAdmin: false,
	}

	//记录登录时间
	go self.RecordLoginTime(userName)

	return me
}

func (self UserService) RecordLoginTime(userName interface{}) (err error) {
	logger := logger.GetLogger()
	_, err = db.MasterDB.Table(new(model.UserLogin)).Where("username=?", userName).
		Update(map[string]interface{}{"login_time": time.Now()})
	if err != nil {
		logger.WithFields(logrus.Fields{"userName": userName}).Info("更新用户登录时间失败")
		return
	}
	return
}