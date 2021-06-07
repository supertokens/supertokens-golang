package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) schema.TypeNormalisedInput {
	getEmailVerificationURL := *config.GetEmailVerificationURL
	if config.GetEmailVerificationURL == nil {
		getEmailVerificationURL = GetEmailVerificationURL(appInfo)
	}
	createAndSendCustomEmail := *config.CreateAndSendCustomEmail
	if config.CreateAndSendCustomEmail == nil {
		createAndSendCustomEmail = CreateAndSendCustomEmail(appInfo)
	}

	override := schema.Override{
		Functions: config.Override.Functions,
		APIs:      config.Override.APIs,
	}

	getEmailForUserId := config.GetEmailForUserID
	return schema.TypeNormalisedInput{
		GetEmailVerificationURL:  getEmailVerificationURL,
		GetEmailForUserID:        getEmailForUserId,
		CreateAndSendCustomEmail: createAndSendCustomEmail,
		Override:                 override,
	}
}
