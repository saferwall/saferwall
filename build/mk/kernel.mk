CURRENT_LINUX_VERSION = v5.1

kernel-clone:		## git clone the linux kernel
	cd vmi \
		&& git clone https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git

kernel-prepare:		## Install required tools to build the kernel
	sudo apt install bison flex libelf-dev

kernel-checkout:	## Checkout the kernel version supported
	cd vmi/linux \
		&& git fetch --all --tags --prune \
		&& git checkout tags/$(CURRENT_LINUX_VERSION)

