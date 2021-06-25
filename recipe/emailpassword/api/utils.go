package api

import (
	"reflect"
	"strings"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	emailpasswordErrors "github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
)

func validateFormFieldsOrThrowError(configFormFields []models.NormalisedFormField, formFieldsRaw []models.FormFieldValue) ([]models.FormFieldValue, error) {
	if formFieldsRaw == nil {
		return nil, errors.BadInputError{Msg: "Missing input param: formFields"}
	}
	if reflect.TypeOf(formFieldsRaw).Kind() == reflect.Array {
		return nil, errors.BadInputError{Msg: "formFields must be an array"}
	}
	var formFields []models.FormFieldValue
	for _, formField := range formFieldsRaw {
		if formField.ID == constants.FormFieldEmailID {
			formFields = append(formFields, models.FormFieldValue{
				ID:    formField.ID,
				Value: strings.TrimSpace(formField.Value),
			})
		}
	}

	return formFields, validateFormOrThrowError(configFormFields, formFields)
}

func validateFormOrThrowError(configFormFields []models.NormalisedFormField, inputs []models.FormFieldValue) error {
	var validationErrors []emailpasswordErrors.ErrorPayload
	if len(configFormFields) != len(inputs) {
		return errors.BadInputError{Msg: "Are you sending too many / too few formFields?"}
	}
	for _, field := range configFormFields {
		var input models.FormFieldValue
		for _, inputField := range inputs {
			if inputField.ID == field.ID {
				input = inputField
				break
			}
		}
		if input.Value == "" && !field.Optional {
			validationErrors = append(validationErrors, emailpasswordErrors.ErrorPayload{ID: field.ID, Error: "Field is not optional"})
		} else {
			err := field.Validate(input.Value)
			if err != nil {
				validationErrors = append(validationErrors, struct {
					ID    string
					Error string
				}{ID: field.ID, Error: *err})
			}
		}
	}
	if len(validationErrors) != 0 {
		return emailpasswordErrors.FieldError{
			Msg:     "Error in input formFields",
			Payload: validationErrors,
		}
	}
	return nil
}
