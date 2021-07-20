# echoapi

# echoapi
A simple HTTP API that echos back any valid JSON input.

## Table Of Contents

* [Intro](#intro)
* [API documentation](#api-documentation)
* [Authentication](#authentication)
* [Monitoring](#monitoring)
* [Configuration](#configuration)
* [Development](#development)
  * [Requirements](#requirements)
  * [How to build](#how-to-build)
  * [How to run tests](#how-to-run-tests)
* [Improvements and limitations](#improvements-and-limitations)

## Intro
The HTTP API echoes back any valid JSON while adding a top-level 'echoed: true' field. If the field is already preset and set to 'true', it will return a 400 HTTP error.  
It also exposes Prometheus metrics at the /metrics endpoint, see documentation below for details.

## API Documentation
### Endpoints:
#### Echo  
`POST|PUT /api/echo`  
* Protected with basic auth
* Returns the POST'd JSON with an added top-level field 'echoed: true'
* If the field is already set to true it returns HTTP 400
* If the JSON is not valid and cannot be unmarshalled it returns HTTP 400  
  
Sample responses:
```
curl -i -u bob:bob123 -d '{"username":"xyz","upload":"xyz"}' http://localhost:9090/api/echo
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 17 Jul 2021 07:11:28 GMT
Content-Length: 50

{"echoed":"true","upload":"xyz","username":"xyz"}
```
```
curl -i -u bob:bob123 -d '{"username":"xyz","upload":"xyz","echoed":"true"}' http://localhost:9090/api/echo
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sat, 17 Jul 2021 07:11:53 GMT
Content-Length: 58

{"error":"Top level echoed field is already set to true"}
```
```
curl -i -u bob:bob123 -d '{"username":"xy' http://localhost:9090/api/echo
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sat, 17 Jul 2021 07:12:11 GMT
Content-Length: 56

{"error":"Malformed request, input must be valid JSON"}
```
```
curl -i -u bob:john -d '{"username":"xyz","upload":"xyz"}' http://localhost:9090/api/echo
HTTP/1.1 401 Unauthorized
Content-Type: text/plain; charset=utf-8
Www-Authenticate: Basic realm="Restricted"
X-Content-Type-Options: nosniff
Date: Sat, 17 Jul 2021 07:28:00 GMT
Content-Length: 14

Unauthorized.
```
#### Metrics  
`GET /metrics`  
* Returns metrics in Prometheus format:
  * echoapi_http_duration_seconds_bucket for latency information.
  * echoapi_http_requests_total for requests counts by code, method and path.  
  
Sample response:
```
curl -s http://localhost:9090/metrics | grep echoapi
# HELP echoapi_http_duration_seconds Duration of HTTP requests.
# TYPE echoapi_http_duration_seconds histogram
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.005"} 0
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.01"} 0
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.025"} 0
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.05"} 0
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.1"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.25"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="0.5"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="1"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="2.5"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="5"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="10"} 8
echoapi_http_duration_seconds_bucket{path="/api/echo",le="+Inf"} 8
echoapi_http_duration_seconds_sum{path="/api/echo"} 0.473520726
echoapi_http_duration_seconds_count{path="/api/echo"} 8
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.005"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.01"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.025"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.05"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.1"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.25"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="0.5"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="1"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="2.5"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="5"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="10"} 2
echoapi_http_duration_seconds_bucket{path="/metrics",le="+Inf"} 2
echoapi_http_duration_seconds_sum{path="/metrics"} 0.000845648
echoapi_http_duration_seconds_count{path="/metrics"} 2
# HELP echoapi_http_requests_total How many HTTP requests processed by status code, method and HTTP path.
# TYPE echoapi_http_requests_total counter
echoapi_http_requests_total{code="200",method="GET",path="/metrics"} 2
echoapi_http_requests_total{code="200",method="POST",path="/api/echo"} 6
echoapi_http_requests_total{code="400",method="POST",path="/api/echo"} 2
```
#### Status Codes
| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 400 | `BAD REQUEST` |
| 401 | `UNAUTHORIZED` |
| 404 | `NOT FOUND` |

## Authentication
The API supports HTTP basic auth via an htpasswd file, for which it supports SHA1 (insecure), Apache salted MD5 or Bcrypt (preferred).  
The path to the file can be configured with the HTPASSWD_PATH environment varaible.  
An example file can be found at [conf/htpasswd](conf/htpasswd)

## Monitoring
The API exposes metrics in Prometheus format at the /metrics endpoint. See the API documentation above for details.

## Configuration
The following environment variables can be set to configure the application:
* `LISTEN_ADDR` - Bind address for the webserver (IP:PORT) [*optional*] [default: :9090]
* `LOG_LEVEL` - Log level (INFO, WARN, DEBUG, TRACE) [*optional*] [default: INFO]
* `HTPASSWD_PATH` - Path to the htpasswd file [*optional*] [default: ./conf/htpasswd]

## Development

### Requirements
* Golang (~>1.15 for local build)
* Docker

### How to build
Build locally, binary written to bin/echoapi:
```
make build
```

Build with Docker:
```
make build-docker
```

Tag and push to the Docker registry:
```
make distribute
```

### How to run tests
To run unit / integration tests:
```
make test
```

### How to run locally
Without Docker:
```
make build
LOG_LEVEL=DEBUG ./bin/echoapi
```

With Docker:
```
make build-docker
docker run --rm -p 8080:8080 -v $PWD/conf/htpasswd:/var/lib/echoapi/htpasswd \
  -e "LISTEN_ADDR=:8080" -e "HTPASSWD_PATH=/var/lib/echoapi/htpasswd" docker.io/gpaggi/echoapi:latest
```


## Improvements and limitations
* Gorilla Mux can be replaced with [httprouter](https://github.com/julienschmidt/httprouter) if performances are of concern.
* Only access logs are being logged to Stdout. Service logs should be implemented for better monitoring.
* The server currently supports only HTTP. HTTPS should be implemented.
* The /metrics endpoint should be exposed on an internal admin port, not to leak information to the outside.
