// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/digitalocean/go-libvirt"
	so "github.com/iamacarpet/go-win64api/shared"
	agent "github.com/saferwall/saferwall/internal/agent"
)

var (
	mu sync.Mutex
)

// VM represents a virtual machine config.
type VM struct {
	// ID identify uniquely the VM.
	ID int32
	// Name of the VM.
	Name string
	// IP address of the VM.
	IP string
	// Operating system used by the guest.
	OS string
	// Server version.
	AgentVersion string
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
func findFreeVM(vms []VM, preferredOS string) *VM {
	var freeVM *VM
	mu.Lock()
	for _, vm := range vms {
		if !vm.InUse && vm.IsHealthy && preferredOS == vm.OS {
			vm.InUse = true
			freeVM = &vm
			break
		}
	}
	mu.Unlock()
	return freeVM
}

func (vm *VM) ping() error {

	// Establish a gRPC connection to the agent server running
	// inside the guest.
	client, err := agent.New(vm.IP + defaultGRPCPort)
	if err != nil {
		return err
	}

	ctx := context.Background()
	pingResult, err := client.Ping(ctx)
	if err != nil {
		return err
	}

	var sysInfo so.OperatingSystem
	err = json.Unmarshal(pingResult.SysInfo, &sysInfo)
	if err != nil {
		return err
	}

	os := strings.ReplaceAll(sysInfo.FriendlyName, "Microsoft", "")
	os = strings.ReplaceAll(os, "Professional", "")
	os = strings.TrimSpace(os)
	vm.OS = os + " " + sysInfo.Architecture
	vm.AgentVersion = pingResult.ServerVersion

	return nil
}

// free makes the VM free for consumption.
func (vm *VM) free() {
	vm.InUse = false
}

// markStale marks the VM as non-healthy.
func (vm *VM) markStale() {
	vm.IsHealthy = false
}
