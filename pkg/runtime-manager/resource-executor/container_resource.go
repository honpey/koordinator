package resource_executor

import (
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type ContainerResourceExecutor struct {
	store                *store.MetaManager
	PodMeta              *v1alpha1.PodSandboxMetadata
	ContainerMata        *v1alpha1.ContainerMetadata
	ContainerAnnotations map[string]string
	ContainerResources   *v1alpha1.LinuxContainerResources
	PodResources         *v1alpha1.LinuxContainerResources
}

func (c *ContainerResourceExecutor) GenerateResourceCheckpoint() interface{} {
	return store.ContainerCheckpoint{
		PodMeta:              c.PodMeta,
		ContainerMata:        c.ContainerMata,
		ContainerAnnotations: c.ContainerAnnotations,
		ContainerResources:   c.ContainerResources,
		PodResources:         c.PodResources,
	}
}

func (c *ContainerResourceExecutor) GenerateHookRequest() interface{} {
	return v1alpha1.ContainerResourceHookRequest{
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
			Name: createContainerRequest.GetConfig().GetMetadata().GetName(),

			Attempt: createContainerRequest.GetConfig().GetMetadata().GetAttempt(),
		}
		c.ContainerAnnotations = createContainerRequest.GetConfig().GetAnnotations()
		c.ContainerResources = transferResource(createContainerRequest.GetConfig().GetLinux().GetResources())
	case *runtimeapi.StartContainerRequest:
		startContainerRequest := request.(*runtimeapi.StartContainerRequest)
		containerID := startContainerRequest.ContainerId
		return c.updateByCheckPoint(containerID)
	case *runtimeapi.UpdateContainerResourcesRequest:
		updateContainerResourcesRequest := request.(*runtimeapi.UpdateContainerResourcesRequest)
		containerID := updateContainerResourcesRequest.ContainerId
		return c.updateByCheckPoint(containerID)
	}
	return nil
}

func (c *ContainerResourceExecutor) ResourceCheckPoint(response interface{}) error {
	// container level resource checkpoint would be triggered during post container create only
	createContainer, ok := response.(*runtimeapi.CreateContainerResponse)
	if !ok {
		return nil
	}
	return c.store.WriteContainerCheckpoint(createContainer.ContainerId,
		c.GenerateResourceCheckpoint().(*store.ContainerCheckpoint))
}
