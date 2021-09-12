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

package hooks

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/twitchtv/twirp"
)

var ctxKey = new(int)

// LoggingHooks creates a new twirp.ServerHooks which logs requests as they are
// routed to Twirp, and logs responses (including response time) when they are
// sent.
//
// This is a demonstration showing how you can use context accessors with hooks.
func LoggingHooks(w io.Writer) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			startTime := time.Now()
			ctx = context.WithValue(ctx, ctxKey, startTime)
			return ctx, nil
		},
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			svc, _ := twirp.ServiceName(ctx)
			meth, _ := twirp.MethodName(ctx)
			fmt.Fprintf(w, "received req svc=%q method=%q\n", svc, meth)
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			startTime := ctx.Value(ctxKey).(time.Time)
			svc, _ := twirp.ServiceName(ctx)
			meth, _ := twirp.MethodName(ctx)
			fmt.Fprintf(w, "response sent svc=%q method=%q time=%q\n", svc, meth, time.Since(startTime))
		},
	}
}
