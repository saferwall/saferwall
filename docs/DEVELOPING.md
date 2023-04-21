# Developing

- Locate the `.env` file and update the `EXTERNAL_IP` to your host local IP address. This variable is used during docker-compose for NSQ to easily consume jobs while you are developing code outside of the containers.
- Run `make dc-up SVC="sandbox"`
- Go to: http://localhost:8091/ui/index.html:
    - Create a new cluster named saferwall (default settings).
    - Create a bucket named `sfw` with default settings.
- As root: `chmod 777 -R /samples`
