package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
)

func Sessioninfo(c *gin.Context) {

	var session *models.SessionContainer = Session.GetSessionFromRequest(c.Request)
	if session == nil {
		c.JSON(500, "no session found")
		return
	}
	sessionData, err := session.GetSessionData()
	if err != nil {
		c.JSON(500, "error in sessiondata")
		return
	}
	c.JSON(200, map[string]interface{}{
		"sessionHandle": session.GetHandle(),
		"userId":        session.GetUserID(),
		"jwtPayload":    session.GetJWTPayload(),
		"sessionData":   sessionData,
	})
}
