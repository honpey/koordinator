package resource_executor

import (
	meta "github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
)

type RuntimeResourceExecutor interface {
	GenerateResourceCheckpoint() interface{}
	GenerateHookRequest() interface{}
	ParseRequest(request interface{}) error
	ResourceCheckPoint(response interface{}) error
}

type RuntimeResourceType string

const (
	RuntimePodResource       RuntimeResourceType = "RuntimePodResource"
	RuntimeContainerResource RuntimeResourceType = "RuntimeContainerResource"
	RuntimeNoopResource      RuntimeResourceType = "RuntimeNoopResource"
)

func NewOnetimeRuntimeResourceExecutor(runtimeResourceType RuntimeResourceType, store *meta.MetaManager) RuntimeResourceExecutor {

	switch runtimeResourceType {
	case RuntimePodResource:
		return &PodResourceExecutor{
			store: store,
		}
	case RuntimeContainerResource:
		return &ContainerResourceExecutor{
			store: store,
		}
	}
	return &NoopResourceExecutor{}
}
