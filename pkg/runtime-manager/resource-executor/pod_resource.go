package resource_executor

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"

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

func (p *PodResourceExecutor) String() string {
	return fmt.Sprintf("%v/%v", p.PodMeta.Name, p.PodMeta.Uid)
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
	klog.Infof("success parse pod Info %v during pod run", p)
	return nil
}

func (p *PodResourceExecutor) ResourceCheckPoint(response interface{}) error {
	runPodSandboxResponse, ok := response.(*runtimeapi.RunPodSandboxResponse)
	if !ok {
		return fmt.Errorf("bad response %v", response)
	}
	podCheckPoint := p.GenerateResourceCheckpoint().(*meta.PodSandboxCheckpoint)
	data, _ := json.Marshal(podCheckPoint)
	err := p.store.WritePodSandboxCheckpoint(runPodSandboxResponse.PodSandboxId, podCheckPoint)
	if err != nil {
		return err
	}
	klog.Infof("success to checkpoint pod level info %v %v",
		runPodSandboxResponse.PodSandboxId, string(data))
	return nil
}
