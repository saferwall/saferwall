certbot-install		## Install Certbot for nginx
	sudo add-apt-repository ppa:certbot/certbot
	sudo apt-get update
	sudo apt-get install python-certbot-nginx
	sudo certbot --nginx -d saferwall.com -d www.saferwall.com --email admin@saferwall.com --agree-tos
