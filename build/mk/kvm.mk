kvm-install:	## Compile and install kvm
	make M=arch/x86/kvm
	make M=arch/x86/kvm modules_install