awscli-install:		## Install aws cli tool and configure it
	sudo apt install curl python -y
	curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
	unzip -o awscli-bundle.zip
	sudo ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws
	aws --version
	rm awscli-bundle.zip
	@echo "========================================================================================="
	@echo "log in to aws console and get your access and secret key, for more information, consult:"
	@echo "https://aws.amazon.com/blogs/security/wheres-my-secret-access-key/"
	aws configure

kops-create-user:	## Create the kops IAM user to provision the cluster
	# create kops group 
	aws iam create-group --group-name kops
	# attach permissions
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess --group-name kops
	# create kops user
	aws iam create-user --user-name kops
	aws iam add-user-to-group --user-name kops --group-name kops
	aws iam create-access-key --user-name kops

kops-install:		## Install Kubernetes Kops
	curl -Lo kops https://github.com/kubernetes/kops/releases/download/v1.16.1/kops-linux-amd64
	chmod +x ./kops
	sudo mv ./kops /usr/local/bin/
	kops version

kops-create-kops-bucket:		## Create s3 bucket for kops
	aws s3api create-bucket --bucket $(KOPS_STATE_STORE) --region $(AWS_REGION)
	# Enable versioning
	aws s3api put-bucket-versioning --bucket $(KOPS_STATE_STORE) --versioning-configuration Status=Enabled
	# Enable encryption
	# aws s3api put-bucket-encryption --bucket $(KOPS_STATE_STORE) --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

kops-create-cluster:			## Create k8s cluster
	kubectl config get-contexts
	aws ec2 describe-availability-zones --region $(AWS_REGION)
	kops create cluster \
		--zones us-east-1a \
		--node-count $(AWS_NODE_COUNT) \
		--node-size $(AWS_NODE_SIZE) \
		--master-size $(AWS_MASTER_SIZE) \
		--cloud aws \
		--name ${KOPS_CLUSTER_NAME} 
	kops edit cluster --name $(KOPS_CLUSTER_NAME)
	kops update cluster --name $(KOPS_CLUSTER_NAME) --yes
	sleep 10m
	kops validate cluster
	kubectl config current-context
	kubectl get nodes

kops-create-efs:				## Create AWS EFS file system
	aws efs create-file-system \
		--creation-token $(AWS_EFS_TOKEN) \
		--performance-mode maxIO \
		--region $(AWS_REGION)

kops-create-mount-targers:		## Create mount targets
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	$(eval SEC_GROUP = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SecurityGroups[?GroupName==`nodes.${KOPS_CLUSTER_NAME}`]' --output text | head -n 1 | cut -d '	' -f1))	
	$(eval SUBNET = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SubnetId' --output text | head -n 1 | cut -f 1 ))
	aws efs create-mount-target \
		--file-system-id $(FS_ID) \
		--subnet-id $(SUBNET) \
		--security-group $(SEC_GROUP) \
		--region $(AWS_REGION)
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
	kops delete cluster --name ${KOPS_CLUSTER_NAME} --yes

kops-update-cluster:		## Update k8s cluster
	kops edit ig --name= nodes
	kops update cluster --yes
	kops rolling-update cluster --yes

kops-tips:		## Some kops commands
	# list clusters with
	kops get cluster
 	# edit this cluster with:
	kops edit cluster ${KOPS_CLUSTER_NAME} 
	# edit your node instance group
	kops edit ig --name=${KOPS_CLUSTER_NAME}  nodes
 	# edit your master instance group:
	kops edit ig --name=${KOPS_CLUSTER_NAME} master-us-east-1a
	# Finally configure your cluster with:
	kops update cluster --name saferwall.k8s.local --yes

helm-init-cert-manager: # Init cert-manager
	# Create the namespace for cert-manager.
	kubectl create namespace cert-manager
	# Install the CustomResourceDefinition
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.2/cert-manager.crds.yaml

saferwall: ## Deploy the cluster
	make awscli-install
	make kops-install
	make kops-create-user
	make kops-create-kops-bucket
	make kops-create-cluster
	make kops-create-efs
	make kops-create-mount-targers
	# Building docker containers
	make backend-release
	make frontend-release
	make consumer-release
	make multiav-release
	make multiav-release-go
	# At this stage, all containers are ready
	make helm-install
	make helm-add-repos
	make helm-init-cert-manager
