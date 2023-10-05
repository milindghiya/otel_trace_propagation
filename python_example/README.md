## Python OPENTELEMETRY TRACE PROPAGATION

This repository consists of three services: serviceA, serviceB, and serviceC. serviceA acts as the main interaction point, exposing an HTTP endpoint. It uses OpenTelemetry SDK to generate traceIds for better observability and makes subsequent requests to both serviceB and serviceC.

## Service Interactions
```
Client -> serviceA (HTTP request)
serviceA -> serviceB (HTTP request using net/http and resty)
serviceA -> serviceC (gRPC request)
```


## Setup & Installation

### Prerequisites

Ensure you have the following installed:
1. Python3
2. gRPC tools
3. virtual environment

### Installation

```
pip install -r requirements.txt
```


## Run

```
# Terminal 1
python3 service_c.py

# Terminal 2
python3 service_b.py

# Terminal 3
python3 service_a.py

curl http://localhost:8080
```



