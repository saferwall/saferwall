MK_VERSION = v1.12.2
MK_DOWNLOAD_URL = https://github.com/kubernetes/minikube/releases/download/$(MK_VERSION)/minikube-linux-amd64

minikube-install:		## Install minikube
	curl -Lo minikube $(MK_DOWNLOAD_URL)
	chmod +x minikube
	sudo cp minikube /usr/local/bin && rm minikube
	minikube version

minikube-up:			## Start minikube cluster.
	sudo apt update 
	sudo apt install -y nfs-common
ifeq ($(MINIKUBE_DRIVER),none)
	sudo apt install -y conntrack
	sudo minikube start --driver=none
	# sudo mv /root/.kube $(HOME)/.kube # this will write over any previous configuration
	sudo chown -R $(USER) $(HOME)/.kube
	sudo chgrp -R $(USER) $(HOME)/.kube
	# sudo mv /root/.minikube $(HOME)/.minikube # this will write over any previous configuration
	sudo chown -R $(USER) $(HOME)/.minikube
	sudo chgrp -R $(USER) $(HOME)/.minikube
else
	minikube start --driver=$(MINIKUBE_DRIVER)  --cpus $(MINIKUBE_CPU) --memory $(MINIKUBE_MEMORY) --disk-size=$(MINIKUBE_DISK_SIZE)
endif
	kubectl version
	kubectl cluster-info
	kubectl config get-contexts  
	kubectl config current-context
	kubectl config use-context minikube
	minikube status
	sudo minikube addons enable ingress

minikube-down:			## Stop and delete minikube cluster.
	sudo minikube stop && sudo minikube delete