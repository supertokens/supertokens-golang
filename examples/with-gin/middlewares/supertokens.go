package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Supertokens() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := supertokens.Middleware(func(rw http.ResponseWriter, r *http.Request) {
			c.Next()
		})
		handler(c.Writer, c.Request)
	}
}

func VerifySession(options *models.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionfunc := session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.Set(strconv.Itoa(models.SessionContext), session.GetSessionFromRequest(r))
			c.Next()
		})
		sessionfunc(c.Writer, c.Request)
	}
}
