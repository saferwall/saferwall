// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/saferwall/saferwall/internal/utils"
)

// VM represents a virtual machine config.
type VM struct {
	// ID identify uniquely the VM.
	ID int32
	// Name of the VM, should match: Windows-10-x64-1 or Windows-7-x86-2.
	Name string
	// IP address of the VM.
	IP string
	// Snapshots list names.
	Snapshots []string
	// InUse represents the availability of the VM.
	InUse bool
	// Indicates if the VM is healthy.
	IsHealthy bool
	// Pointer to the domain object.
	Dom *libvirt.Domain
}

// findFreeVM iterates over the list of available VM and find
// one which is currently not in use.
func (s *Service) findFreeVM(preferredOS string) *VM {
	var freeVM *VM
	s.mu.Lock()
	for _, vm := range s.vms {
		// Todo: use `preferredOS` when looking for free VMs.
		if !vm.InUse && vm.IsHealthy {
			vm.InUse = true
			freeVM = &vm
			break
		}
	}
	s.mu.Unlock()
	return freeVM
}

// freeVM makes the VM free for consumption.
func (s *Service) freeVM(vm *VM) {
	vm.InUse = false
}

// markStale marks the VM as non-healthy.
func (s *Service) markStale(vm *VM) {
	vm.IsHealthy = false
}

// OS() parses the name of the VM and return a pretty name.
func (vm *VM) OS() string {
	r := `(?P<os>W\w+)-(?P<version>\d{1,2})-(?P<platform>\x86|x64)-(?P<number>\d{1,2})`
	m := utils.RegSubMatchToMapString(r, vm.Name)
	if m["os"] != "" || m["version"] != "" || m["platform"] != "" {
		return "windows-7-x64"
	} else {
		return m["os"] + "-" + m["version"] + "-" + m["platform"]
	}
}
