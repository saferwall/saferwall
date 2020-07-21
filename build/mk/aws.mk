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