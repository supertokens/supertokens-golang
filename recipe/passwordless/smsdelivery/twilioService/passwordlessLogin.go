package twilioService

import (
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func getPasswordlessLoginSmsContent(input smsdelivery.PasswordlessLoginType) smsdelivery.TwilioGetContentResult {
	return smsdelivery.TwilioGetContentResult{
		Body:          getPasswordlessLoginSmsBody(input.CodeLifetime, input.UrlWithLinkCode, input.UserInputCode),
		ToPhoneNumber: input.PhoneNumber,
	}
}

func getPasswordlessLoginSmsBody(codeLifetime uint64, urlWithLinkCode *string, userInputCode *string) string {
	var message string = ""

	if urlWithLinkCode != nil && userInputCode != nil {
		message = `Enter OTP: ` + *userInputCode + ` OR click this link: ` + *urlWithLinkCode + ` to login.`
	} else if urlWithLinkCode != nil {
		message = `Click this link: ` + *urlWithLinkCode + ` to login.`
	} else {
		message = `Enter OTP: ` + *userInputCode + ` to login.`
	}
	message += ` It will expire in ` + supertokens.HumaniseMilliseconds(codeLifetime) + `.`
	return message
}
