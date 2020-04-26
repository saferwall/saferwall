elastic-drop-db:		## Delete all indexes.
	curl -X DELETE 'http://localhost:9200/_all'