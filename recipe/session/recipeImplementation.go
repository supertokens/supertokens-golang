package session

import (
	defaultErrors "errors"
	"net/http"
	"reflect"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var handshakeInfoLock sync.Mutex

func makeRecipeImplementation(querier supertokens.Querier, config models.TypeNormalisedInput) models.RecipeImplementation {
	var recipeImplHandshakeInfo *models.HandshakeInfo = nil
	getHandshakeInfo(recipeImplHandshakeInfo, config, querier)

	return models.RecipeImplementation{
		CreateNewSession: func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (*models.SessionContainer, error) {
			response, err := createNewSessionHelper(recipeImplHandshakeInfo, config, querier, userID, jwtPayload, sessionData)
			if err != nil {
				return nil, err
			}

			// TODO: For POC: we are setting cookies to this `res` inside that function.
			// But since it's not a pointer, will it work?
			attachCreateOrRefreshSessionResponseToRes(config, res, *response)
			sessionContainerInput := MakeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			return NewSessionContainer(querier, config, sessionContainerInput), nil
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
					return nil, nil
				}
				return nil, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
			}

			accessToken := getAccessTokenFromCookie(req)
			if accessToken == nil {
				return nil, errors.TryRefreshTokenError{
					Msg: "Access token has expired. Please call the refresh API",
				}
			}

			antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
			if doAntiCsrfCheck == nil {
				doAntiCsrfCheckBool := req.Method != http.MethodGet
				doAntiCsrfCheck = &doAntiCsrfCheckBool
			}

			response, err := getSessionHelper(recipeImplHandshakeInfo, config, querier, *accessToken, antiCsrfToken, *doAntiCsrfCheck, getRidFromHeader(req) != nil)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					// TODO: For POC: will this set these cookies / headers in the final response
					// sent by our API / the user?
					clearSessionFromCookie(config, res)
				}
				return nil, err
			}

			if !reflect.DeepEqual(response.AccessToken, models.CreateOrRefreshAPIResponseToken{}) {
				setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInJWT)
				attachAccessTokenToCookie(config, res, response.AccessToken.Token, response.AccessToken.Expiry)
				accessToken = &response.AccessToken.Token
			}
			sessionContainerInput := MakeSessionContainerInput(*accessToken, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := NewSessionContainer(querier, config, sessionContainerInput)
			return sessionContainer, nil
		},

		RefreshSession: func(req *http.Request, res http.ResponseWriter) (*models.SessionContainer, error) {
			inputIdRefreshToken := getIDRefreshTokenFromCookie(req)
			if inputIdRefreshToken == nil {
				return nil, errors.UnauthorizedError{Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?"}
			}

			inputRefreshToken := getRefreshTokenFromCookie(req)
			if inputRefreshToken == nil {
				clearSessionFromCookie(config, res)
				return nil, errors.UnauthorizedError{Msg: "Refresh token not found. Are you sending the refresh token in the request as a cookie?"}
			}

			antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
			response, err := refreshSessionHelper(recipeImplHandshakeInfo, config, querier, *inputRefreshToken, antiCsrfToken, getRidFromHeader(req) != nil)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) || defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
					clearSessionFromCookie(config, res)
				}
				return nil, err
			}
			attachCreateOrRefreshSessionResponseToRes(config, res, *response)
			sessionContainerInput := MakeSessionContainerInput(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			sessionContainer := NewSessionContainer(querier, config, sessionContainerInput)
			return sessionContainer, nil
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

		GetAccessTokenLifeTimeMS: func() (uint64, error) {
			handShake, err := getHandshakeInfo(recipeImplHandshakeInfo, config, querier)
			if err != nil {
				return 0, err
			}
			return handShake.AccessTokenValidity, nil
		},

		GetRefreshTokenLifeTimeMS: func() (uint64, error) {
			handShake, err := getHandshakeInfo(recipeImplHandshakeInfo, config, querier)
			if err != nil {
				return 0, err
			}
			return handShake.RefreshTokenValidity, nil
		},
	}
}

func getHandshakeInfo(recipeImplHandshakeInfo *models.HandshakeInfo, config models.TypeNormalisedInput, querier supertokens.Querier) (models.HandshakeInfo, error) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	if recipeImplHandshakeInfo == nil {
		antiCsrf := config.AntiCsrf
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
			SigningKeyLastUpdated:          getCurrTimeInMS(),
		}
	}
	return *recipeImplHandshakeInfo, nil
}

func updateJwtSigningPublicKeyInfo(recipeImplHandshakeInfo *models.HandshakeInfo, newKey string, newExpiry uint64) {
	handshakeInfoLock.Lock()
	defer handshakeInfoLock.Unlock()
	if recipeImplHandshakeInfo == nil {
		if recipeImplHandshakeInfo.JWTSigningPublicKeyExpiryTime != newExpiry ||
			recipeImplHandshakeInfo.JWTSigningPublicKey != newKey {
			recipeImplHandshakeInfo.SigningKeyLastUpdated = getCurrTimeInMS()
		}
		recipeImplHandshakeInfo.JWTSigningPublicKey = newKey
		recipeImplHandshakeInfo.JWTSigningPublicKeyExpiryTime = newExpiry
	}
}
