awscli-install:		## Install aws cli tool
	sudo apt install python-pip -y
	pip install awscli
	aws configure

aws-create-user:	## Create user to provision the cluster
	aws iam create-group --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess --group-name kops
	aws iam create-user --user-name kops
	aws iam add-user-to-group --user-name kops --group-name kops
	aws iam create-access-key --user-name kops

kops-install:		## Install Kubernetes Kops
	curl -Lo kops https://github.com/kubernetes/kops/releases/download/v1.16.0/kops-linux-amd64
	chmod +x ./kops
	sudo mv ./kops /usr/local/bin/
	kops version

ZONE = us-east-1
kops-create-kops-bucket:		## Create s3 bucket for kops
	aws s3api create-bucket --bucket kops-saferwall-com-state-store --region $(ZONE)
	aws s3api put-bucket-versioning --bucket kops-saferwall-com-state-store --versioning-configuration Status=Enabled

NODE_COUNT = 2
NODE_SIZE = t2.xlarge
FS_TOKEN = saferwall-efs
kops-create-cluster:			## Create k8s cluster
	kubectl config get-contexts
	aws ec2 describe-availability-zones --region $(ZONE)
	kops create cluster \
		--zones us-east-1a \
		--node-count $(NODE_COUNT) \
		--node-size $(NODE_SIZE) \
		${NAME} 
	kops edit cluster ${NAME}
	kops update cluster ${NAME} --yes
	sleep 10m
	kops validate cluster
	kubectl config current-context
	kubectl get nodes

kops-create-efs:				## Create AWS EFS file system
	aws efs create-file-system \
		--creation-token $(FS_TOKEN) \
		--performance-mode maxIO \
		--region us-east-1

kops-create-mount-targers:		## Create mount targets
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	$(eval SEC_GROUP = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SecurityGroups[?GroupName==`nodes.${NAME}`]' --output text | head -n 1 | cut -d '	' -f1))	
	$(eval SUBNET = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SubnetId' --output text | head -n 1 | cut -f 1 ))
	aws efs create-mount-target \
		--file-system-id $(FS_ID) \
		--subnet-id $(SUBNET) \
		--security-group $(SEC_GROUP) \
		--region us-east-1 
	aws efs describe-mount-targets --file-system-id $(FS_ID)

kops-delete-mount-targets:		## Delete mount targets
		$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
		$(eval MOUNT_TARGET_ID = $(shell aws efs describe-mount-targets --file-system-id $(FS_ID) --query 'MountTargets[0].MountTargetId'))
		aws efs delete-mount-target --mount-target-id $(MOUNT_TARGET_ID) ; exit 0

kops-delete-file-system:		## Delete file system
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	aws efs delete-file-system --file-system-id $(FS_ID) ; exit 0

kops-delete-cluster:	## Delete k8s cluster
	make kops-delete-mount-targets
	sleep 1m
	make kops-delete-file-system 
	kops delete cluster --name ${NAME} --yes

kops-update-cluster:		## Update k8s cluster
	kops edit ig --name= nodes
	kops update cluster --yes
	kops rolling-update cluster --yes

kops-tips:		## Some kops commands
	# list clusters with
	kops get cluster
 	# edit this cluster with:
	kops edit cluster ${NAME} 
	# edit your node instance group
	kops edit ig --name=${NAME}  nodes
 	# edit your master instance group:
	kops edit ig --name=${NAME} master-us-east-1a
	# Finally configure your cluster with:
	kops update cluster --name saferwall.k8s.local --yes

kops-init-cert-manager: # Init cert-manager
	helm repo add jetstack https://charts.jetstack.io
	kubectl create namespace cert-manager
	kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
	kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.12/deploy/manifests/00-crds.yaml

saferwall: ## Deploy the cluster
	make awscli-install
	make aws-create-user
	make kops-install
	make kops-create-kops-bucket
	make kops-create-cluster
	make kops-create-efs
	make kops-create-mount-targers
	make kops-init-cert-manager
	make helm-install

# edit gp2 to: allowVolumeExpansion: true