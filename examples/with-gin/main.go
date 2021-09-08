package main

import (
	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/examples/with-gin/server"
)

func main() {
	// supertokens init here
	config.Init()

	// adding of superotkens middleware, cors and APIs
	server.Init()
}
