package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func HandleRefreshAPI(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.RefreshPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	err := apiImplementation.RefreshPOST(options)
	if err != nil {
		return err
	}
	supertokens.Send200Response(options.Res, nil)
	return nil
}
