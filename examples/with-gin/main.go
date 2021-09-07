package main

import (
	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/examples/with-gin/logger"
	"github.com/supertokens/supertokens-golang/examples/with-gin/server"
)

func main() {
	config.Init()
	logger.Init()
	server.Init()
}
