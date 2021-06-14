package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var handshakeInfo *models.HandshakeInfo

func MakeRecipeImplementation(querier supertokens.Querier, config models.TypeNormalisedInput) models.RecipeImplementation {
	return models.RecipeImplementation{
		CreateNewSession: func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (models.SessionContainer, error) {
			response, err := CreateNewSession(querier, userID, jwtPayload, sessionData)
			if err != nil {
				return models.SessionContainer{}, err
			}
			attachCreateOrRefreshSessionResponseToRes(config, res, response)
			session := MakeSession(response.AccessToken.Token, response.Session.Handle, response.Session.UserID, response.Session.UserDataInJWT, res)
			return MakeSessionContainer(querier, config, session), nil
		},
		GetSession: func(req *http.Request, res http.ResponseWriter, options *models.VerifySessionOptions) (*models.SessionContainer, error) {
			var doAntiCsrfCheck *bool
			if options.AntiCsrfCheck != nil {
				doAntiCsrfCheck = options.AntiCsrfCheck
			}
			idRefreshToken := getIDRefreshTokenFromCookie(req)
			if idRefreshToken == nil {
				if options != nil && *options.SessionRequired == false {
					return nil, UnauthorizedError{
						Msg: "Session does not exist. Are you sending the session tokens in the request as cookies?",
					}
				}
			}
			accessToken := getAccessTokenFromCookie(req)
			if accessToken == nil {
				return nil, TryRefreshTokenError{
					Msg: "Access token has expired. Please call the refresh API",
				}
			}
			// TODO
			// antiCsrfToken := getAntiCsrfTokenFromHeaders(req)
			if doAntiCsrfCheck == nil {
				doAntiCsrfCheckBool := req.Method != http.MethodGet
				doAntiCsrfCheck = &doAntiCsrfCheckBool
			}
			// response, err := GetSession()
			return nil, nil
		},
	}
}

func GetHandshakeInfo() models.HandshakeInfo {
	// TODO
	return models.HandshakeInfo{}
}

func UpdateJwtSigningPublicKeyInfo(newKey string, newExpiry uint64) {
	// TODO
}
