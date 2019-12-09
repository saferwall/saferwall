LINUX_KERNEL_VERSION = 5.4
LINUX_KERNEL_ARCHIVE = linux-$(LINUX_KERNEL_VERSION).tar.gz
LINUX_KERNEL_DOWNLOAD_URL = https://github.com/torvalds/linux/archive/v$(LINUX_KERNEL_VERSION).tar.gz
LINUX_DIR = sandbox/src/linux

kvm-download:		## Download linux kernel which has KVM inside.
	mkdir -p $(LINUX_DIR)
	wget $(LINUX_KERNEL_DOWNLOAD_URL) -O $(LINUX_DIR)/$(LINUX_KERNEL_ARCHIVE)
	tar zxvf $(LINUX_DIR)/$(LINUX_KERNEL_ARCHIVE) -C $(LINUX_DIR) --strip 1
	rm $(LINUX_DIR)/$(LINUX_KERNEL_ARCHIVE)

linux-install-kernel:	## Install last (long) stable linux kernel.
	# Though you can download the last version from kvm git tree,
	# which might not be merged yet into the linux master tree,
	# Here we are just using the one which ships along with the
	# kernel. At the end, we will keep building only the kvm kernel
	# module and load it every time we have a new release.
	sudo apt-get install build-essential libncurses-dev bison flex libssl-dev libelf-dev
	cat /proc/version
	cd $(LINUX_DIR) \
		&& cp -v /boot/config-$(uname -r) .config \
		# Compiles the main kernel
		&& make -j $(nproc) \
		# Install linux kernel modules
		make -j $(nproc) modules_install \
		# Install the linux kernel
		sudo make -j $(nproc) install
		

kvm-mod-install:	## Compile and install kvm
	make M=arch/x86/kvm
	make M=arch/x86/kvm modules_install
	sudo rmmod kvm-intel kvm
	# just loading kvm-intel will load kvm automatically
	sudo modprobe kvm-intel

kvm-tools-install: 	## Install KVM related tools
	sudo apt install virt-manager qemu-kvm bridge-utils -y
	sudo kvm-ok
	kvm --version

kvm-libvirt:
	sudo virsh net-list --all
	sudo service libvirtd status

kvm-qemu:		## Clone qemu
	cd vmi \
		&& git clone git://git.qemu-project.org/qemu.git
