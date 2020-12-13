### Grafana

http://razorpay-spotlight.in:3000/d/I4PoSO1Mz/razorpay-metrics?orgId=1&refresh=5s

The json file present is a dump of the json model of the dashboard above

### Auth less dashboard access

Modify the following lines in the `/etc/grafana/grafana.ini`

```
[auth.anonymous]
enabled = true

org_name = Main Org.

org_role = Viewer
```

`systemctl restart grafana-server`
