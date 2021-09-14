package sessmodels

import "net/http"

type RecipeInterface struct {
	CreateNewSession            func(res http.ResponseWriter, userID string, jwtPayload map[string]interface{}, sessionData map[string]interface{}) (SessionContainer, error)
	GetSession                  func(req *http.Request, res http.ResponseWriter, options *VerifySessionOptions) (*SessionContainer, error)
	RefreshSession              func(req *http.Request, res http.ResponseWriter) (SessionContainer, error)
	GetSessionInformation       func(sessionHandle string) (SessionInformation, error)
	RevokeAllSessionsForUser    func(userID string) ([]string, error)
	GetAllSessionHandlesForUser func(userID string) ([]string, error)
	RevokeSession               func(sessionHandle string) (bool, error)
	RevokeMultipleSessions      func(sessionHandles []string) ([]string, error)
	UpdateSessionData           func(sessionHandle string, newSessionData map[string]interface{}) error
	UpdateJWTPayload            func(sessionHandle string, newJWTPayload map[string]interface{}) error
	GetAccessTokenLifeTimeMS    func() (uint64, error)
	GetRefreshTokenLifeTimeMS   func() (uint64, error)
}
