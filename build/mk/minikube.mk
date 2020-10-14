MK_VERSION = v1.14.0
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
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube start --driver=none
else
	minikube start --driver=$(MINIKUBE_DRIVER)  --cpus $(MINIKUBE_CPU) --memory $(MINIKUBE_MEMORY) --disk-size=$(MINIKUBE_DISK_SIZE)
endif
	kubectl version
	kubectl cluster-info
	kubectl config get-contexts  
	kubectl config current-context
	kubectl config use-context minikube
	minikube status
ifneq ($(MINIKUBE_DRIVER),none)
	sudo minikube addons enable ingress
endif
	

minikube-down:			## Stop and delete minikube cluster.
	sudo -E minikube stop
	sudo -E minikube delete
	sudo rm -rf /tmp/*