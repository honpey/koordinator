package resource_executor

import (
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/store"
)

type NoopResourceExecutor struct {
	store *store.MetaManager
}

func (n *NoopResourceExecutor) GenerateResourceCheckpoint() interface{} {
	return v1alpha1.ContainerResourceHookRequest{}
}

func (n *NoopResourceExecutor) GenerateHookRequest() interface{} {
	return store.ContainerCheckpoint{}
}

func (n *NoopResourceExecutor) ParseRequest(request interface{}) error {
	return nil
}
func (n *NoopResourceExecutor) ResourceCheckPoint(response interface{}) error {

	return nil
}
