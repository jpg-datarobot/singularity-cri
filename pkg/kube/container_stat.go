// Copyright (c) 2018 Sylabs, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kube

import (
	"fmt"

	"github.com/containerd/cgroups"
	"github.com/sylabs/cri/pkg/fs"
)

// ContainerStat holds information about container resources usage.
type ContainerStat struct {
	// Writable layer fs usage.
	Fs *fs.UsageInfo
	// Total memory used by container in bytes
	Memory uint64
	// Total CPU used in nanoseconds.
	CPU uint64
}

// Stat fetches information about container resources usage. This method
// implies that cpuacct and memory cgroups controllers are mounted on host
// at /sys/fs/cgroups/cpuacct and  /sys/fs/cgroups/memory respectively.
func (c *Container) Stat() (*ContainerStat, error) {
	fsInfo, err := fs.Usage(c.baseDir())
	if err != nil {
		return nil, fmt.Errorf("could not get fs usage: %v", err)
	}
	cgroup, err := cgroups.Load(cgroups.V1, cgroups.PidPath(c.Pid()))
	if err != nil {
		return nil, fmt.Errorf("could not load cgroups: %v", err)
	}

	metrics, err := cgroup.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		return nil, fmt.Errorf("could not fetch metrics: %v", err)
	}

	var cpuTotal uint64
	var memoryTotal uint64
	if metrics.CPU != nil && metrics.CPU.Usage != nil {
		cpuTotal = metrics.CPU.Usage.Total
	}
	if metrics.Memory != nil && metrics.Memory.Usage != nil {
		memoryTotal = metrics.Memory.Usage.Usage
	}

	return &ContainerStat{
		Fs:     fsInfo,
		Memory: memoryTotal,
		CPU:    cpuTotal,
	}, nil
}
