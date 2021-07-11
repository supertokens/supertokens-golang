package supertokens

import "net/http"

func Middleware() func(req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc) {
	return func(req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc) {
		superTokensInstance, err := getInstanceOrThrowError()
		if err != nil {
			panic("supertokens not initialised - " + err.Error())
		}
		middleware := superTokensInstance.Middleware(theirHandler)
		middleware(res, req)
	}
}

func GetAllCORSHeaders() []string {
	superTokensInstance, err := getInstanceOrThrowError()
	if err != nil {
		panic("supertokens not initialised - " + err.Error())
	}
	return superTokensInstance.GetAllCORSHeaders()
}

func ErrorHandler(err error, req *http.Request, res http.ResponseWriter) bool {
	superTokensInstance, insterr := getInstanceOrThrowError()
	if insterr != nil {
		panic("supertokens not initialised - " + insterr.Error())
	}
	return superTokensInstance.ErrorHandler(err, req, res)
}
