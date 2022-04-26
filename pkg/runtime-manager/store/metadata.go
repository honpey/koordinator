package store

import "github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"

// PodSandboxCheckpoint is almost the same with v1alpha.RunPodSandboxHookRequest
type PodSandboxCheckpoint struct {
	PodMeta        *v1alpha1.PodSandboxMetadata
	RuntimeHandler string
	Labels         map[string]string
	Annotations    map[string]string
	CgroupParent   string
	Overhead       *v1alpha1.LinuxContainerResources
	Resources      *v1alpha1.LinuxContainerResources
}

// ContainerCheckpoint is almost the same with v1alpha.ContainerResourceHookRequest
type ContainerCheckpoint struct {
	PodMeta              *v1alpha1.PodSandboxMetadata
	ContainerMata        *v1alpha1.ContainerMetadata
	ContainerAnnotations map[string]string
	ContainerResources   *v1alpha1.LinuxContainerResources
	PodResources         *v1alpha1.LinuxContainerResources
}
