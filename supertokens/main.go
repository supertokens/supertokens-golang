package supertokens

import (
	"net/http"
)

func Init(config TypeInput) error {
	return supertokensInit(config)
}

func Middleware(theirHandler http.Handler) http.Handler {
	instance, err := getInstanceOrThrowError()
	if err != nil {
		panic("Please call SupertokensInit before using the middleware")
	}
	return instance.middleware(theirHandler)
}

// TODO: Make this match (err error, w http.ResponseWriter, r *http.Request) or the standard http.Error() signature.
func ErrorHandler(err error, req *http.Request, res http.ResponseWriter) {
	instance, instanceErr := getInstanceOrThrowError()
	if instanceErr != nil {
		panic("Please call SupertokensInit before using the error handler")
	}
	instance.errorHandler(err, req, res)
}

func GetAllCORSHeaders() []string {
	instance, err := getInstanceOrThrowError()
	if err != nil {
		panic("Please call SupertokensInit before using the GetAllCORSHeaders function")
	}
	return instance.getAllCORSHeaders()
}

func GetUserCount(includeRecipeIds *[]string) (int, error) {
	return getUserCount(includeRecipeIds)
}

func GetUsersOldestFirst(limit *int, paginationToken *string, includeRecipeIds *[]string) (*UserPaginationResult, error) {
	return getUsers("ASC", limit, paginationToken, includeRecipeIds)
}

func GetUsersNewestFirst(limit *int, paginationToken *string, includeRecipeIds *[]string) (*UserPaginationResult, error) {
	return getUsers("DESC", limit, paginationToken, includeRecipeIds)
}
