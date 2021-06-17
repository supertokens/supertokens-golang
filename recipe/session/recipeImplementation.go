package session

import (
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	recipeImplHandshakeInfo *models.HandshakeInfo
	staticConfig            models.TypeNormalisedInput
)

func MakeRecipeImplementation(querier supertokens.Querier, config models.TypeNormalisedInput) models.RecipeImplementation {
	staticConfig = config
	return models.RecipeImplementation{
		CreateNewSession: func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (*models.SessionContainer, error) {
			response, err := createNewSessionHelper(querier, userID, jwtPayload, sessionData)
			if err != nil {
				return nil, err
			}
			attachCreateOrRefreshSessionResponseToRes(config, res, response)
			session := MakeSession(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			return NewSessionContainer(querier, config, session), nil
		},
		GetSession: func(req *http.Request, res http.ResponseWriter, options *models.VerifySessionOptions) (*models.SessionContainer, error) {
			var doAntiCsrfCheck *bool
			if options.AntiCsrfCheck != nil {
				doAntiCsrfCheck = options.AntiCsrfCheck
			}
			idRefreshToken := getIDRefreshTokenFromCookie(req)
			if idRefreshToken == nil {
				if options != nil && options.SessionRequired != nil &&
					!(*options.SessionRequired) {
					return nil, errors.MakeUnauthorizedError("Session does not exist. Are you sending the session tokens in the request as cookies?")
				}
			}
			accessToken := getAccessTokenFromCookie(req)
			if accessToken == nil {
				return nil, errors.TryRefreshTokenError{
					Msg: "Access token has expired. Please call the refresh API",
				}
			}
			if doAntiCsrfCheck == nil {
				doAntiCsrfCheckBool := req.Method != http.MethodGet
				doAntiCsrfCheck = &doAntiCsrfCheckBool
			}
			antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
			response, err := getSessionHelper(querier, *accessToken, antiCsrfToken, *doAntiCsrfCheck, getRidFromHeader(req) != nil)
			if err != nil {
				if errors.IsUnauthorizedError(err) {
					clearSessionFromCookie(config, res)
				}
				return nil, err
			}
			if reflect.DeepEqual(response.AccessToken, models.CreateOrRefreshAPIResponseToken{}) {
				setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInJWT)
				attachAccessTokenToCookie(config, res, response.AccessToken.Token, response.AccessToken.Expiry)
				accessToken = &response.AccessToken.Token
			}
			session := MakeSession(*accessToken, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := NewSessionContainer(querier, config, session)
			return sessionContainer, nil
		},
		RefreshSession: func(req *http.Request, res http.ResponseWriter) (models.SessionContainer, error) {
			inputIdRefreshToken := getIDRefreshTokenFromCookie(req)
			if inputIdRefreshToken == nil {
				return models.SessionContainer{}, errors.MakeUnauthorizedError("Session does not exist. Are you sending the session tokens in the request as cookies?")
			}
			inputRefreshToken := getRefreshTokenFromCookie(req)
			if inputRefreshToken == nil {
				return models.SessionContainer{}, errors.MakeUnauthorizedError("Refresh token not found. Are you sending the refresh token in the request as a cookie?")
			}
			antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
			response, err := refreshSessionHelper(querier, *inputRefreshToken, antiCsrfToken, getRidFromHeader(req) != nil)
			if err != nil {
				if errors.IsUnauthorizedError(err) || errors.IsTokenTheftDetectedError(err) {
					clearSessionFromCookie(config, res)
				}
				return models.SessionContainer{}, err
			}
			attachCreateOrRefreshSessionResponseToRes(config, res, response)
			session := MakeSession(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := NewSessionContainer(querier, config, session)
			return *sessionContainer, nil
		},
		RevokeAllSessionsForUser: func(userID string) ([]string, error) {
			return revokeAllSessionsForUserHelper(querier, userID)
		},
		GetAllSessionHandlesForUser: func(userID string) ([]string, error) {
			return getAllSessionHandlesForUserHelper(querier, userID)
		},
		RevokeSession: func(sessionHandle string) (bool, error) {
			return revokeSessionHelper(querier, sessionHandle)
		},
		RevokeMultipleSessions: func(sessionHandles []string) ([]string, error) {
			return revokeMultipleSessionsHelper(querier, sessionHandles)
		},
		GetSessionData: func(sessionHandle string) (interface{}, error) {
			return getSessionDataHelper(querier, sessionHandle)
		},
		UpdateSessionData: func(sessionHandle string, newSessionData interface{}) error {
			return updateSessionDataHelper(querier, sessionHandle, newSessionData)
		},
		GetJWTPayload: func(sessionHandle string) (interface{}, error) {
			return getJWTPayloadHelper(querier, sessionHandle)
		},
		UpdateJWTPayload: func(sessionHandle string, newJWTPayload interface{}) error {
			return updateJWTPayloadHelper(querier, sessionHandle, newJWTPayload)
		},
	}
}

func GetHandshakeInfo(querier supertokens.Querier) (models.HandshakeInfo, error) {
	if recipeImplHandshakeInfo == nil {
		antiCsrf := staticConfig.AntiCsrf
		path, err := supertokens.NewNormalisedURLPath("/recipe/handshake")
		if err != nil {
			return models.HandshakeInfo{}, err
		}
		response, err := querier.SendPostRequest(*path, nil)
		if err != nil {
			return models.HandshakeInfo{}, err
		}
		recipeImplHandshakeInfo = &models.HandshakeInfo{
			JWTSigningPublicKey:            response["jwtSigningPublicKey"].(string),
			AntiCsrf:                       antiCsrf,
			AccessTokenBlacklistingEnabled: response["accessTokenBlacklistingEnabled"].(bool),
			JWTSigningPublicKeyExpiryTime:  response["jwtSigningPublicKeyExpiryTime"].(uint64),
			AccessTokenValidity:            response["accessTokenValidity"].(uint64),
			RefreshTokenValidity:           response["refreshTokenValidity"].(uint64),
		}
	}
	return *recipeImplHandshakeInfo, nil
}

func UpdateJwtSigningPublicKeyInfo(newKey string, newExpiry uint64) {
	if recipeImplHandshakeInfo == nil {
		recipeImplHandshakeInfo.JWTSigningPublicKey = newKey
		recipeImplHandshakeInfo.JWTSigningPublicKeyExpiryTime = newExpiry
	}
}
