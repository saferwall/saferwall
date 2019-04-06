# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change. 

Please note we have a code of conduct, please follow it in all your interactions with the project.

### Repostiory Layout
* __pkg__ : library code to use by external applications.
* __website__ : saferwall website and documentation.
* __ui__ : (frontend) vue.js dashboard.
* __web__ : (backend) go web application.
* __build__ : docker files, makefiles, kubernetes deployements.
* __api__ : proto buffer specs.

### Requirements

- Install docker:
```
sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce -y
```