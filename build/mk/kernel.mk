CURRENT_LINUX_VERSION = v5.1

kernel-clone:		## git clone the linux kernel
	cd vmi \
		&& git clone https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git

kernel-prepare:		## Install required tools to build the kernel
	sudo apt install build-essential libncurses-dev bison flex libssl-dev libelf-dev

kernel-checkout:	## Checkout the kernel version supported
	cd vmi/linux \
		&& git fetch --all --tags --prune \
		&& git checkout tags/$(CURRENT_LINUX_VERSION)

kernel-compile:		## Compile the linux kernel
	make -j $(nproc)

kernel-install:		## Install the linux kernel
	sudo make modules_install
	sudo make install
	sudo update-initramfs -c -k 5.1.0
	sudo update-grub
