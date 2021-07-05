package supertokens

import "net/http"

func Middleware() func(req *http.Request, res http.ResponseWriter, next http.HandlerFunc) {
	return func(req *http.Request, res http.ResponseWriter, next http.HandlerFunc) {
		superTokensInstance, err := getInstanceOrThrowError()
		if err != nil {
			errhandler := ErrorHandler()
			errhandler(err, req, res, next)
		}
		middleware := superTokensInstance.Middleware(next)
		middleware(res, req)
	}
}

func ErrorHandler() func(err error, req *http.Request, res http.ResponseWriter, next http.HandlerFunc) {
	return func(err error, req *http.Request, res http.ResponseWriter, next http.HandlerFunc) {
		superTokensInstance, err := getInstanceOrThrowError()
		if err != nil {
			errhandler := ErrorHandler()
			errhandler(err, req, res, next)
		}
		errhandler := superTokensInstance.ErrorHandler(err)
		errhandler(err, req, res, next)
	}
}

func GetAllCORSHeaders() []string {
	superTokensInstance, _ := getInstanceOrThrowError()
	return superTokensInstance.GetAllCORSHeaders()
}
