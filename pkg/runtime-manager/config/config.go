package config

type FailurePolicyType string

const (
	// PolicyFail returns error to caller when got an error cri hook server
	PolicyFail FailurePolicyType = "Fail"
	// PolicyIgnore transfer cri request to containerd/dockerd when got an error to cri serer
	PolicyIgnore FailurePolicyType = "Ignore"
)

type RuntimeHookType string

const (
	defaultRuntimeHookConfigPath string = "/etc/runtime/runtimehook.d"
)

const (
	PreRunPodSandbox            RuntimeHookType = "PreRunPodSandbox"
	PreStartContainer           RuntimeHookType = "PreStartContainer"
	PostStartContainer          RuntimeHookType = "PostStartContainer"
	PreUpdateContainerResources RuntimeHookType = "PreUpdateContainerResources"
	PostStopContainer           RuntimeHookType = "PostStopContainer"
	NoneRuntimeHookType         RuntimeHookType = "NoneRuntimeHookType"
)

type RuntimeHookConfig struct {
	RemoteEndpoint string            `json:"remote-endpoint,omitempty"`
	FailurePolicy  FailurePolicyType `json:"failure-policy,omitempty"`
	RuntimeHooks   []RuntimeHookType `json:"runtime-hooks,omitempty"`
}

type RuntimeHookConfigs struct {
	configs []RuntimeHookConfig
}
