package supertokens

import (
	"net/url"
	"strings"
)

type NormalisedURLPath struct {
	value string
}

func NewNormalisedURLPath(url string) (*NormalisedURLPath, error) {
	val, err := normaliseURLPathOrThrowError(url)
	if err != nil {
		return nil, err
	}
	return &NormalisedURLPath{
		value: val,
	}, nil
}

func (n NormalisedURLPath) GetAsStringDangerous() string {
	return n.value
}

func (n NormalisedURLPath) StartsWith(other NormalisedURLPath) bool {
	return strings.HasPrefix(n.value, other.value)
}

func (n NormalisedURLPath) AppendPath(other NormalisedURLPath) NormalisedURLPath {
	return NormalisedURLPath{value: n.value + other.value}
}

func (n NormalisedURLPath) Equals(other NormalisedURLPath) bool {
	return n.value == other.value
}

func (n NormalisedURLPath) IsARecipePath() bool {
	return n.value == "/recipe" || strings.HasPrefix(n.value, "/recipe/")
}

func normaliseURLPathOrThrowError(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		if (domainGiven(input) || strings.HasPrefix(input, "localhost")) &&
			!strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			input = "http://" + input
			return normaliseURLPathOrThrowError(input)
		}

		if !strings.HasPrefix(input, "/") {
			input = "/" + input
		}
		return normaliseURLPathOrThrowError("http://example.com" + input)
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	input = urlObj.Path
	input = strings.TrimSuffix(input, "/")

	return input, nil
}

func domainGiven(input string) bool {
	// If no dot, return false.
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "/") {
		return false
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return true
	}
	if urlObj.Hostname() == "" {
		urlObj, err = url.Parse("http://" + input)
		if err != nil {
			return false
		}
	}
	return strings.Contains(urlObj.Hostname(), ".")
}
