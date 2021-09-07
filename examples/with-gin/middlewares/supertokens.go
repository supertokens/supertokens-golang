package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const sessionContext = "supertokens_session_key"

func Supertokens() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := supertokens.Middleware(func(rw http.ResponseWriter, r *http.Request) {
			c.Next()
		})
		handler(c.Writer, c.Request)
	}
}

func VerifySession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionfunc := session.VerifySession(nil)
		sessionfunc(c.Writer, c.Request, func(rw http.ResponseWriter, r *http.Request) {
			sessionContainer, err := session.GetSession(c.Request, c.Writer, nil)
			if err != nil {
				fmt.Println(err)
				c.Abort()
			}
			if sessionContainer != nil {
				c.Set(sessionContext, sessionContainer)
			}
			c.Next()
		})
	}
}
