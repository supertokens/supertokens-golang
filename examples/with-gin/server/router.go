package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/examples/with-gin/controllers"
	"github.com/supertokens/supertokens-golang/examples/with-gin/middlewares"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     append([]string{"content-type"}, supertokens.GetAllCORSHeaders()...),
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	auth := router.Group("/auth", middlewares.Supertokens())
	auth.POST("/*action")
	auth.GET("/*action")
	router.GET("/sessioninfo", middlewares.VerifySession(), controllers.Sessioninfo)

	return router
}
