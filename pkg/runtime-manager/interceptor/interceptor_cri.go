package interceptor

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"
	"net"
	"reflect"
	"time"

	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (ci *CriInterceptor) getHookType(serviceType RuntimeServiceType) config.RuntimeRequestPath {
	switch serviceType {
	case RunPodSandbox:
		return config.RunPodSandbox
	case StartContainer:
		return config.StartContainer
	case StopContainer:
		return config.StopContainer
	case UpdateContainerResources:
		return config.UpdateContainerResources
	}
	return config.NoneRuntimeHookPath
}

func (ci *CriInterceptor) transfer(r *runtimeapi.LinuxContainerResources) *v1alpha1.LinuxContainerResources {
	linuxResource := &v1alpha1.LinuxContainerResources{
		CpuPeriod:              r.GetCpuPeriod(),
		CpuQuota:               r.GetCpuQuota(),
		CpuShares:              r.GetCpuShares(),
		MemoryLimitInBytes:     r.GetMemoryLimitInBytes(),
		OomScoreAdj:            r.GetOomScoreAdj(),
		CpusetCpus:             r.GetCpusetCpus(),
		CpusetMems:             r.GetCpusetMems(),
		Unified:                r.GetUnified(),
		MemorySwapLimitInBytes: r.GetMemorySwapLimitInBytes(),
	}

	for _, item := range r.GetHugepageLimits() {
		linuxResource.HugepageLimits = append(linuxResource.HugepageLimits, &v1alpha1.HugepageLimit{
			PageSize: item.GetPageSize(),
			Limit:    item.GetLimit(),
		})
	}

	return linuxResource
}

func (ci *CriInterceptor) generateHookRequest(request interface{},
	hookRequestPath config.RuntimeRequestPath) (interface{}, error) {
	switch hookRequestPath {
	case config.RunPodSandbox:
		runPodSandboxRequest, ok := request.(*runtimeapi.RunPodSandboxRequest)
		if !ok {
			klog.Errorf("fail to transfer %v to runtimeapi.RunPodSandboxRequest", reflect.TypeOf(request))
			return nil, fmt.Errorf("hh")
		}
		return &v1alpha1.RunPodSandboxHookRequest{
			PodMeta: &v1alpha1.PodSandboxMetadata{
				Name:      runPodSandboxRequest.GetConfig().GetMetadata().GetName(),
				Namespace: runPodSandboxRequest.GetConfig().GetMetadata().GetNamespace(),
			},
			RuntimeHandler: runPodSandboxRequest.GetRuntimeHandler(),
			Annotations:    runPodSandboxRequest.GetConfig().GetAnnotations(),
			Labels:         runPodSandboxRequest.GetConfig().GetLabels(),
			CgroupParent:   runPodSandboxRequest.GetConfig().GetLinux().GetCgroupParent(),
		}, nil
	case config.StartContainer:
		_, ok := request.(*runtimeapi.StartContainerRequest)
		if !ok {
			klog.Errorf("fail to transfer %v to runtimeapi.RunPodSandboxRequest", reflect.TypeOf(request))
			return nil, fmt.Errorf("hh")
		}
		return &v1alpha1.ContainerResourceHookRequest{
			// TODO: add the pod/containerd  info when container create
		}, nil

	case config.StopContainer:
		_, ok := request.(*runtimeapi.StopContainerRequest)
		if !ok {

			klog.Errorf("fail to transfer %v to runtimeapi.RunPodSandboxRequest", reflect.TypeOf(request))
			return nil, fmt.Errorf("hh")
		}
		return &v1alpha1.ContainerResourceHookRequest{}, nil
	case config.UpdateContainerResources:
		_, ok := request.(*runtimeapi.UpdateContainerResourcesRequest)
		if !ok {
			klog.Errorf("fail to transfer %v to runtimeapi.RunPodSandboxRequest", reflect.TypeOf(request))
			return nil, fmt.Errorf("hh")
		}
		return &v1alpha1.ContainerResourceHookRequest{}, nil
	}
	return nil, fmt.Errorf("fail")
}

func (ci *CriInterceptor) interceptRuntimeRequest(serviceType RuntimeServiceType,
	ctx context.Context, request interface{}, handler grpc.UnaryHandler) (interface{}, error) {

	requestPath := ci.getHookType(serviceType)

	if preHookType := requestPath.PreHookType(); preHookType != config.NoneRuntimeHookType {
		if hookRequest, err := ci.generateHookRequest(request, requestPath); err != nil {
			ci.dispatcher.Dispatch(ctx, requestPath, hookRequest)
		} else {
			klog.Errorf("fail to create the")
		}
	}

	res, err := handler(ctx, request)
	klog.Infof("%v %v", res, err)
	// should record the pod infoif

	if postHookType := requestPath.PostHookType(); postHookType != config.NoneRuntimeHookType {
		if hookRequest, err := ci.generateHookRequest(request, requestPath); err != nil {
			ci.dispatcher.Dispatch(ctx, requestPath, hookRequest)
		} else {
			klog.Infof("fail to create the")

		}
	}
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
