package interceptor

type ServiceType int
type RuntimeServiceType int
type ImageServiceType int

const (
	RuntimeService ServiceType = iota
	ImageService
)

const (
	Version RuntimeServiceType = iota
	RunPodSandbox
	StopPodSandbox
	RemovePodSandbox
	PodSandboxStatus
	ListPodSandbox
	CreateContainer
	StartContainer
	StopContainer
	RemoveContainer
	ContainerStatus
	ListContainers
	UpdateContainerResources
	ContainerStats
	ListContainerStats
	Status
	ReopenContainerLog
	ExecSync
	Exec
	Attach
	PortForward
	PodSandboxStats
	ListPodSandboxStats
	UpdateRuntimeConfig
)
