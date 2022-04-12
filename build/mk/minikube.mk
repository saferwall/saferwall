MK_VERSION = v1.25.2
MK_DOWNLOAD_URL = https://github.com/kubernetes/minikube/releases/download/$(MK_VERSION)/minikube-linux-amd64

minikube-install:		## Install minikube
	@minikube version | grep $(MK_VERSION); \
		if [ $$? -eq 1 ]; then \
			curl -Lo minikube $(MK_DOWNLOAD_URL); \
			chmod +x minikube; \
			sudo cp minikube /usr/local/bin && rm minikube; \
			minikube version; \
		else \
            echo "${GREEN} [*] Minikube already installed ${RESET}"; \
		fi

minikube-up:			## Start minikube cluster.
	sudo apt update -qq
	sudo apt install -qq nfs-common
ifeq ($(MINIKUBE_DRIVER),none)
	sudo apt install -qq conntrack
	CHANGE_MINIKUBE_NONE_USER=true sudo -E minikube start --driver=none
else
	minikube start --driver=$(MINIKUBE_DRIVER) --cpus $(MINIKUBE_CPU) \
		--memory $(MINIKUBE_MEMORY) --disk-size=$(MINIKUBE_DISK_SIZE)
endif
	kubectl version
	kubectl cluster-info
	kubectl config get-contexts
	kubectl config current-context
	kubectl config use-context minikube
	minikube status
ifneq ($(MINIKUBE_DRIVER),none)
	minikube addons enable ingress
endif

minikube-down:			## Stop and delete minikube cluster.
ifeq ($(MINIKUBE_DRIVER),none)
	sudo -E minikube stop
	sudo -E minikube delete
else
	minikube stop
	minikube delete
endif
