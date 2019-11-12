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
	aws efs describe-file-systems

kops-create-mount-targers:
	aws ec2 describe-instances --query 'Reservations[*].Instances[*].{Instance:InstanceId,Subnet:SubnetId,SecurityGroups:SecurityGroups}'
	aws efs create-mount-target \
	--file-system-id fs-04f47285 \
	--subnet-id subnet-0806281d75ce0c581 \
	--security-group sg-0bc97a80b78efb528 \
	--region us-east-1 
	aws efs describe-mount-targets --file-system-id fs-04f47285

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


kops-delete-cluster:		## delete k8s cluster
	kops delete cluster --name ${NAME} --yes

kops-update-cluster:
	kops edit ig --name= nodes
	kops update cluster --yes



# aws  ec2 describe-vpcs    



