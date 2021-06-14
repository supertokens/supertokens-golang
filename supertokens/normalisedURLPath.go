package supertokens

import (
	"fmt"
	"net/url"
	"strings"
)

type NormalisedURLPath struct {
	value string
}

func NewNormalisedURLPath(url string) (*NormalisedURLPath, error) {
	val, err := NormaliseURLPathOrThrowError(url)
	if err != nil {
		return nil, err
	}
	return &NormalisedURLPath{
		value: val,
	}, nil
}

func (n *NormalisedURLPath) GetAsStringDangerous() string {
	return n.value
}

func (n *NormalisedURLPath) StartsWith(other NormalisedURLPath) bool {
	return strings.HasPrefix(n.value, other.value)
}

func (n *NormalisedURLPath) AppendPath(other NormalisedURLPath) NormalisedURLPath {
	return NormalisedURLPath{value: n.value + other.value}
}

func (n *NormalisedURLPath) Equals(other NormalisedURLPath) bool {
	return n.value == other.value
}

func (n *NormalisedURLPath) IsARecipePath() bool {
	return n.value == "/recipe" || strings.HasPrefix(n.value, "/recipe/")
}

func NormaliseURLPathOrThrowError(input string) (string, error) {
	input = strings.ToLower(strings.Trim(input, ""))
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		fmt.Println("converting to proper URL")
		if (domainGiven(input) || strings.HasPrefix(input, "localhost")) &&
			!strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			input = "http://" + input
			return NormaliseURLPathOrThrowError(input)
		}

		if input[:1] != "/" {
			input = "/" + input
		}
		return NormaliseURLPathOrThrowError("http://example.com" + input)
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	input = urlObj.Path
	if input != "" && input[len(input)-1:] == "/" {
		return input[:len(input)-1], nil
	}

	return input, nil
}

func domainGiven(input string) bool {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "/") {
		return false
	}

	urlObj, err := url.Parse(input)
	if err != nil {
		return false
	}
	if urlObj.Host == "" {
		urlObj, err = url.Parse("http://" + input)
		if err != nil {
			return false
		}
	}
	return strings.Index(urlObj.Host, ".") != -1
}
