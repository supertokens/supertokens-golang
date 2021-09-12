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

package haberdasher

// Go generate directives are a convenient way to describe compilation of proto
// files. They let users run 'go generate ./...', no Makefile necessary.
//
// This particular particular one invokes protoc using '$GOPATH/src'.
//
// This is used to tell protoc where to look up .proto files, through
// --proto_path.
//
// It is also used to tell protoc where to put output generated files, through
// --twirp_out and --go_out.

//go:generate protoc --proto_path=$GOPATH/src --twirp_out=$GOPATH/src --go_out=$GOPATH/src github.com/twitchtv/twirp-example/rpc/haberdasher/haberdasher.proto
