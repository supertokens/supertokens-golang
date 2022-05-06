package backwardcompatibility

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(appInfo supertokens.NormalisedAppinfo, createAndSendCustomEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext)) emaildeliverymodels.EmailDeliveryInterface {
	sendEmail := func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			createAndSendCustomEmail(evmodels.User{
				ID:    input.EmailVerification.User.ID,
				Email: input.EmailVerification.User.Email,
			}, input.EmailVerification.EmailVerifyLink, userContext)
		} else {
			return errors.New("should never come here")
		}
		return nil
	}

	return emaildeliverymodels.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
