package user

import (
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/pkg/jwt"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRouter(r *gin.RouterGroup) {
	r.GET("/user/list", ListUser)
	r.GET("/user/verifytoken", VerifyToken)
	r.GET("/user/login", Login)
	r.POST("/user/register", RegisterUser)
	// 	r.PUT("/project/update", UpdateProject)
	// 	r.DELETE("/project/delete", DeleteProject)
}

type Project struct {
	Name      string `json:"name"`
	ProjectID int64  `json:"project_id"`
}

// @BasePath /api/v1
// PingUser godoc
// @Summary ListUser
// @Schemes
// @Description List User
// @Tags ListUser
// @Accept json
// @Produce json
// @Param username query string false "UserName"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /user/list [get]
// @ID ListUser
func ListUser(ctx *gin.Context) {
	username := ctx.Query("username")

	result := []model.User{}
	q := query.Use(global.Config.DB.DB()).User
	if username == "" {
		err := q.WithContext(ctx).Scan(&result)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, result, error_msg)
			return
		}
	} else {
		err := q.WithContext(ctx).Where(q.Username.Eq(username)).Scan(&result)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, result, error_msg)
			return
		}

	}
	if len(result) == 0 {
		msg := fmt.Sprintf("The UserName: [%s] is not Found!", username)
		util.ReturnMsg(ctx, http.StatusNotFound, result, msg)
	} else {
		msg := "success"
		util.ReturnMsg(ctx, http.StatusOK, result, msg)
	}
}

// @BasePath /api/v1
// PingUser godoc
// @Summary VerifyToken
// @Schemes
// @Description Verify Token
// @Tags VerifyToken
// @Accept json
// @Produce json
// @Param token query string true "token"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /user/verifytoken [get]
// @ID VerifyToken
func VerifyToken(ctx *gin.Context) {
	token := ctx.Query("token")

	mc, err := jwt.ParseToken(token)
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		util.ReturnMsg(ctx, http.StatusUnauthorized, "", msg)
		return
	}
	util.ReturnMsg(ctx, http.StatusOK, mc.Username, "success")
}

// @BasePath /api/v1
// PingUser godoc
// @Summary Login
// @Schemes
// @Description User Login
// @Tags Login
// @Accept json
// @Produce json
// @Param username query string true "UserName"
// @Param password query string true "PassWord"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /user/login [get]
// @ID Login
func Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	q := query.Use(global.Config.DB.DB()).User
	usernameList, err := q.WithContext(ctx).Where(q.Username.Eq(username)).First()

	if err != nil {
		if usernameList == nil {
			error_msg := fmt.Sprintf("The UserName: [%s] is not Found!", username)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(usernameList.Password), []byte(password))
	if err != nil {
		util.ReturnMsg(ctx, http.StatusUnauthorized, "", "账号或者密码错误")
		return
	} else {
		tokenString, _ := jwt.GenToken(username, password)
		util.ReturnMsg(ctx, http.StatusOK, tokenString, "success")

	}
}

type CreateUserRequestBody struct {
	UserName string `json:"username" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

// @BasePath /api/v1
// PingUser godoc
// @Summary RegisterUser
// @Schemes
// @Description Register User
// @Tags RegisterUser
// @Accept json
// @Produce json
// @Param ReqeustBody body CreateUserRequestBody true "Register User"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /user/register [post]
// @ID RegisterUser
func RegisterUser(ctx *gin.Context) {
	body := CreateUserRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	q := query.Use(global.Config.DB.DB()).User
	_, err = q.WithContext(ctx).Where(q.Username.Eq(body.UserName)).First()

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The UserName: [%s] already exists!", body.UserName)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

		return
	}

	user := model.User{}
	//加密密码
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(body.PassWord), bcrypt.DefaultCost)

	user.Username = body.UserName

	user.Password = string(hashPwd)

	err = q.WithContext(ctx).Create(&user)

	if err != nil {
		msg := fmt.Sprintf("The UserName: [%s] Register Failed!, err: %s", body.UserName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", msg)
	} else {
		ok_msg := fmt.Sprintf("The UserName: [%s] Register successfully!", body.UserName)
		util.ReturnMsg(ctx, http.StatusOK, user, ok_msg)

	}

}
