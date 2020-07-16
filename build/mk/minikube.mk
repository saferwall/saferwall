DOWNLOAD_URL = https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

minikube-install:		## Install minikube
	curl -Lo minikube $(DOWNLOAD_URL)
	chmod +x minikube
	sudo cp minikube /usr/local/bin && rm minikube
	minikube version

minikube-start:			## Start minikube
	minikube start --driver=$(MK_DRIVER)  --cpus $(MK_CPU) --memory $(MK_MEM) --disk-size=$(MK_DISK)GB
	kubectl proxy --address='0.0.0.0' --disable-filter=true &
