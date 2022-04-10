// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/digitalocean/go-libvirt"
	"github.com/saferwall/saferwall/internal/vmmanager"
)

func main() {

	var flagNetwork = flag.String("network", "unix", "network to use")
	var flagAddress = flag.String("address", "192.168.20.24", "IP address of the target server")
	var flagPort = flag.String("port", "22", "ssh port number")
	var flagUser = flag.String("user", "linux", "username for the ssh session")

	flag.Parse()

	s, err := vmmanager.New(*flagNetwork, *flagAddress, *flagPort, *flagUser)
	if err != nil {
		log.Fatalf("failed to create new libvirt conn: %v", err)
	}

	out, _ := s.Domains()

	log.Print(out)

	v, err := s.Conn.Version()
	if err != nil {
		log.Fatalf("failed to retrieve libvirt version: %v", err)
	}
	fmt.Println("Version:", v)

	dd, _, err := s.Conn.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive)
	if err != nil {
		log.Fatalf("failed to list all domains: %v", err)
	}

	rNames, err := s.Conn.DomainSnapshotListNames(dd[0], 5, 0)
	if err != nil {
		log.Fatalf("failed to retrieve snapshot names: %v", err)
	}
	fmt.Println(rNames)

	err = s.Revert(dd[0], "clean-state")
	if err != nil {
		log.Fatalf("failed to revert snapshot: %v", err)
	}

	osType, err := s.Conn.DomainGetOsType(dd[0])
	if err != nil {
		log.Fatalf("failed to retrieve os type: %v", err)
	}
	fmt.Println(osType)

	rParams, err := s.Conn.DomainGetGuestInfo(dd[0], uint32(libvirt.DomainGuestInfoOs), 0)
	if err != nil {
		log.Fatalf("failed to retrieve guest info: %v", err)
	}
	fmt.Println(rParams)

	faces, ret, err := s.Conn.ConnectListAllInterfaces(1, libvirt.ConnectListInterfacesActive)
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	fmt.Println(faces)
	fmt.Println(ret)

	rIfaces, err := s.Conn.DomainInterfaceAddresses(
		dd[0], uint32(libvirt.DomainInterfaceAddressesSrcLease), 0)
	if err != nil {
		log.Printf("failed to retrieve domains: %v", err)
	}

	fmt.Println(rIfaces)
	fmt.Println(ret)

	networks, ret, err := s.Conn.ConnectListAllNetworks(1, libvirt.ConnectListNetworksActive)
	if err != nil {
		log.Fatalf("failed to retrieve betworks: %v", err)
	}

	fmt.Println(networks)
	fmt.Println(ret)

	fmt.Println("ID\tName\t\tUUID")
	fmt.Printf("--------------------------------------------------------\n")
	for _, d := range dd {
		fmt.Printf("%d\t%s\t%x\n", d.ID, d.Name, d.UUID)
	}

	if err := s.Conn.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect: %v", err)
	}
}
