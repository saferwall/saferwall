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
	echo "Copy the SecretAccessKey and AccessKeyID from the output."
	aws configure

KOPS_VERSION=1.21.0
kops-install:		## Install Kubernetes Kops
	curl -Lo kops https://github.com/kubernetes/kops/releases/download/v$(KOPS_VERSION)/kops-linux-amd64
	chmod +x ./kops
	sudo mv ./kops /usr/local/bin/
	kops version

kops-create-kops-bucket:		## Create s3 bucket for kops
	@echo "Creating S3 Bucket to store kops state"
	aws s3api create-bucket --bucket $(KOPS_STATE_S3_BUCKET_NAME) --region $(AWS_REGION)
	@echo "Enable bucket versioning"
	aws s3api put-bucket-versioning --bucket $(KOPS_STATE_STORE) --versioning-configuration Status=Enabled
	# Enable encryption
	# aws s3api put-bucket-encryption --bucket $(KOPS_STATE_STORE) --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

kops-create-cluster:			## Create k8s cluster
	kubectl config get-contexts
	aws ec2 describe-availability-zones --region $(AWS_REGION)
	kops create cluster \
		--master-count $(AWS_MASTER_COUNT) \
		--master-size $(AWS_MASTER_SIZE) \
		--master-zones $(AWS_MASTER_ZONES) \
		--node-count $(AWS_NODE_COUNT) \
		--node-size $(AWS_NODE_SIZE) \
		--zones $(AWS_NODE_ZONES) \
		--cloud aws \
		--topology private \
		--networking calico \
		--name ${KOPS_CLUSTER_NAME}
	kops edit cluster --name $(KOPS_CLUSTER_NAME)
	kops update cluster --name $(KOPS_CLUSTER_NAME) --yes
	sleep 10m
	kops export kubecfg --user
	kops validate cluster
	kubectl config current-context
	kubectl get nodes

kops-create-efs:				## Create AWS EFS file system
	aws efs create-file-system \
		--creation-token $(AWS_EFS_TOKEN) \
		--performance-mode maxIO \
		--region $(AWS_REGION)

kops-create-mount-targers:		## Create mount targets
	sudo apt update -qq && sudo apt install -qq jq
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	$(eval SEC_GROUP = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SecurityGroups[?GroupName==`nodes.${KOPS_CLUSTER_NAME}`]' --output text | head -n 1 | cut -d '	' -f1))	
	$(eval SUBNET = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].[SecurityGroups[0].GroupName==`nodes.${KOPS_CLUSTER_NAME}`, SubnetId]'  | jq '.[1][0][1]'  | tr -d '"'))
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
	sleep 30s
	make kops-delete-file-system
	kops delete cluster --name ${KOPS_CLUSTER_NAME} --yes

kops-update-cluster:		## Update k8s cluster
	kops edit ig --name= nodes
	kops update cluster --yes
	kops rolling-update cluster --yes

kops-cluster:			## Init ios cluster: create user, kops bucket,
	make kops-create-user
	make kops-create-kops-bucket
	make kops-create-cluster
	make kops-create-efs
	make kops-create-mount-targers

kops-tips:		## Some kops commands
	# list clusters with
	kops get clusters
 	# edit this cluster with:
	kops edit cluster ${KOPS_CLUSTER_NAME}
	# edit your node instance group
	kops edit ig --name=${KOPS_CLUSTER_NAME} nodes
 	# edit your master instance group:
	kops edit ig --name=${KOPS_CLUSTER_NAME} master-us-east-1a
	# Finally configure your cluster with:
	kops update cluster --name saferwall.k8s.local --yes

kops-export-yaml:  ### Export Kops cluster config
	kops get --name ${KOPS_CLUSTER_NAME} -o yaml > cluster-desired-config.yaml

saferwall: ## Deploy the cluster
	make awscli-install
	make kops-install
	make kops-cluster
	make helm-install
	make helm-add-repos
	make helm-update-dep
	make k8s-init-cert-manager
	# Install a release
	make helm-release
