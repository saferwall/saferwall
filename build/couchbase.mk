couchbase-run:		## Run couchbase docker container instance.
	docker run -d --name db -p 8091-8094:8091-8094 -p 11210:11210 couchbase/server

couchbase-start:	## Run exiting couchbase `db` container.
	docker start db

couchbase-k8s-init:	## Init kubernetes operator ops
	kubectl create -f crd.yaml
	kubectl create -f cluster-role-sa.yaml
	kubectl create -f cluster-role-user.yaml
	kubectl create serviceaccount couchbase-operator --namespace default
	kubectl create clusterrolebinding couchbase-operator --clusterrole couchbase-operator --serviceaccount default:couchbase-operator
	kubectl create -f operator.yaml
	kubectl create -f secret.yaml
	cbopctl create -f couchbase-cluster.yamG STATUS
	kubectl delete deployment couchbase-operator
