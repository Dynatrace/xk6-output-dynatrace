
# xk6-output-dynatrace
k6 extension for publishing test-run metrics to Dynatrace 


### Usage

To build k6 binary with the Prometheus remote write output extension use:
```
xk6 build --with github.com/henrikrexed/xk6-output-dynatrace@latest 
```

Then run new k6 binary with:
```
export K6_DYNATRACE_URL=http://<environmentid>.live.dynatrace.com 
export K6_DYNATRACE_APITOKEN=<Dynatrace API token>
the api token needs to have the scope: metric ingest v2
./k6 run script.js -o output-dynatrace
```


### On sample rate

k6 processes its outputs once per second and that is also a default flush period in this extension. The number of k6 builtin metrics is 26 and they are collected at the rate of 50ms. In practice it means that there will be around 1000-1500 samples on average per each flush period in case of raw mapping. If custom metrics are configured, that estimate will have to be adjusted.


