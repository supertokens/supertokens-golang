package session

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionContainerInput struct {
	sessionHandle string
	userID        string
	userDataInJWT interface{}
	res           http.ResponseWriter
	accessToken   string
}

func MakeSessionContainerInput(accessToken string, sessionHandle string, userID string, userDataInJWT interface{}, res http.ResponseWriter) SessionContainerInput {
	return SessionContainerInput{
		sessionHandle: sessionHandle,
		userID:        userID,
		userDataInJWT: userDataInJWT,
		res:           res,
		accessToken:   accessToken,
	}
}

func NewSessionContainer(querier supertokens.Querier, config models.TypeNormalisedInput, session SessionContainerInput) *models.SessionContainer {
	return &models.SessionContainer{
		RevokeSession: func() error {
			success, err := revokeSessionHelper(querier, session.sessionHandle)
			if err != nil {
				return err
			}
			if success {
				clearSessionFromCookie(config, session.res)
			}
			return nil
		},

		GetSessionData: func() (interface{}, error) {
			data, err := getSessionDataHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return nil, err
			}
			return data, nil
		},

		UpdateSessionData: func(newSessionData interface{}) error {
			err := updateSessionDataHelper(querier, session.sessionHandle, newSessionData)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return err
			}
			return nil
		},
		UpdateJWTPayload: func(newJWTPayload interface{}) error {
			path, err := supertokens.NewNormalisedURLPath("/recipe/session/regenerate")
			if err != nil {
				return err
			}
			response, err := querier.SendPostRequest(*path, map[string]interface{}{
				"accessToken":   session.accessToken,
				"userDataInJWT": newJWTPayload,
			})
			if err != nil {
				return err
			}
			if response["status"].(string) == errors.UnauthorizedErrorStr {
				clearSessionFromCookie(config, session.res)
				return errors.UnauthorizedError{Msg: "Session has probably been revoked while updating JWT payload"}
			}
			return nil
		},
		GetUserID: func() string {
			return session.userID
		},
		GetJWTPayload: func() interface{} {
			return session.userDataInJWT
		},
		GetHandle: func() string {
			return session.sessionHandle
		},
		GetAccessToken: func() string {
			return session.accessToken
		},
	}
}
