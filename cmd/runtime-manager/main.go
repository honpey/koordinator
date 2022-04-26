package main

import (
	"fmt"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/dispatcher"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/interceptor"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/tools"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/utils"
	"time"
)

func main() {
	fmt.Printf("Enter the runtime manager\n")
	go tools.NewMockRuntimeHookServer()
	config := config.NewConfigManager()
	config.Setup()
	interceptor.NewCriInterceptor(dispatcher.NewRuntimeDispatcher(
		utils.NewClientManager(), config)).Setup()
	/*
		client, _ := tools.NewRuntimeHookClient(pkg..DefaultRpcServerPath)
		go client.Start() */
	for {
		time.Sleep(100 * time.Minute)
	}
}
