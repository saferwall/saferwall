// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package vmmanager

import (
	"fmt"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
)

const (
	// Timeout used to connect to the libvirt server.
	dialTimeout = 20 * time.Second

	// Flags field currently unused in libvirt.
	flagsUnused = 0

	// Maximum snapshot to get.
	maxSnapshotLen = 3
)

type VMManager struct {
	Conn *libvirt.Libvirt
}

// Domain represents a domain.
type Domain struct {
	// The domain object
	Dom *libvirt.Domain
	// IP address of the VM.
	IP string
	// Snapshot list the VM has.
	Snapshots []string
}

// New creates a new libvirt RPC connection.  It dials libvirt
// either on the local machine or the remote one depending on
// the network parameter ("tcp" for rmote and "unix" for local),
func New(network, address, port, user string) (VMManager, error) {

	var err error
	var conn *libvirt.Libvirt

	switch network {
	case "unix":
		dialer := dialers.NewLocal(dialers.WithLocalTimeout(dialTimeout))
		conn = libvirt.NewWithDialer(dialer)
		err = conn.Connect()
	case "tcp":
		dialer := dialers.NewRemote(address, dialers.UsePort(port))
		conn = libvirt.NewWithDialer(dialer)
		uri := fmt.Sprintf("qemu+ssh://%s@%s/system", user, address)
		err = conn.ConnectToURI(libvirt.ConnectURI(uri))

	}

	if err != nil {
		return VMManager{}, err
	}

	return VMManager{Conn: conn}, nil
}

// Domains retrieves the list of domains.
func (vmm *VMManager) Domains() ([]Domain, error) {
	dd, _, err := vmm.Conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	for _, d := range dd {
		addresses, err := vmm.Conn.DomainInterfaceAddresses(
			d, uint32(libvirt.DomainInterfaceAddressesSrcLease), flagsUnused)
		if err != nil {
			return nil, err
		}

		names, err := vmm.Conn.DomainSnapshotListNames(d, maxSnapshotLen, 0)
		if err != nil {
			return nil, err
		}

		domains = append(domains, Domain{
			Dom:       &d,
			IP:        addresses[0].Addrs[0].Addr,
			Snapshots: names,
		})
	}

	return domains, nil
}

// Revert reverts the domain to a particular snapshot.
func (vmm *VMManager) Revert(dom libvirt.Domain, name string) error {

	snap, err := vmm.Conn.DomainSnapshotLookupByName(dom, name, flagsUnused)
	if err != nil {
		return err
	}
	return vmm.Conn.DomainRevertToSnapshot(snap, uint32(libvirt.DomainSnapshotRevertRunning))
}
