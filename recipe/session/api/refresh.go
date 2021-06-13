package api

import "github.com/supertokens/supertokens-golang/recipe/session/schema"

func HandleRefreshAPI(apiImplementation schema.APIImplementation, options schema.APIOptions) {
	if apiImplementation.RefreshPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return
	}
	apiImplementation.RefreshPOST(options)
	// TODO: send200Response(options.res, {})
}
