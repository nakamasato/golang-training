# Expose Prometheus metrics

https://prometheus.io/docs/guides/go-application/

Install dependencies:

```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

Run app

```
go run main.go
```

Check

```
curl http://localhost:2112/metrics
```

```
curl http://localhost:2112/metrics | grep myapp_processed_ops_total
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  5159    0  5159    0     0  1679k      0 --:--:-- --:--:-- --:--:-- 1679k
# HELP myapp_processed_ops_total The total number of processed events
# TYPE myapp_processed_ops_total counter
myapp_processed_ops_total 87
```
