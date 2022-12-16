package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shengbox/oauth2-server/ctrl"

	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	sessionStroe := cookie.NewStore([]byte("secret"))
	g := gin.Default()
	g.LoadHTMLGlob("web/*")
	g.Static("/static", "static")
	g.Use(sessions.Sessions("mysession", sessionStroe))
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)
	ginserver.InitServer(manager)
	ginserver.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)

	ginserver.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		return
	})
	ginserver.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		session, err := sessionStroe.Get(r, "mysession")
		user := session.Values["userID"]
		userID = user.(string)
		return
	})
	ginserver.DefaultConfig.ErrorHandleFunc = func(context *gin.Context, err error) {
		context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token 无效"})
		context.Abort()
	}
	ginserver.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	ginserver.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	g.GET("/captcha", ctrl.GetCaptcha)
	g.GET("/login", ctrl.LoginPage)
	g.POST("/login", ctrl.Login)
	g.GET("/token", ginserver.HandleTokenRequest)
	g.POST("/token", ginserver.HandleTokenRequest)
	g.GET("/authorize", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("userID") == nil {
			loginPage := fmt.Sprintf("/login?client_id=%s&response_type=%s&redirect_uri=%s&state=%s", c.Query("client_id"), c.Query("response_type"), c.Query("redirect_uri"), c.Query("state"))
			c.Redirect(http.StatusTemporaryRedirect, loginPage)
		}
		ginserver.HandleAuthorizeRequest(c)
	})
	api := g.Group("/")
	{
		api.Use(ginserver.HandleTokenVerify())
		api.GET("/user", func(c *gin.Context) {
			ti, exists := c.Get(ginserver.DefaultConfig.TokenKey)
			if exists {
				ti := ti.(*models.Token)
				c.JSON(http.StatusOK, gin.H{"clientId": ti.GetClientID(), "userId": ti.GetUserID()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 404, "message": "not found"})
		})
	}
	g.GET("/", ctrl.Index)
	g.GET("/logout", ctrl.Logout)
	_ = g.Run(":9096")
}
