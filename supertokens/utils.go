package supertokens

import (
	"errors"
	"regexp"
)

func IsAnIPAddress(ipaddress string) (bool, error) {
	return regexp.MatchString(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, ipaddress)
}

func NormaliseInputAppInfoOrThrowError(appInfo AppInfo) (NormalisedAppinfo, error) {
	if appInfo == nil {
		return NormalisedAppinfo{}, errors.New("Please provide the appInfo object when calling supertokens.init")
	}
	if appInfo.apiDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your apiDomain inside the appInfo object when calling supertokens.init")
	}
	if appInfo.appName == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your appName inside the appInfo object when calling supertokens.init")
	}
	if appInfo.websiteDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
	}
	return NormalisedAppinfo{}, nil
}
