package supertokens

import (
	"errors"
	"net/url"
	"strings"
)

type NormalisedURLDomain struct {
	value string
}

func NewNormalisedURLDomain(url string, ignoreProtocol bool) (*NormalisedURLDomain, error) {
	val, err := NormaliseURLDomainOrThrowError(url, ignoreProtocol)
	if err != nil {
		return nil, err
	}
	return &NormalisedURLDomain{
		value: val,
	}, nil
}

func (n NormalisedURLDomain) GetAsStringDangerous() string {
	return n.value
}

func NormaliseURLDomainOrThrowError(input string, ignoreProtocol bool) (string, error) {
	input = strings.ToLower(strings.Trim(input, ""))

	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") && !strings.HasPrefix(input, "supertokens://") {
		if strings.HasPrefix(input, "/") {
			return "", errors.New("Please provide a valid domain name")
		}
		if strings.HasPrefix(input, ".") {
			input = input[1:]
		}
		if (strings.Contains(input, ".") || strings.HasPrefix(input, "localhost")) && !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			input = "https://" + input
			return NormaliseURLDomainOrThrowError(input, true)
		}
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
		input = urlObj.Scheme + "://" + urlObj.Host
	}
	return input, nil
}
