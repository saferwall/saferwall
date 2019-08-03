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
	minikube start --cpus 4 --memory 8192 --disk-size=60GB

k8s-deploy-saferwall:	k8s-deploy-nfs-server k8s-deploy-minio k8s-deploy-cb k8s-deploy-nsq k8s-deploy-backend k8s-deploy-consumer k8s-deploy-multiav ## Deploy all stack in k8s

k8s-deploy-nfs-server:	## Deploy NFS server in a newly created k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl create -f nfs-server.yaml \
	&& kubectl create -f samples-pv.yaml \
	&& kubectl create -f samples-pvc.yaml

k8s-deploy-cb:	## Deploy couchbase in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f couchbase-sc.yaml ; \
	kubectl create -f couchbase-pv.yaml ; \
	kubectl create -f couchbase-pvc.yaml ; \
	kubectl create -f crd.yaml ; \
	kubectl create -f operator-role.yaml ; \
	kubectl create serviceaccount couchbase-operator --namespace default ; \
	kubectl create rolebinding couchbase-operator --role couchbase-operator --serviceaccount default:couchbase-operator ; \
	kubectl create -f admission.yaml ; \
	kubectl create -f secret.yaml ; \
	kubectl create -f operator-deployment.yaml ; \
	kubectl create -f couchbase-cluster.yaml  

k8s-deploy-nsq:			## Deploy NSQ in a newly created k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl create -f nsqlookupd.yaml \
	&& kubectl create -f nsqd.yaml \
	&& kubectl create -f nsqadmin.yaml
	
k8s-deploy-minio:		## Deploy minio
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl create -f minio-standalone-pvc.yaml ; \
	kubectl create -f minio-standalone-deployment.yaml ; \
	kubectl create -f minio-standalone-service.yaml

k8s-deploy-fsecure:		## Deploy fsecure
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl delete deployments fsecure ; \
	kubectl apply -f multiav-fsecure.yaml

k8s-deploy-multiav:		## Deploy multiav in a newly created k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl apply -f multiav-clamav.yaml \
	&& kubectl apply -f multiav-avira.yaml \
	&& kubectl apply -f multiav-eset.yaml \
	&& kubectl apply -f multiav-kaspersky.yaml \
	&& kubectl apply -f multiav-comodo.yaml \
	&& kubectl apply -f multiav-fsecure.yaml \
	&& kubectl apply -f multiav-bitdefender.yaml \
	&& kubectl apply -f multiav-avast.yaml \
	&& kubectl apply -f multiav-windefender.yaml \
	&& kubectl apply -f multiav-symantec.yaml

k8s-apply-multiav:		## Delete multiav from k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl apply -f multiav-clamav.yaml \
	&& kubectl apply -f multiav-avira.yaml \
	&& kubectl apply -f multiav-eset.yaml \
	&& kubectl apply -f multiav-kaspersky.yaml \
	&& kubectl apply -f multiav-comodo.yaml \
	&& kubectl apply -f multiav-fsecure.yaml \
	&& kubectl apply -f multiav-bitdefender.yaml \
	&& kubectl apply -f multiav-avast.yaml \
	&& kubectl apply -f multiav-windefender.yaml \
	&& kubectl apply -f multiav-symantec.yaml


k8s-deploy-backend:		## Deploy backend in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl delete deployments backend ;\
	kubectl apply -f backend.yaml

k8s-deploy-consumer:		## Deploy consumer in kubernetes cluster
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl apply -f consumer.yaml

k8s-delete-nsq:
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl delete svc nsqd nsqadmin nsqlookupd
	kubectl delete deployments nsqadmin 
	kubectl delete deployments nsqadmin

k8s-delete-cb:		## Delete all couchbase related objects from k8s
	kubectl delete cbc cb-saferwall
	kubectl delete deployment couchbase-operator
	kubectl delete crd couchbaseclusters.couchbase.com
	kubectl delete secret cb-saferwall-auth

k8s-delete:			## delete all
	kubectl delete deployments,service backend -l app=web
	kubectl delete service backend
	kubectl delete service consumer
	kubectl delete deployments consumer ; kubectl apply -f consumer.yaml
	kubectl delete deployments avast ; kubectl apply -f multiav-avast.yaml
	kubectl delete deployments avira ; kubectl apply -f multiav-avira.yaml
	kubectl delete deployments bitdefender ; kubectl apply -f multiav-bitdefender.yaml
	kubectl delete deployments comodo ; kubectl apply -f multiav-comodo.yaml
	kubectl delete deployments eset ; kubectl apply -f multiav-eset.yaml
	kubectl delete deployments symantec ; kubectl apply -f multiav-symantec.yaml
	kubectl delete deployments kaspersky ; kubectl apply -f multiav-kaspersky.yaml
	kubectl delete deployments windefender ; kubectl apply -f multiav-windefender.yaml
	kubectl delete cbc cb-saferwall ; kubectl create -f couchbase-cluster.yaml
	kubectl delete deployments backend ; kubectl apply -f backend.yaml
