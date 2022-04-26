package resource_executor

import (
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func transferResource(r *runtimeapi.LinuxContainerResources) *v1alpha1.LinuxContainerResources {
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
