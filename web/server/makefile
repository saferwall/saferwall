# Install glide
install-glide:
	curl https://glide.sh/get | sh


# Download and install Docker
install-docker:
	sudo apt-get update
	sudo apt-get install apt-transport-https ca-certificates \
    	curl software-properties-common	-y
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo apt-key fingerprint 0EBFCD88
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $$(lsb_release -cs) stable"
	sudo apt-get update
	sudo apt-get install docker-ce -y

# Install Couchbase in a docker container and run it
init-cb:
	sudo docker run -d --name db -p 8091-8094:8091-8094 -p 11210:11210 couchbase

# Run exiting container
start-cb:
	sudo docker start db