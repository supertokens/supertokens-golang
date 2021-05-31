package emailverification

import "github.com/supertokens/supertokens-golang/supertokens"

func getEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) string {
	return appInfo.WebsiteDomain.Value + appInfo.WebsiteBasePath.Value + "/verify-email"
}

func createAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) { //doubt

}
