k8s-kubectl-install:	## Install kubectl
	sudo apt-get update && sudo apt-get install -y apt-transport-https
	curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
	echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
	sudo apt-get update
	sudo apt-get install -y kubectl

k8s-minikube-install:	## Install minikube
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	chmod +x minikube
	sudo cp minikube /usr/local/bin && rm minikube

k8s-minikube-start:		## Start minikube
	minikube start --cpus 4 --memory 8192

k8s-deploy-cb:	## Deploy couchbase in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f crd.yaml ; \
	kubectl create -f operator-role.yaml ; \
	kubectl create serviceaccount couchbase-operator --namespace default ; \
	kubectl create rolebinding couchbase-operator --role couchbase-operator --serviceaccount default:couchbase-operator ; \
	kubectl create -f secret.yaml ; \
	kubectl create -f operator-deployment.yaml ; \
	kubectl create -f couchbase-cluster.yaml  

k8s-delte-cb:
	kubectl delete cbc cb-saferwall
	kubectl delete deployment couchbase-operator
	kubectl delete crd couchbaseclusters.couchbase.com
	kubectl delete secret cb-saferwall-auth

k8s-deploy-nsq:			## Deploy NSQ in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f nsqlookupd.yaml ; \
	kubectl create -f nsqd.yaml ; \
	kubectl create -f nsqadmin.yaml
	
k8s-deploy-minio:		## Deploy minio
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f minio-standalone-pvc.yaml ; \
	kubectl create -f minio-standalone-deployment.yaml ; \
	kubectl create -f minio-standalone-service.yaml

k8s-deploy-multiav:		## Deploy multiav
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f multiav-clamav.yaml

k8s-deploy-backend:		## Deploy backend in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl apply -f backend.yaml

k8s-deploy-consumer:		## Deploy consumer in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl apply -f consumer.yaml

k8s-delete-nsq:
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl delete svc nsqd nsqadmin nsqlookupd
	kubectl delete deployments nsqadmin 
	kubectl delete deployments nsqadmin 

k8s-delete:
	kubectl delete deployments,service backend -l app=web
	kubectl delete service backend
	kubectl delete service consumer

	kubectl delete deployments consumer
	kubectl delete deployments avast
	kubectl delete deployments backend

	kubectl apply -f multiav-avast.yaml
	kubectl apply -f backend.yaml
	kubectl apply -f consumer.yaml