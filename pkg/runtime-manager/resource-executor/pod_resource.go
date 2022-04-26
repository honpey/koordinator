package resource_executor

import (
	"fmt"
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
	meta "github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type PodResourceExecutor struct {
	store          *meta.MetaManager
	PodMeta        *v1alpha1.PodSandboxMetadata
	RuntimeHandler string
	Annotations    map[string]string
	Labels         map[string]string
	CgroupParent   string
}

func (p *PodResourceExecutor) GenerateResourceCheckpoint() interface{} {
	return &store.PodSandboxCheckpoint{
		PodMeta:        p.PodMeta,
		RuntimeHandler: p.RuntimeHandler,
		Annotations:    p.Annotations,
		Labels:         p.Labels,
	}
}

func (p *PodResourceExecutor) GenerateHookRequest() interface{} {
	return &v1alpha1.RunPodSandboxHookRequest{
		PodMeta:        p.PodMeta,
		RuntimeHandler: p.RuntimeHandler,
		Annotations:    p.Annotations,
		Labels:         p.Labels,
	}
}

func (p *PodResourceExecutor) ParseRequest(request interface{}) error {
	runPodSandboxRequest, ok := request.(*runtimeapi.RunPodSandboxRequest)
	if !ok {
		return fmt.Errorf("bad request")
	}
	p.PodMeta = &v1alpha1.PodSandboxMetadata{
		Name:      runPodSandboxRequest.GetConfig().GetMetadata().GetName(),
		Namespace: runPodSandboxRequest.GetConfig().GetMetadata().GetNamespace(),
	}
	p.RuntimeHandler = runPodSandboxRequest.GetRuntimeHandler()
	p.Annotations = runPodSandboxRequest.GetConfig().GetAnnotations()
	p.Labels = runPodSandboxRequest.GetConfig().GetLabels()
	p.CgroupParent = runPodSandboxRequest.GetConfig().GetLinux().GetCgroupParent()
	return nil
}

func (p *PodResourceExecutor) ResourceCheckPoint(response interface{}) error {
	runPodSandboxResponse, ok := response.(*runtimeapi.RunPodSandboxResponse)
	if !ok {
		return fmt.Errorf("bad response %v", response)
	}
	p.store.WritePodSandboxCheckpoint(runPodSandboxResponse.PodSandboxId, p.GenerateResourceCheckpoint().(*meta.PodSandboxCheckpoint))
	return nil
}
