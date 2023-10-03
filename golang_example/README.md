## GOLANG OPENTELEMETRY TRACE PROPAGATION

This repository consists of three services: serviceA, serviceB, and serviceC. serviceA acts as the main interaction point, exposing an HTTP endpoint. It uses OpenTelemetry SDK to generate traceIds for better observability and makes subsequent requests to both serviceB and serviceC.

## Service Interactions

Client -> serviceA (HTTP request)
serviceA -> serviceB (HTTP request using net/http and resty)
serviceA -> serviceC (gRPC request)



## Setup & Installation

### Prerequisites

Ensure you have the following installed:
1. GoLang (version 1.20+)
2. gRPC tools


### Installation

```
go mod tidy
```


## Build

```
make all
```


## Run

```
# Terminal 1
./bin/servicea

# Terminal 2
./bin/serviceb 

# Terminal 3
./bin/servicec

curl http://localhost:8080
```



