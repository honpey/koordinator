package dispatcher

import (
	"context"
	"encoding/json"
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/utils"
	"k8s.io/klog/v2"
)

// RuntimeDispatcher is
type RuntimeDispatcher struct {
	cm          *utils.HookServerClientManager
	hookManager *config.Manager
}

func NewRuntimeDispatcher(cm *utils.HookServerClientManager, hookManager *config.Manager) *RuntimeDispatcher {
	return &RuntimeDispatcher{
		cm:          cm,
		hookManager: hookManager,
	}
}

func (rd *RuntimeDispatcher) genHookServerRequest() {
}

func (rd *RuntimeDispatcher) dispatchInternal(ctx context.Context, hookType config.RuntimeHookType,
	client *utils.RuntimeHookClient, request interface{}) (response interface{}, err error) {

	switch hookType {
	case config.PreRunPodSandbox:
		data, _ := json.Marshal(request.(*v1alpha1.RunPodSandboxHookRequest))
		klog.Infof("runPodSandbox request %v", string(data))
		response, err = client.PreRunPodSandboxHook(ctx, request.(*v1alpha1.RunPodSandboxHookRequest))
		if err != nil {
			klog.Infof("show error info: %v %v", response, err)
		}
	case config.PreStartContainer:
		data, _ := json.Marshal(request.(*v1alpha1.ContainerResourceHookRequest))
		klog.Infof("preStartContainer request %v", string(data))
		response, err = client.PreStartContainerHook(ctx, request.(*v1alpha1.ContainerResourceHookRequest))
		if err != nil {
			klog.Infof("show error info: %v %v", response, err)
		}
	case config.PreUpdateContainerResources:
		data, _ := json.Marshal(request.(*v1alpha1.ContainerResourceHookRequest))
		klog.Infof("preUpdateContainer request %v", string(data))
		response, err = client.PreUpdateContainerResourcesHook(ctx, request.(*v1alpha1.ContainerResourceHookRequest))
		if err != nil {
			klog.Infof("show error info: %v %v", response, err)
		}
	case config.PostStartContainer:
		data, _ := json.Marshal(request.(*v1alpha1.ContainerResourceHookRequest))
		klog.Infof("postStartContainer request %v", string(data))
		response, err = client.PostStartContainerHook(ctx, request.(*v1alpha1.ContainerResourceHookRequest))
		if err != nil {
			klog.Infof("show error info: %v %v", response, err)
		}
	case config.PostStopContainer:
		data, _ := json.Marshal(request.(*v1alpha1.ContainerResourceHookRequest))
		klog.Infof("postStopContainer request %v", string(data))
		response, err = client.PostStopContainerHook(ctx, request.(*v1alpha1.ContainerResourceHookRequest))
		if err != nil {
			klog.Infof("show error info: %v %v", response, err)
		}
	}
	return nil, nil
}

func (rd *RuntimeDispatcher) Dispatch(ctx context.Context, runtimeRequestPath config.RuntimeRequestPath,
	stage config.RuntimeHookStage, request interface{}) (interface{}, error) {
	hookServers := rd.hookManager.GetAllHook()
	for _, hookServer := range hookServers {
		for _, hookType := range hookServer.RuntimeHooks {
			if !hookType.OccursOn(runtimeRequestPath) {
				continue
			}
			if hookType.HookStage() != stage {
				continue
			}
			client, err := rd.cm.RuntimeHookClient(utils.HookServerPath{
				Path: hookServer.RemoteEndpoint,
			})
			if err != nil {
				klog.Infof("fail to create the client %v", err)
				continue
			}
			response, err := rd.dispatchInternal(ctx, hookType, client, request)
			// TODO: multi
			klog.V(6).Infof("%v %v", response, err)
		}
	}
	return nil, nil
}
