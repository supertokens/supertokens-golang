// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package haberdasherserver

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/supertokens/supertokens-golang/examples/with-twirp/haberdasher"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/twitchtv/twirp"
)

// New returns a new Haberdasher which returns random Hats of the requested
// size.
func New() *randomHaberdasher {
	return new(randomHaberdasher)
}

// randomHaberdasher is our implementation of the generated
// rpc/haberdasher.Haberdasher interface. This is where the real "business
// logic" lives.
type randomHaberdasher struct{}

func (h *randomHaberdasher) MakeHat(ctx context.Context, size *haberdasher.Size) (*haberdasher.Hat, error) {
	sessionContainer := session.GetSessionFromRequestContext(ctx)
	if sessionContainer == nil {
		fmt.Println("no session exists!")
	} else {
		// session exists!
		fmt.Println(sessionContainer.GetUserID())
	}
	// When returning an error, it's best to use the error constructors defined in
	// the Twirp package so that the client gets a well-structured error response.
	if size.Inches <= 0 {
		return nil, twirp.InvalidArgumentError("Inches", "I can't make a hat that small!")
	}
	colors := []string{"white", "black", "brown", "red", "blue"}
	names := []string{"bowler", "baseball cap", "top hat", "derby"}
	return &haberdasher.Hat{
		Size:  size.Inches,
		Color: colors[rand.Intn(len(colors))],
		Name:  names[rand.Intn(len(names))],
	}, nil
}
