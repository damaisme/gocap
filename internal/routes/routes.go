package routes

import (
	"github.com/damaisme/gocap/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/")
	{
		api.GET("/", handlers.GetIndex)

		api.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})

		api.GET("/voc", func(c *gin.Context) {
			c.HTML(http.StatusOK, "voc.html", nil)
		})

		api.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "Hallooo")
		})

		api.GET("/generate_204", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "http://gocap.local/")
		})

		api.GET("/logout", handlers.Logout)

		api.POST("/login", handlers.Login)

		api.POST("/buy_voc", handlers.BuyVoucher)

		api.POST("/logout", handlers.Logout)

		api.GET("/finish", handlers.Finish)

	}
}
