package interceptor

import (
	"context"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/resource-executor"
	"k8s.io/klog/v2"
	"net"
	"time"

	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (ci *CriInterceptor) getRuntimeHookInfo(serviceType RuntimeServiceType) (config.RuntimeRequestPath,
	resource_executor.RuntimeResourceType) {
	switch serviceType {
	case RunPodSandbox:
		return config.RunPodSandbox, resource_executor.RuntimePodResource
	case CreateContainer:
		// No Nook point in create container, but we need store the container info during container create
		return config.NoneRuntimeHookPath, resource_executor.RuntimeContainerResource
	case StartContainer:
		return config.StartContainer, resource_executor.RuntimeContainerResource
	case StopContainer:
		return config.StopContainer, resource_executor.RuntimeContainerResource
	case UpdateContainerResources:
		return config.UpdateContainerResources, resource_executor.RuntimeContainerResource
	}
	return config.NoneRuntimeHookPath, resource_executor.RuntimeNoopResource
}

func (ci *CriInterceptor) interceptRuntimeRequest(serviceType RuntimeServiceType,
	ctx context.Context, request interface{}, handler grpc.UnaryHandler) (interface{}, error) {

	runtimeHookPath, runtimeResourceType := ci.getRuntimeHookInfo(serviceType)
	resourceExecutor := resource_executor.NewOnetimeRuntimeResourceExecutor(runtimeResourceType, ci.MetaManager)

	if err := resourceExecutor.ParseRequest(request); err != nil {
		klog.Errorf("fail to parse request %v %v", request, err)
	}

	// pre call hook server
	ci.dispatcher.Dispatch(ctx, runtimeHookPath, config.PreHook, resourceExecutor.GenerateHookRequest())

	// call the backend runtime engine
	res, err := handler(ctx, request)
	if err == nil {
		klog.Infof("%v %v success", res, err)
		// store checkpoint info basing request
		// checkpoint only when response success
		if err := resourceExecutor.ResourceCheckPoint(res); err != nil {
			klog.Errorf("fail to checkpoint %v", err)
		}
	} else {
		klog.Errorf("%v %v", res, err)
	}

	// post call hook server
	ci.dispatcher.Dispatch(ctx, runtimeHookPath, config.PostHook, resourceExecutor.GenerateHookRequest())

	return res, err
}

func dialer(ctx context.Context, addr string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "unix", addr)
}

func (ci *CriInterceptor) Init(sockPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, sockPath, grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		klog.Infof("err to create  %v\n", err)
		return
	}
	ci.runtimeClient = runtimeapi.NewRuntimeServiceClient(conn)
}
