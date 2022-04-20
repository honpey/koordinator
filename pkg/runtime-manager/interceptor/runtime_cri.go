package interceptor

import (
	"context"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (ci *CriInterceptor) Version(ctx context.Context, req *runtimeapi.VersionRequest) (*runtimeapi.VersionResponse, error) {
	return ci.runtimeClient.Version(ctx, req)
}

func (ci *CriInterceptor) RunPodSandbox(ctx context.Context, req *runtimeapi.RunPodSandboxRequest) (*runtimeapi.RunPodSandboxResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(RunPodSandbox, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.RunPodSandbox(ctx, req.(*runtimeapi.RunPodSandboxRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.RunPodSandboxResponse), err
}
func (ci *CriInterceptor) StopPodSandbox(ctx context.Context, req *runtimeapi.StopPodSandboxRequest) (*runtimeapi.StopPodSandboxResponse, error) {

	rsp, err := ci.interceptRuntimeRequest(StopPodSandbox, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.StopPodSandbox(ctx, req.(*runtimeapi.StopPodSandboxRequest))
		})

	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.StopPodSandboxResponse), err
}

func (ci *CriInterceptor) RemovePodSandbox(ctx context.Context, req *runtimeapi.RemovePodSandboxRequest) (*runtimeapi.RemovePodSandboxResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(RemovePodSandbox, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.RemovePodSandbox(ctx, req.(*runtimeapi.RemovePodSandboxRequest))

		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.RemovePodSandboxResponse), err
}

func (ci *CriInterceptor) PodSandboxStatus(ctx context.Context, req *runtimeapi.PodSandboxStatusRequest) (*runtimeapi.PodSandboxStatusResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(PodSandboxStatus, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.PodSandboxStatus(ctx, req.(*runtimeapi.PodSandboxStatusRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.PodSandboxStatusResponse), err
}

func (ci *CriInterceptor) ListPodSandbox(ctx context.Context, req *runtimeapi.ListPodSandboxRequest) (*runtimeapi.ListPodSandboxResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ListPodSandbox, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ListPodSandbox(ctx, req.(*runtimeapi.ListPodSandboxRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ListPodSandboxResponse), err
}

func (ci *CriInterceptor) CreateContainer(ctx context.Context, req *runtimeapi.CreateContainerRequest) (*runtimeapi.CreateContainerResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(CreateContainer, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.CreateContainer(ctx, req.(*runtimeapi.CreateContainerRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.CreateContainerResponse), err
}

func (ci *CriInterceptor) StartContainer(ctx context.Context, req *runtimeapi.StartContainerRequest) (*runtimeapi.StartContainerResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(StartContainer, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.StartContainer(ctx, req.(*runtimeapi.StartContainerRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.StartContainerResponse), err
}

func (ci *CriInterceptor) StopContainer(ctx context.Context, req *runtimeapi.StopContainerRequest) (*runtimeapi.StopContainerResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(StopContainer, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.StopContainer(ctx, req.(*runtimeapi.StopContainerRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.StopContainerResponse), err
}

func (ci *CriInterceptor) RemoveContainer(ctx context.Context, req *runtimeapi.RemoveContainerRequest) (*runtimeapi.RemoveContainerResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(RemoveContainer, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.RemoveContainer(ctx, req.(*runtimeapi.RemoveContainerRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.RemoveContainerResponse), err
}

func (ci *CriInterceptor) ContainerStatus(ctx context.Context, req *runtimeapi.ContainerStatusRequest) (*runtimeapi.ContainerStatusResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ContainerStatus, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ContainerStatus(ctx, req.(*runtimeapi.ContainerStatusRequest))
		})
	if err != nil {
		return nil, err
	}

	return rsp.(*runtimeapi.ContainerStatusResponse), err
}

func (ci *CriInterceptor) ListContainers(ctx context.Context, req *runtimeapi.ListContainersRequest) (*runtimeapi.ListContainersResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ListContainers, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ListContainers(ctx, req.(*runtimeapi.ListContainersRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ListContainersResponse), err
}

func (ci *CriInterceptor) UpdateContainerResources(ctx context.Context, req *runtimeapi.UpdateContainerResourcesRequest) (*runtimeapi.UpdateContainerResourcesResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(UpdateContainerResources, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.UpdateContainerResources(ctx, req.(*runtimeapi.UpdateContainerResourcesRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.UpdateContainerResourcesResponse), err
}

func (ci *CriInterceptor) ContainerStats(ctx context.Context, req *runtimeapi.ContainerStatsRequest) (*runtimeapi.ContainerStatsResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ContainerStats, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ContainerStats(ctx, req.(*runtimeapi.ContainerStatsRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ContainerStatsResponse), err
}
func (ci *CriInterceptor) ListContainerStats(ctx context.Context, req *runtimeapi.ListContainerStatsRequest) (*runtimeapi.ListContainerStatsResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ListContainerStats, ctx, req,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ListContainerStats(ctx, req.(*runtimeapi.ListContainerStatsRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ListContainerStatsResponse), err
}
func (ci *CriInterceptor) Status(ctx context.Context, req *runtimeapi.StatusRequest) (*runtimeapi.StatusResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(Status, ctx, req,

		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.Status(ctx, req.(*runtimeapi.StatusRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.StatusResponse), err
}

//  the following 8 functions are new added ones

func (ci *CriInterceptor) ReopenContainerLog(ctx context.Context, in *runtimeapi.ReopenContainerLogRequest) (*runtimeapi.ReopenContainerLogResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ReopenContainerLog, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ReopenContainerLog(ctx, req.(*runtimeapi.ReopenContainerLogRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ReopenContainerLogResponse), err
}
func (ci *CriInterceptor) ExecSync(ctx context.Context, in *runtimeapi.ExecSyncRequest) (*runtimeapi.ExecSyncResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ExecSync, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ExecSync(ctx, req.(*runtimeapi.ExecSyncRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ExecSyncResponse), err
}
func (ci *CriInterceptor) Exec(ctx context.Context, in *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(Exec, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.Exec(ctx, req.(*runtimeapi.ExecRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ExecResponse), err
}

// Attach
// TODO: Attach is long-link?
func (ci *CriInterceptor) Attach(ctx context.Context, in *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(Attach, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.Attach(ctx, req.(*runtimeapi.AttachRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.AttachResponse), err
}
func (ci *CriInterceptor) PortForward(ctx context.Context, in *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(PortForward, ctx, in,

		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.PortForward(ctx, req.(*runtimeapi.PortForwardRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.PortForwardResponse), err
}

/*
func (ci *CriInterceptor) PodSandboxStats(ctx context.Context, in *runtimeapi.PodSandboxStatsRequest) (*runtimeapi.PodSandboxStatsResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(PodSandboxStats, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.PodSandboxStats(ctx, req.(*runtimeapi.PodSandboxStatsRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.PodSandboxStatsResponse), err
}

func (ci *CriInterceptor) ListPodSandboxStats(ctx context.Context, in *runtimeapi.ListPodSandboxStatsRequest) (*runtimeapi.ListPodSandboxStatsResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(ListPodSandbox, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.ListPodSandboxStats(ctx, req.(*runtimeapi.ListPodSandboxStatsRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.ListPodSandboxStatsResponse), err
}
*/
func (ci *CriInterceptor) UpdateRuntimeConfig(ctx context.Context, in *runtimeapi.UpdateRuntimeConfigRequest) (*runtimeapi.UpdateRuntimeConfigResponse, error) {
	rsp, err := ci.interceptRuntimeRequest(UpdateRuntimeConfig, ctx, in,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return ci.runtimeClient.UpdateRuntimeConfig(ctx, req.(*runtimeapi.UpdateRuntimeConfigRequest))
		})
	if err != nil {
		return nil, err
	}
	return rsp.(*runtimeapi.UpdateRuntimeConfigResponse), err
}
