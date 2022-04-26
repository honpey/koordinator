package resource_executor

import (
	"encoding/json"
	"fmt"

	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
)

type ContainerResourceExecutor struct {
	store                *store.MetaManager
	PodMeta              *v1alpha1.PodSandboxMetadata
	ContainerMata        *v1alpha1.ContainerMetadata
	ContainerAnnotations map[string]string
	ContainerResources   *v1alpha1.LinuxContainerResources
	PodResources         *v1alpha1.LinuxContainerResources
}

func (c *ContainerResourceExecutor) String() string {
	return fmt.Sprintf("pod(%v/%v)container(%v)",
		c.PodMeta.Name, c.PodMeta.Uid,
		c.ContainerMata.Name)
}

func (c *ContainerResourceExecutor) GenerateResourceCheckpoint() interface{} {
	return &store.ContainerCheckpoint{
		PodMeta:              c.PodMeta,
		ContainerMata:        c.ContainerMata,
		ContainerAnnotations: c.ContainerAnnotations,
		ContainerResources:   c.ContainerResources,
		PodResources:         c.PodResources,
	}
}

func (c *ContainerResourceExecutor) GenerateHookRequest() interface{} {
	return &v1alpha1.ContainerResourceHookRequest{
		PodMeta:              c.PodMeta,
		ContainerMata:        c.ContainerMata,
		ContainerAnnotations: c.ContainerAnnotations,
		ContainerResources:   c.ContainerResources,
		PodResources:         c.PodResources,
	}
}

func (c *ContainerResourceExecutor) updateByCheckPoint(containerID string) error {
	containerCheckPoint, err := c.store.GetContainerCheckpoint(containerID)
	if err != nil {
		return err
	}
	c.PodMeta = containerCheckPoint.PodMeta
	c.PodResources = containerCheckPoint.PodResources
	c.ContainerMata = containerCheckPoint.ContainerMata
	c.ContainerResources = containerCheckPoint.ContainerResources
	c.ContainerAnnotations = containerCheckPoint.ContainerAnnotations
	klog.Infof("get container info successful %v", containerID)
	return nil
}

func (c *ContainerResourceExecutor) ParseRequest(request interface{}) error {
	switch request.(type) {
	case *runtimeapi.CreateContainerRequest:
		// get the pod info from local store
		createContainerRequest := request.(*runtimeapi.CreateContainerRequest)
		podID := createContainerRequest.PodSandboxId
		podCheckPoint, err := c.store.GetPodSandboxCheckpoint(podID)
		if err != nil {
			return err
		}
		c.PodMeta = podCheckPoint.PodMeta
		c.PodResources = podCheckPoint.Resources

		// construct the container info
		c.ContainerMata = &v1alpha1.ContainerMetadata{
			Name:    createContainerRequest.GetConfig().GetMetadata().GetName(),
			Attempt: createContainerRequest.GetConfig().GetMetadata().GetAttempt(),
		}
		c.ContainerAnnotations = createContainerRequest.GetConfig().GetAnnotations()
		c.ContainerResources = transferResource(createContainerRequest.GetConfig().GetLinux().GetResources())
		klog.Infof("success parse container info %v during container create", c)
	case *runtimeapi.StartContainerRequest:
		startContainerRequest := request.(*runtimeapi.StartContainerRequest)
		containerID := startContainerRequest.ContainerId
		err := c.updateByCheckPoint(containerID)
		if err != nil {
			return err
		}
		klog.Infof("success parse container Info %v during container start", c)
	case *runtimeapi.UpdateContainerResourcesRequest:
		updateContainerResourcesRequest := request.(*runtimeapi.UpdateContainerResourcesRequest)
		containerID := updateContainerResourcesRequest.ContainerId
		err := c.updateByCheckPoint(containerID)
		if err != nil {
			return err
		}
		klog.Infof("success parse container Info %v during container update resource", c)
	}
	return nil
}

func (c *ContainerResourceExecutor) ResourceCheckPoint(response interface{}) error {
	// container level resource checkpoint would be triggered during post container create only
	createContainer, ok := response.(*runtimeapi.CreateContainerResponse)
	if !ok {
		return nil
	}
	containerCheckpoint := c.GenerateResourceCheckpoint().(*store.ContainerCheckpoint)
	data, _ := json.Marshal(containerCheckpoint)

	err := c.store.WriteContainerCheckpoint(createContainer.ContainerId, containerCheckpoint)
	if err != nil {
		return err
	}
	klog.Infof("success to checkpoint container level info %v %v",
		createContainer.ContainerId, string(data))
	return nil

}
