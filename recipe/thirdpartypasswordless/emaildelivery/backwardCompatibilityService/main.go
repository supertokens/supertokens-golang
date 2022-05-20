package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	passwordlessBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl tplmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext), sendPasswordlessLoginEmail func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error) emaildelivery.EmailDeliveryInterface {
	if sendEmailVerificationEmail != nil {
		inputSendEmailVerificationEmail := sendEmailVerificationEmail
		sendEmailVerificationEmail = func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
			userInfo, err := (*recipeInterfaceImpl.GetUserByID)(user.ID, userContext)
			if err != nil {
				return // FIXME: No error handling here
			}
			if userInfo.Email == nil {
				return // FIXME: No error handling here
			}
			inputSendEmailVerificationEmail(evmodels.User{ID: userInfo.ID, Email: *userInfo.Email}, emailVerificationURLWithToken, userContext)
		}
	}
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
