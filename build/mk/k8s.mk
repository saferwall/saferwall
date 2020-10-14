kubectl-install:		## Install kubectl
	sudo apt-get update && sudo apt-get install -y apt-transport-https
	curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
	echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
	sudo apt-get update
	sudo apt-get install -y kubectl
	kubectl version --client

KUBECTX_URL=https://github.com/ahmetb/kubectx/releases/download/v0.9.0/kubectx_v0.9.0_linux_x86_64.tar.gz
KUBENS_URL = https://github.com/ahmetb/kubectx/releases/download/v0.9.0/kubens_v0.9.0_linux_x86_64.tar.gz
k8s-kubectx-install:		## Install kubectx/kubens
	wget -N $(KUBECTX_URL) -O /tmp/kubectx.tar.gz
	tar zxvf /tmp/kubectx.tar.gz -C /tmp
	sudo mv /tmp/kubectx /usr/local/bin/
	chmod +x /usr/local/bin/kubectx
	wget -N $(KUBENS_URL) -O /tmp/kubens.tar.gz
	tar zxvf /tmp/kubens.tar.gz -C /tmp
	sudo mv /tmp/kubens /usr/local/bin/
	chmod +x /usr/local/bin/kubens

k8s-prepare:	k8s-kubectl-install k8s-kube-capacity k8s-minikube-start ## Install minikube, kubectl, kube-capacity and start a cluster

k8s-deploy-saferwall:	k8s-deploy-nfs-server k8s-deploy-minio k8s-deploy-cb k8s-deploy-nsq k8s-deploy-backend k8s-deploy-consumer k8s-deploy-multiav ## Deploy all stack in k8s

k8s-deploy-nfs-server:	## Deploy NFS server in a newly created k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl apply -f nfs-server.yaml \
	&& kubectl apply -f samples-pv.yaml \
	&& kubectl apply -f samples-pvc.yaml

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
	kubectl apply -f couchbase-cluster.yaml  

k8s-deploy-nsq:			## Deploy NSQ in a newly created k8s cluster
	cd  $(ROOT_DIR)/build/k8s \
	&& kubectl apply -f nsqlookupd.yaml \
	&& kubectl apply -f nsqd.yaml \
	&& kubectl apply -f nsqadmin.yaml
	
k8s-deploy-minio:		## Deploy minio
	cd  $(ROOT_DIR)/build/k8s ; \
	kubectl apply -f minio-standalone-pvc.yaml ; \
	kubectl apply -f minio-standalone-deployment.yaml ; \
	kubectl apply -f minio-standalone-service.yaml

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
	&& kubectl apply -f multiav-symantec.yaml \
	&& kubectl apply -f multiav-sophos.yaml \
	&& kubectl apply -f multiav-mcafee.yaml \
	&& kubectl apply -f seccomp-profile.yaml \
	&& kubectl apply -f seccomp-installer.yaml \
	&& kubectl apply -f multiav-windefender.yaml

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
	kubectl delete cbc cb-saferwall ; \
	kubectl delete deployment couchbase-operator-admission ; \
	kubectl delete deployment couchbase-operator  ; \
	kubectl delete crd couchbaseclusters.couchbase.com  ; \
	kubectl delete secret cb-saferwall-auth ; \
	kubectl delete pvc couchbase-pvc ; \
	kubectl delete pv couchbase-pv ; \
	kubectl delete sc couchbase-sc

k8s-delete-multiav:		## Delete all multiav related objects from k8s
	cd  $(ROOT_DIR)/build/k8s ; \
		kubectl delete deployments avast ; kubectl apply -f multiav-avast.yaml ; \
		kubectl delete deployments avira ; kubectl apply -f multiav-avira.yaml ; \
		kubectl delete deployments bitdefender ; kubectl apply -f multiav-bitdefender.yaml ; \
		kubectl delete deployments comodo ; kubectl apply -f multiav-comodo.yaml ; \
		kubectl delete deployments eset ; kubectl apply -f multiav-eset.yaml ; \
		kubectl delete deployments fsecure ; kubectl apply -f multiav-fsecure.yaml ; \
		kubectl delete deployments symantec ; kubectl apply -f multiav-symantec.yaml ; \
		kubectl delete deployments kaspersky ; kubectl apply -f multiav-kaspersky.yaml ; \
		kubectl delete deployments windefender ; kubectl apply -f multiav-windefender.yaml

k8s-delete:			## delete all
	kubectl delete deployments,service backend -l app=web
	kubectl delete service backend
	kubectl delete service consumer
	kubectl delete deployments consumer ; kubectl apply -f consumer.yaml

	kubectl delete cbc cb-saferwall ; kubectl create -f couchbase-cluster.yaml
	kubectl delete deployments backend ; kubectl apply -f backend.yaml

k8s-pf:				## Port fordward
	kubectl port-forward --namespace default $(POD) $(PORT):$(PORT)

k8s-kube-capacity: 	## Install kube-capacity
	wget https://github.com/robscott/kube-capacity/releases/download/0.4.0/kube-capacity_0.4.0_Linux_x86_64.tar.gz -P /tmp
	cd /tmp \
		&& tar zxvf kube-capacity_0.4.0_Linux_x86_64.tar.gz \
		&& sudo mv kube-capacity /usr/local/bin \
		&& kube-capacity

k8s-delete-all-objects: ## Delete all objects
	kubectl delete "$(kubectl api-resources --namespaced=true --verbs=delete -o name | tr "\n" "," | sed -e 's/,$//')" --all

k8s-dump-tls-secrets: ## Dump TLS secrets
	sudo apt install jq -y
	$(eval HELM_RELEASE_NAME := $(shell sudo helm ls --filter saferwall --output json | jq '.[0].name' | tr -d '"'))
	$(eval HELM_SECRET_TLS_NAME := $(HELM_RELEASE_NAME)-tls)
	kubectl get secret $(HELM_SECRET_TLS_NAME) -o jsonpath="{.data['ca\.crt']}" | base64 --decode  >> ca.crt
	kubectl get secret $(HELM_SECRET_TLS_NAME) -o jsonpath="{.data['tls\.crt']}" | base64 --decode  >> tls.crt
	kubectl get secret $(HELM_SECRET_TLS_NAME) -o jsonpath="{.data['tls\.key']}" | base64 --decode  >> tls.key
	openssl pkcs12 -export -out saferwall.p12 -inkey tls.key -in tls.crt -certfile ca.crt

k8s-init-cert-manager: ## Init cert-manager
	# Create the namespace for cert-manager.
	kubectl create namespace cert-manager
	# Install the CustomResourceDefinition
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.2/cert-manager.yaml

k8s-cert-manager-rm-crd: ## Delete cert-manager crd objects.
	kubectl get crd | grep cert-manager | xargs --no-run-if-empty kubectl delete crd
	kubectl delete namespace cert-manager