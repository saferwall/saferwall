kvm-mod-install:	## Compile and install kvm
	make M=arch/x86/kvm
	make M=arch/x86/kvm modules_install
	sudo rmmod kvm-intel kvm
	# just loading kvm-intel will load kvm automatically
	sudo modprobe kvm-intel

kvm-tools-install: ## Install KVM related tools
	sudo apt install virt-manager qemu-kvm bridge-utils -y
	sudo kvm-ok

kvm-libvirt:
	sudo virsh net-list --all
	sudo service libvirtd status

kvm-qemu:			## Clone qemu
	cd vmi \
		&& git clone git://git.qemu-project.org/qemu.git