package session

import (
	"encoding/json"
	defaultErrors "errors"
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionContainerInput struct {
	sessionHandle string
	userID        string
	userDataInJWT map[string]interface{}
	res           http.ResponseWriter
	accessToken   string
}

func makeSessionContainerInput(accessToken string, sessionHandle string, userID string, userDataInJWT map[string]interface{}, res http.ResponseWriter) SessionContainerInput {
	return SessionContainerInput{
		sessionHandle: sessionHandle,
		userID:        userID,
		userDataInJWT: userDataInJWT,
		res:           res,
		accessToken:   accessToken,
	}
}

func newSessionContainer(querier supertokens.Querier, config sessmodels.TypeNormalisedInput, session *SessionContainerInput) sessmodels.SessionContainer {

	return sessmodels.SessionContainer{
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

		GetSessionData: func() (map[string]interface{}, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return nil, err
			}
			return sessionInformation.SessionData, nil
		},

		UpdateSessionData: func(newSessionData map[string]interface{}) error {
			err := updateSessionDataHelper(querier, session.sessionHandle, newSessionData)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return err
			}
			return nil
		},

		UpdateJWTPayload: func(newJWTPayload map[string]interface{}) error {
			if newJWTPayload == nil {
				newJWTPayload = map[string]interface{}{}
			}
			response, err := querier.SendPostRequest("/recipe/session/regenerate", map[string]interface{}{
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

			responseByte, err := json.Marshal(response)
			if err != nil {
				return err
			}
			var resp sessmodels.GetSessionResponse
			err = json.Unmarshal(responseByte, &resp)
			if err != nil {
				return err
			}

			session.userDataInJWT = resp.Session.UserDataInJWT
			if !reflect.DeepEqual(resp.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
				session.accessToken = resp.AccessToken.Token
				setFrontTokenInHeaders(session.res, resp.Session.UserID, resp.AccessToken.Expiry, resp.Session.UserDataInJWT)
				attachAccessTokenToCookie(config, session.res, resp.AccessToken.Token, resp.AccessToken.Expiry)
			}
			return nil
		},
		GetUserID: func() string {
			return session.userID
		},
		GetJWTPayload: func() map[string]interface{} {
			return session.userDataInJWT
		},
		GetHandle: func() string {
			return session.sessionHandle
		},
		GetAccessToken: func() string {
			return session.accessToken
		},
		GetTimeCreated: func() (uint64, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return 0, err
			}
			return sessionInformation.TimeCreated, nil
		},
		GetExpiry: func() (uint64, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return 0, err
			}
			return sessionInformation.Expiry, nil
		},
	}
}
