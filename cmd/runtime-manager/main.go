package main

import (
	"fmt"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/dispatcher"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/interceptor"
	"time"
)

func main() {
	fmt.Printf("Enter the runtime manager\n")
	interceptor.NewCriInterceptor(&dispatcher.RuntimeDispatcher{}).Setup3()
	/*
		client, _ := tools.NewRuntimeHookClient(pkg..DefaultRpcServerPath)
		go client.Start()
	*/
	config := config.NewConfigManager()
	config.Setup()
	for {
		time.Sleep(100 * time.Minute)
	}

}
