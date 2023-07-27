/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package supertokens

import (
	"net/url"
	"strings"
)

type NormalisedURLPath struct {
	value string
}

func NewNormalisedURLPath(url string) (NormalisedURLPath, error) {
	val, err := normaliseURLPathOrThrowError(url)
	if err != nil {
		return NormalisedURLPath{}, err
	}
	return NormalisedURLPath{
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
	parts := strings.Split(n.value, "/")
	if len(parts) > 2 {
		if parts[2] == "recipe" {
			return true
		}
	}
	if len(parts) > 1 {
		if parts[1] == "recipe" {
			return true
		}
	}
	return false
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
