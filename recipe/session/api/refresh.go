package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func HandleRefreshAPI(apiImplementation schema.APIImplementation, options schema.APIOptions) {
	if apiImplementation.RefreshPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return
	}
	apiImplementation.RefreshPOST(options)
	supertokens.Send200Response(options.Res, nil)
}
