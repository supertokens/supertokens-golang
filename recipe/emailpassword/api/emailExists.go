package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func EmailExists(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.EmailExistsGET == nil {
		// TODO: add tests - does their actual API get called?
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	email := options.Req.URL.Query().Get("email")
	if email == "" {
		return supertokens.BadInputError{Msg: "Please provide the email as a GET param"}
	}
	result, err := apiImplementation.EmailExistsGET(email, options)
	if err != nil {
		return err
	}
	supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
		"exists": result.OK.Exists,
	})

	return nil
}
