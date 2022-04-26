package tools

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"net"
	"os"

	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
)

const (
	socketPath = "/tmp/runtimeserver.sock"
)

type MockHookRuntimeServer struct {
	v1alpha1.UnimplementedRuntimeHookServiceServer
}

func (s *MockHookRuntimeServer) PreRunPodSandboxHook(context.Context, *v1alpha1.RunPodSandboxHookRequest) (*v1alpha1.RunPodSandboxHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreRunPodSandboxHook not implemented")
}
func (s *MockHookRuntimeServer) PreStartContainerHook(context.Context, *v1alpha1.ContainerResourceHookRequest) (*v1alpha1.ContainerResourceHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreStartContainerHook not implemented")
}
func (s *MockHookRuntimeServer) PostStartContainerHook(context.Context, *v1alpha1.ContainerResourceHookRequest) (*v1alpha1.ContainerResourceHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostStartContainerHook not implemented")
}
func (s *MockHookRuntimeServer) PostStopContainerHook(context.Context, *v1alpha1.ContainerResourceHookRequest) (*v1alpha1.ContainerResourceHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostStopContainerHook not implemented")
}
func (s *MockHookRuntimeServer) PreUpdateContainerResourcesHook(context.Context, *v1alpha1.ContainerResourceHookRequest) (*v1alpha1.ContainerResourceHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreUpdateContainerResourcesHook not implemented")
}

func NewMockRuntimeHookServer() {
	os.Remove(socketPath)
	klog.Infof("start mock hook server on %v", socketPath)
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		klog.Infof("fail to create  %v", err)
		return
	}
	srv := grpc.NewServer()
	v1alpha1.RegisterRuntimeHookServiceServer(srv, &MockHookRuntimeServer{})
	err = srv.Serve(lis)
	klog.Infof("fail to trigger mock server %v", err)
}
