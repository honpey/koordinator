package tools

import (
	"context"
	"fmt"
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"time"

	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/utils"
)

func GenerateMockRequest(m *utils.RuntimeHookClient) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		rsp, err := m.PreRunPodSandboxHook(ctx, &v1alpha1.RunPodSandboxHookRequest{})
		fmt.Printf("clien run presanbox hook : %v %v\n", rsp, err)
		rsp2, err := m.PreStartContainerHook(ctx, &v1alpha1.ContainerResourceHookRequest{})
		fmt.Printf("clien run container hook : %v %v\n", rsp2, err)
		cancel()
		time.Sleep(time.Second)
	}
}
