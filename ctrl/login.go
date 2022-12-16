package ctrl

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginParam struct {
	Captcha  string `form:"captcha" json:"captcha" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var param LoginParam
	if err := c.ShouldBind(&param); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"msg": "用户名或密码错误"})
		return
	}

	if !CaptchaVerify(c, param.Captcha) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"msg":           "验证码错误",
			"username":      param.Username,
			"client_id":     c.PostForm("client_id"),
			"redirect_uri":  c.PostForm("redirect_uri"),
			"response_type": c.PostForm("response_type"),
		})
		return
	}
	if param.Username == "admin" && param.Password == "admin" {
		session := sessions.Default(c)
		session.Set("userID", param.Username)
		session.Save()
		if c.PostForm("client_id") == "" {
			c.Redirect(http.StatusFound, "https://www.zjzwfw.gov.cn")
		} else {
			c.Redirect(http.StatusFound, "/authorize?response_type=code&client_id="+c.PostForm("client_id")+"&redirect_uri="+c.PostForm("redirect_uri"))
		}
	} else {
		c.HTML(http.StatusOK, "login.html", gin.H{"msg": "用户名或密码错误"})
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"client_id":     c.Query("client_id"),
		"redirect_uri":  c.Query("redirect_uri"),
		"response_type": c.Query("response_type"),
	})
}

func Index(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userID") == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
