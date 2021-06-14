package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func EmailVerify(apiImplementation schema.APIImplementation, options schema.APIOptions) error {
	var result map[string]interface{}
	if options.Req.Method == http.MethodPost {

		if apiImplementation.VerifyEmailPOST == nil {
			options.OtherHandler(options.Res, options.Req)
			return nil
		}

		body, err := ioutil.ReadAll(options.Req.Response.Body)
		if err != nil {
			return err
		}
		var readBody map[string]interface{}
		err = json.Unmarshal(body, &readBody)
		if err != nil {
			return err
		}
		token, ok := readBody["token"]
		if !ok {
			return errors.GeneralError{Msg: "Please provide the email verification token"}
		}
		if reflect.ValueOf(token).Kind() != reflect.String {
			return errors.GeneralError{Msg: "The email verification token must be a string"}
		}

		response, err := apiImplementation.VerifyEmailPOST(token.(string), options)
		if err != nil {
			return err
		}

		if response.Ok != nil {
			result = map[string]interface{}{
				"status": "OK",
			}
		} else {
			result = map[string]interface{}{
				"status": "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR",
			}
		}
	} else {
		if apiImplementation.IsEmailVerifiedGET == nil {
			options.OtherHandler(options.Res, options.Req)
			return nil
		}

		response, err := apiImplementation.IsEmailVerifiedGET(options)
		if err != nil {
			return err
		}

		result = map[string]interface{}{
			"status":     "OK",
			"isVerified": response,
		}

	}

	supertokens.Send200Response(options.Res, result)
	return nil
}
