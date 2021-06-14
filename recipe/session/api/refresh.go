package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func HandleRefreshAPI(apiImplementation models.APIImplementation, options models.APIOptions) {
	if apiImplementation.RefreshPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return
	}
	apiImplementation.RefreshPOST(options)
	supertokens.Send200Response(options.Res, nil)
}
