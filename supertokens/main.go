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
		panic("Please call supertokens.Init function before using the Middleware")
	}
	return instance.middleware(theirHandler)
}

func ErrorHandler(err error, req *http.Request, res http.ResponseWriter) error {
	instance, instanceErr := getInstanceOrThrowError()
	if instanceErr != nil {
		return instanceErr
	}
	return instance.errorHandler(err, req, res)
}

func GetAllCORSHeaders() []string {
	instance, err := getInstanceOrThrowError()
	if err != nil {
		panic("Please call supertokens.Init before using the GetAllCORSHeaders function")
	}
	return instance.getAllCORSHeaders()
}

func GetUserCount(includeRecipeIds *[]string) (float64, error) {
	return getUserCount(includeRecipeIds)
}

func GetUsersOldestFirst(paginationToken *string, limit *int, includeRecipeIds *[]string) (UserPaginationResult, error) {
	return getUsers("ASC", paginationToken, limit, includeRecipeIds)
}

func GetUsersNewestFirst(paginationToken *string, limit *int, includeRecipeIds *[]string) (UserPaginationResult, error) {
	return getUsers("DESC", paginationToken, limit, includeRecipeIds)
}
