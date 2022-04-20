package config

import "strings"

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

type RuntimeRequestPath string

const (
	RunPodSandbox            RuntimeRequestPath = "RunPodSandbox"
	StartContainer           RuntimeRequestPath = "StartContainer"
	UpdateContainerResources RuntimeRequestPath = "UpdateContainerResources"
	StopContainer            RuntimeRequestPath = "StopContainer"
	NoneRuntimeHookPath      RuntimeRequestPath = "NoneRuntimeHookPath"
)

func (ht RuntimeHookType) OccursOn(path RuntimeRequestPath) bool {
	switch ht {
	case PreRunPodSandbox:
		if path == RunPodSandbox {
			return true
		}
	case PreStartContainer:
		if path == StartContainer {
			return true
		}
	case PostStartContainer:
		if path == StartContainer {
			return true
		}
	case PreUpdateContainerResources:
		if path == UpdateContainerResources {
			return true
		}
	case PostStopContainer:
		if path == StopContainer {
			return true
		}
	}
	return false
}

func (hp RuntimeRequestPath) PreHookType() RuntimeHookType {
	if hp == RunPodSandbox {
		return PreRunPodSandbox
	}
	return NoneRuntimeHookType
}

func (hp RuntimeRequestPath) PostHookType() RuntimeHookType {
	if hp == RunPodSandbox {
		return NoneRuntimeHookType
	}
	return NoneRuntimeHookType
}

type RuntimeHookStage string

const (
	PreHook     RuntimeHookStage = "PreHook"
	PostHook    RuntimeHookStage = "PostHook"
	UnknownHook RuntimeHookStage = "UnknownHook"
)

func (ht RuntimeHookType) HookStage() RuntimeHookStage {
	if strings.HasPrefix(string(ht), "Pre") {
		return PreHook
	} else if strings.HasPrefix(string(ht), "Post") {
		return PostHook
	}
	return UnknownHook
}
