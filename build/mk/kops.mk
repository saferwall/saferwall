kops-install:			## Install Kubernetes Kops
	curl -Lo kops https://github.com/kubernetes/kops/releases/download/$$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
	chmod +x ./kops
	sudo mv ./kops /usr/local/bin/
	kops version

aws-cli-install:		## Install aws cli tool
	sudo apt-get update
	sudo apt-get install awscli -y

kops-create-user:		## Create user to provision the cluster
	aws iam create-group --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops
	aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops
	aws iam create-user --user-name kops
	aws iam add-user-to-group --user-name kops --group-name kops
	aws iam create-access-key --user-name kops

kops-create-bucket:		## create s3 bucket
	aws s3api create-bucket --bucket kops-saferwall-com-state-store --region us-east-1
	aws s3api put-bucket-versioning --bucket kops-saferwall-com-state-store --versioning-configuration Status=Enabled

kops-create-efs:		## create AWS EFS
	aws efs create-file-system \
		--creation-token saferwall-efs \
		--performance-mode generalPurpose \
		--throughput-mode bursting \
		--region us-west-2 \
		--tags Key=Name,Value="Test File System" Key=developer,Value=rhoward \

kops-create-cluster:	## create k8s cluster
	aws ec2 describe-availability-zones --region us-west-2
	kops create cluster --zones us-west-2a ${NAME}
	kops edit cluster ${NAME}
	kops update cluster ${NAME} --yes
	sleep 5m
	kops validate cluster
	kubectl get nodes

kops-create-efs-provisioner:
	kubectl create configmap efs-provisioner \
		--from-literal=file.system.id=fs-47a2c22e \
		--from-literal=aws.region=us-west-2 \
		--from-literal=dns.name=""
		--from-literal=provisioner.name=example.com/aws-efs

kops-delete-cluster:		## delete k8s cluster
	kops delete cluster --name ${NAME} --yes
