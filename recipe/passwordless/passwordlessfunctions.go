package passwordless

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

var PasswordlessLoginEmailSentForTest bool = false

func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	return func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
		if supertokens.IsRunningInTestMode() {
			PasswordlessLoginEmailSentForTest = true
			// if running in test mode, we do not want to send this.
			return nil
		}

		// FIXME: What to do here!!??
		return nil
	}
}
