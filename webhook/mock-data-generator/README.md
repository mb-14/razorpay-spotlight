### mock-data-generator
This golang code generates webhooks events and directly loads them into influxDB (http://localhost:8086)


### Instructions
- Make sure influx DB is running on your local (or run this on the droplet)
```
go build -o generator .
./generator -help

# To generate events for the last 10 days which are spaced out by 5000 milliseconds
./generator -mode=backfill -interval=5000 -duration=10

# To geneate failure events
./generator -mode=backfill -interval=5000 -duration=10 -event=payment_failed
``
