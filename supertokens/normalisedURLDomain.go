package supertokens

import (
	"errors"
	"net/url"
	"strings"
)

type NormalisedURLDomain struct {
	Value string
}

func (n NormalisedURLDomain) GetAsStringDangerous() string {
	return n.Value
}

func NormaliseURLDomainOrThrowError(input string, ignoreProtocol bool) (string, error) {
	input = strings.ToLower(strings.Trim(input, ""))

	if strings.HasPrefix(input, "http://") != true && strings.HasPrefix(input, "https://") != true && strings.HasPrefix(input, "supertokens://") != true {
		return "", errors.New("converting to proper URL")
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	if ignoreProtocol {
		isAnIP, err := IsAnIPAddress(urlObj.Host)
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(urlObj.Host, "localhost") || isAnIP {
			input = "http://" + urlObj.Host
		} else {
			input = "https://" + urlObj.Host
		}
	} else {
		input = urlObj.Scheme + "//" + urlObj.Host
	}
	return input, nil
}
