k8s-kubectl-install:	## Install Kubectl
	sudo apt-get update && sudo apt-get install -y apt-transport-https
	curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
	echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
	sudo apt-get update
	sudo apt-get install -y kubectl

k8s-minikube-install:	## Install Minikube
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x minikube
	sudo cp minikube /usr/local/bin && rm minikube

k8s-deploy-minio:		## Deploy minio
	kubectl create -f https://github.com/minio/minio/blob/master/docs/orchestration/kubernetes/minio-standalone-pvc.yaml?raw=true
	kubectl create -f https://github.com/minio/minio/blob/master/docs/orchestration/kubernetes/minio-standalone-deployment.yaml?raw=true
	kubectl create -f https://github.com/minio/minio/blob/master/docs/orchestration/kubernetes/minio-standalone-service.yaml?raw=true

k8s-deploy-couchbase:	## Deploy couchbase
	kubectl create -f crd.yaml
	kubectl create -f cluster-role-sa.yaml
	kubectl create -f cluster-role-user.yaml
	kubectl create serviceaccount couchbase-operator --namespace default
	kubectl create clusterrolebinding couchbase-operator --clusterrole couchbase-operator --serviceaccount default:couchbase-operator
	kubectl create -f operator.yaml
	kubectl create -f secret.yaml
	cbopctl create -f couchbase-cluster.yaml
	kubectl create -f backend.yaml
	
	kubectl delete deployment couchbase-operator
	kubectl delete cbc cb-saferwall
	kubectl delete deployments backend