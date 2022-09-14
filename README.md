# blog-memory-leaks
This repository contains sample code linked from a company blogpost.

# Running
To start the application, and execute some requests against it, run `docker-compose up -d`

# Metrics & Traces
Metrics are available in [local Prometheus](http://localhost:9090/graph?g0.expr=go_memstats_heap_alloc_bytes&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=15m)

Traces are available in [local Zipkin](http://localhost:9411/zipkin/)

