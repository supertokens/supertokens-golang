package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
)

func EmailVerify(apiImplementation schema.APIInterface, options schema.APIOptions) error {
	var result map[string]string
	if options.Req.Method == http.MethodPost {
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
			return errors.New("Please provide the email verification token")
		}
		if reflect.ValueOf(token).Kind() != reflect.String {
			return errors.New("The email verification token must be a string")
		}

		response := apiImplementation.VerifyEmailPOST(token.(string), options)
		if response["status"] == "OK" {
			result["status"] = "OK"
		} else {
			for k, v := range response {
				result[k] = v.(string)
			}
		}
	} else {
		response := apiImplementation.IsEmailVerifiedGET(options)
		for k, v := range response {
			result[k] = v.(string)
		}
	}
	// todo: send200Response(options.res, result);
	return nil
}
