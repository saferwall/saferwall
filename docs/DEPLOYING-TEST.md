# Using Saferwall from a Vagrant Box

The steps below describes how to setup a vagrant box that have `saferwall` kubernetes cluster deployed on it.

 The `virtualbox` VM will requires 6GB of RAM and it is about 11GB of disk.

## 1. Download VirtualBox

- Download and install VirtualBox from the following url: https://www.virtualbox.org/wiki/Downloads

## 2. Download Vagrant

- Download and install Vagrant from the following url: https://www.vagrantup.com/downloads

## 3. Getting the Vagrant box

- Open a shell and run: `vagrant init saferwall/saferwall`
- Then: `vagrant up`