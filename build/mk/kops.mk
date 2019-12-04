kops-install:			## Install Kubernetes Kops
	curl -Lo kops https://github.com/kubernetes/kops/releases/download/$$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
	chmod +x ./kops
	sudo mv ./kops /usr/local/bin/
	kops version

aws-cli-install:		## Install aws cli tool
	sudo apt-get update
	sudo apt-get install awscli -y
	aws configure

kops-create-user:		## Create user to provision the cluster
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

kops-create-bucket:		## create s3 bucket
	aws s3api create-bucket --bucket kops-saferwall-com-state-store --region us-east-1
	aws s3api put-bucket-versioning --bucket kops-saferwall-com-state-store --versioning-configuration Status=Enabled

kops-create-cluster:	## create k8s cluster
	aws ec2 describe-availability-zones --region us-east-1
	kops create cluster --zones us-east-1a ${NAME}
	kops edit cluster ${NAME}
	kops update cluster ${NAME} --yes
	sleep 8m
	kops validate cluster
	kubectl get nodes

kops-create-efs:		## create AWS EFS file system
	aws efs create-file-system \
		--creation-token saferwall-efs \
		--performance-mode maxIO \
		--region us-east-1

kops-create-mount-targers:	## Create mount targets
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	$(eval SEC_GROUP = $(shell aws ec2 describe-instances --query 'Reservations[*].Instances[*].SecurityGroups[?GroupName==`nodes.saferwall.k8s.local`]' --output text | head -n 1 | cut -d '	' -f1))	
	$(eval SUBNET = $(shell aws ec2 describe-instances --query 'Reservations[0].Instances[0].SubnetId'))
	aws efs create-mount-target \
		--file-system-id $(FS_ID) \
		--subnet-id $(SUBNET) \
		--security-group $(SEC_GROUP) \
		--region us-east-1 
	aws efs describe-mount-targets --file-system-id $(FS_ID)

kops-delete-mount-targets:		## Delete mount targets
		$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
		$(eval MOUNT_TARGET_ID = $(shell aws efs describe-mount-targets --file-system-id $(FS_ID) --query 'MountTargets[0].MountTargetId'))
		aws efs delete-mount-target --mount-target-id $(MOUNT_TARGET_ID)

kops-delete-file-system:		## Delete file system
	$(eval FS_ID = $(shell aws efs describe-file-systems --query 'FileSystems[0].FileSystemId'))
	aws efs delete-file-system --file-system-id $(FS_ID)

kops-create-efs-provisioner:		## Create efs provisioner
	cd ~ \
		&& git clone https://github.com/kubernetes-incubator/external-storage \
		&& cd external-storage/aws/efs/deploy/ \
		&& kubectl apply -f rbac.yaml \
		
		# Modify manifest.yaml. In the configmap section change the
		# file.system.id: and aws.region: to match the details of the
		#  EFS you created. Change dns.name if you want to mount by your
		# own DNS name and not by AWS's *file-system-id*.efs.*aws-region*.amazonaws.com.
		# In the deployment section change the server: to the DNS endpoint of the EFS you created.
		&& kubectl apply -f manifest.yaml

kops-delete-cluster:		## Delete k8s cluster
	kops delete cluster --name ${NAME} --yes

kops-update-cluster:		## Update k8s cluster
	kops edit ig --name= nodes
	kops update cluster --yes
