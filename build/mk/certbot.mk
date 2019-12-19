certbot-install:	## Install Certbot for nginx
	sudo add-apt-repository ppa:certbot/certbot
	sudo apt-get update
	sudo apt-get install python-certbot-nginx -y
	sudo certbot --nginx -d saferwall.com -d www.saferwall.com --email admin@saferwall.com --agree-tos


	sudo git clone https://github.com/letsencrypt/letsencrypt /opt/letsencrypt
	cd /opt/letsencrypt
	sudo -H ./letsencrypt-auto certonly --standalone -d dev.api.saferwall.com -d www.dev.api.saferwall.com
