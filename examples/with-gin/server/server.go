package server

import (
	"log"

	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
)

func Init() {
	config := config.GetConfig()

	r := newRouter()
	err := r.Run(config.GetString("server.apiPort"))
	if err != nil {
		log.Println("error running server => ", err)
	}
}
