package api

import (
	"encoding/json"
	defaultErrors "errors"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
)

func validateFormFieldsOrThrowError(configFormFields []epmodels.NormalisedFormField, formFieldsRaw []interface{}) ([]epmodels.TypeFormField, error) {
	if formFieldsRaw == nil {
		return nil, defaultErrors.New("Missing input param: formFields")
	}

	if len(formFieldsRaw) == 0 {
		return nil, defaultErrors.New("formFields must be an array")
	}

	var formFields []epmodels.TypeFormField
	for _, rawFormField := range formFieldsRaw {
		jsonformField, err := json.Marshal(rawFormField)
		if err != nil {
			return nil, err
		}
		var formField epmodels.TypeFormField
		err = json.Unmarshal(jsonformField, &formField)
		if err != nil {
			return nil, err
		}

		if formField.ID == "email" {
			formFields = append(formFields, epmodels.TypeFormField{
				ID:    formField.ID,
				Value: strings.TrimSpace(formField.Value),
			})
		} else {
			formFields = append(formFields, epmodels.TypeFormField{
				ID:    formField.ID,
				Value: formField.Value,
			})
		}
	}

	return formFields, validateFormOrThrowError(configFormFields, formFields)
}

func validateFormOrThrowError(configFormFields []epmodels.NormalisedFormField, inputs []epmodels.TypeFormField) error {
	var validationErrors []errors.ErrorPayload
	if len(configFormFields) != len(inputs) {
		return defaultErrors.New("Are you sending too many / too few formFields?")
	}
	for _, field := range configFormFields {
		var input epmodels.TypeFormField
		for _, inputField := range inputs {
			if inputField.ID == field.ID {
				input = inputField
				break
			}
		}
		if input.Value == "" && !field.Optional {
			validationErrors = append(validationErrors, errors.ErrorPayload{ID: field.ID, Error: "Field is not optional"})
		} else {
			err := field.Validate(input.Value)
			if err != nil {
				validationErrors = append(validationErrors, errors.ErrorPayload{
					ID:    field.ID,
					Error: *err,
				})
			}
		}
	}
	if len(validationErrors) != 0 {
		return errors.FieldError{
			Msg:     "Error in input formFields",
			Payload: validationErrors,
		}
	}
	return nil
}
