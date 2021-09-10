package supertokens

import (
	"errors"
	"net/url"
	"strings"
)

type NormalisedURLDomain struct {
	value string
}

func NewNormalisedURLDomain(url string) (NormalisedURLDomain, error) {
	val, err := normaliseURLDomainOrThrowError(url, false)
	if err != nil {
		return NormalisedURLDomain{}, err
	}
	return NormalisedURLDomain{
		value: val,
	}, nil
}

func (n NormalisedURLDomain) GetAsStringDangerous() string {
	return n.value
}

func normaliseURLDomainOrThrowError(input string, ignoreProtocol bool) (string, error) {
	input = strings.ToLower(strings.Trim(input, ""))

	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") && !strings.HasPrefix(input, "supertokens://") {
		if strings.HasPrefix(input, "/") {
			return "", errors.New("please provide a valid domain name")
		}
		input = strings.TrimPrefix(input, ".")
		if (strings.Contains(input, ".") || strings.HasPrefix(input, "localhost")) && !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			input = "https://" + input
			return normaliseURLDomainOrThrowError(input, true)
		}
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	if ignoreProtocol {
		isAnIP, err := IsAnIPAddress(urlObj.Hostname())
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
