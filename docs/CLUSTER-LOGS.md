# VIEWING THE CLUSTER LOGS

- Port forward kibana service: `make k8s-pf-kibana`
- Check if Kibana's dashboard is up and running: http://localhost:5601/
- Next, we need to create an index, go to: http://localhost:5601/app/management/kibana/indexPatterns/create
- In the `index pattern name field`, select `filebeat-*`:
![CreateIndex](https://i.imgur.com/VfhpyWr.png)
- Afterwards, you will be asked to select a time field for use with the global time filter:
![CreateIndexTimestamp](https://i.imgur.com/0BLAp9k.png)
- Click on the left sidebar, then select `Discover` and that's pretty much it. You can select which fields you want to view on the table, together with some filters.
- For example, we have a saved query to see the workers (or consumers) only, by using this query: `container.labels.io_kubernetes_container_name: "consumer"`
- Then you can filter only the columns you need:
![Logs](https://i.imgur.com/2s2Tn5L.png)
