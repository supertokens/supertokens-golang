package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/examples/with-gin/logger"
	"github.com/supertokens/supertokens-golang/examples/with-gin/server"
)

func main() {
	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	logger.Init()
	server.Init()
}
