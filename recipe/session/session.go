package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Session struct {
	sessionHandle        string
	userID               string
	userDataInJWT        interface{}
	res                  *http.ResponseWriter
	accessToken          string
	recipeImplementation schema.RecipeImplementation
}

func MakeSession(recipeImplementation schema.RecipeImplementation, accessToken string, sessionHandle string, userID string, userDataInJWT interface{}, res *http.ResponseWriter) *Session {
	return &Session{
		sessionHandle:        sessionHandle,
		userID:               userID,
		userDataInJWT:        userDataInJWT,
		res:                  res,
		accessToken:          accessToken,
		recipeImplementation: recipeImplementation,
	}
}

// RevokeSession function used to revoke a session for this session
func (session *Session) RevokeSession() error {
	success, err := revokeSession(session.recipeImplementation, session.sessionHandle)
	if err != nil {
		return err
	}
	if success {
		clearSessionFromCookie(session.recipeImplementation.Config, *session.res)
	}
	return nil
}

// GetSessionData function used to get session data for this session
func (session *Session) GetSessionData() (map[string]interface{}, error) {
	data, err := getSessionData(session.recipeImplementation, session.sessionHandle)
	if err != nil {
		if errors.IsUnauthorizedError(err) {
			clearSessionFromCookie(session.recipeImplementation.Config, *session.res)
		}
		return nil, err
	}
	return data, nil
}

// UpdateSessionData function used to update session data for this session
func (session *Session) UpdateSessionData(newSessionData map[string]interface{}) error {
	err := updateSessionData(session.recipeImplementation, session.sessionHandle, newSessionData)
	if err != nil {
		if errors.IsUnauthorizedError(err) {
			clearSessionFromCookie(session.recipeImplementation.Config, *session.res)
		}
		return err
	}
	return nil
}

// GetUserID function gets the user for this session
func (session *Session) GetUserID() string {
	return session.userID
}

// GetJWTPayload function gets the jwt payload for this session
func (session *Session) GetJWTPayload() interface{} {
	return session.userDataInJWT
}

// GetHandle function gets the session handle for this session
func (session *Session) GetHandle() string {
	return session.sessionHandle
}

// GetAccessToken function gets the access token for this session
func (session *Session) GetAccessToken() string {
	return session.accessToken
}

// UpdateJWTPayload function used to update jwt payload for this session
func (session *Session) UpdateJWTPayload(newJWTPayload map[string]interface{}) error {
	path, err := supertokens.NewNormalisedURLPath("/recipe/session/regenerate")
	if err != nil {
		return err
	}
	response, err := session.recipeImplementation.Querier.SendPostRequest(*path, map[string]interface{}{
		"accessToken":   session.accessToken,
		"userDataInJWT": newJWTPayload,
	})
	if response["status"] == UnauthorizedError {
		clearSessionFromCookie(session.recipeImplementation.Config, *session.res)
		return errors.UnauthorizedError{
			Msg: "Session has probably been revoked while updating JWT payload",
		}
	}
	return nil
}
