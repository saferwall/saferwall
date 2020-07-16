DOWNLOAD_URL = https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

minikube-install:		## Install minikube
	curl -Lo minikube $(DOWNLOAD_URL)
	chmod +x minikube
	sudo cp minikube /usr/local/bin && rm minikube
	minikube version

minikube-start:			## Start minikube
ifeq ($(MINIKUBE_DRIVER),none)
	sudo apt update 
	sudo apt install -y conntrack
	sudo minikube start --driver=none
else
	minikube start --driver=$(MINIKUBE_DRIVER)  --cpus $(MINIKUBE_CPU) --memory $(MINIKUBE_MEMORY) --disk-size=$(MINIKUBE_DISK_SIZE)GB
endif
