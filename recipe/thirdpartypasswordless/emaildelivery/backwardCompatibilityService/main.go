package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	passwordlessBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(appInfo supertokens.NormalisedAppinfo, sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext), sendPasswordlessLoginEmail func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error) emaildelivery.EmailDeliveryInterface {

	emailVerificationService := emailVerificationBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail)
	passwordlessService := passwordlessBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendPasswordlessLoginEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationService.SendEmail)(input, userContext)

		} else if input.PasswordlessLogin != nil {
			return (*passwordlessService.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
