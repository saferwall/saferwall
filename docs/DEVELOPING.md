# Developing

- Install docker and docker compose.
- Install make command
- Locate the `.env` file and update the `EXTERNAL_IP` to your host local IP address. This variable is used during docker-compose for NSQ to easily consume jobs while you are developing code outside of the containers.
- Run `make dc-up SVC="sandbox"`
- Go to: http://localhost:8091/ui/index.html:
    - Create a new cluster named saferwall (default settings).
    - Create a bucket named `sfw` with default settings.
- As root: `mkdir /sample && chmod 777 -R /samples`
- Install a browser extension that disables CORS.
- Now you can sign up in: http://localhost:8000/
