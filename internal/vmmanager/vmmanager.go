// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Terraform provider for libvirt contains nice usage of the go-libvirt library.
// https://github.com/dmacvicar/terraform-provider-libvirt.git

package vmmanager

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"golang.org/x/crypto/ssh"
)


const (
	// Timeout used to connect to the libvirt server.
	dialTimeout = 20 * time.Second

	// Flags field currently unused in libvirt.
	flagsUnused = 0

	// Maximum snapshot to get.
	maxSnapshotLen = 3

	defaultUnixSock = "/var/run/libvirt/libvirt-sock"
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

// create a net SSH connection.
func dialSSH(hostname, username, port, sshKeyPath string) (net.Conn, error) {

	sshKey, err := ioutil.ReadFile(os.ExpandEnv(sshKeyPath))
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(sshKey)
	if err != nil {
		return nil, err
	}
	authMethods := make([]ssh.AuthMethod, 0)
	authMethods = append(authMethods, ssh.PublicKeys(signer))
	hostKeyCallback := ssh.InsecureIgnoreHostKey()

	cfg := ssh.ClientConfig{
		User:            username,
		HostKeyCallback: hostKeyCallback,
		Auth:            authMethods,
		Timeout:         dialTimeout,
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), &cfg)
	if err != nil {
		return nil, err
	}

	address := defaultUnixSock
	c, err := sshClient.Dial("unix", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt on the remote host: %w", err)
	}

	return c, nil
}

// New creates a new libvirt RPC connection.  It dials libvirt
// either on the local machine or the remote one depending on
// the transport parameter "unix" for local and "ssh" for remote connections.
func New(transport, address, port, user, sshKeyPath string) (VMManager, error) {

	var err error
	var conn *libvirt.Libvirt

	switch transport {
	case "unix":
		dialer := dialers.NewLocal(dialers.WithLocalTimeout(dialTimeout))
		conn = libvirt.NewWithDialer(dialer)
	case "ssh":
		c, err := dialSSH(address, user, port, sshKeyPath)
		if err != nil {
			return VMManager{}, err
		}
		dialer := dialers.NewAlreadyConnected(c)
		conn = libvirt.NewWithDialer(dialer)
	}

	err = conn.Connect()
	if err != nil {
		return VMManager{}, err
	}

	return VMManager{Conn: conn}, nil
}

// Domains retrieves the list of domains.
func (vmm *VMManager) Domains() ([]Domain, error) {
	flags := libvirt.ConnectListDomainsActive
	dd, _, err := vmm.Conn.ConnectListAllDomains(1, flags)
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

		flags := libvirt.DomainSnapshotListActive
		names, err := vmm.Conn.DomainSnapshotListNames(d, maxSnapshotLen, uint32(flags))
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
