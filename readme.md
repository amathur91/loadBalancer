# RGate #

This is a minimalistic config driven load balancer for docker containers implemented in Go.
Configuration of the services is defined in yml file and passed to the 
service during bootstrap. Load balancer scans the running containers after every X minutes.
X is defined in **knobs.go**.
Once the load balancer recognizes during scanning that particular services last alive response is beyond the threshold,
it marks the service status as Inactive and further calls for the impacted service will get default response.


Features:
* YML driven service declaration
* Docker container are identified with Labels
* Statistics for the API calls with P95/P99 
* Automatically shuts off the calls when containers are not available
* Services are activated when containers are available and ready.
* API Service to Docker Mapping is driven by config YML
* Default Response with HTTP Code is based on YAML config

## Sample Config YAML
```yml
backends:
- name: orders
  match_labels:
    app_name: orders
    env: production
- name: payment
  match_labels:
    app_name: payment
    env: production
default_response:
  body: "“This is not reachable”"
  status_code: 403
routes:
  - backend: orders
    path_prefix: /api/orders
  - backend: payment
    path_prefix: /api/payment
```

## Usage
* Define the YAML with your backend services, default response and path to service mapping.
* Use the compiled binary as => 
```bash
sudo ./loadbalancer --config config.yml --port 8081
```

## Steps to build
```go
go build
```

## Configs
* All the configurable parameters are in knobs like time interval for service discovery, threshold for service alive


## Local Testing
For the purpose of local testing. Custom docker images of nginx are added to this repo.
These images include additional labels which are used for service discovery.
* Build the docker images in orders and payment folders
```bash
sudo docker build -t orders .
sudo docker build -t payment .
```

* Docker Images have labels that are used for mapping with config.yml file.
```Dockerfile
FROM nginx
LABEL app_name=orders
LABEL env=production
```

* Now run the docker containers.
```bash
sudo docker run --rm -d -p 8082:80 --name orders orders
sudo docker run --rm -d -p 8083:80 --name payment payment
```

* Run the service using the local config.yml provided.
```bash
sudo ./loadbalancer --config config.yml --port 8081
```

* Use Postman or Curl to hit the below Endpoint
```bash
curl -v http://localhost:8081/api/orders
curl -v http://localhost:8081/api/payment
curl -v http://localhost:8081/stats
```

* When the containers are up the above API will return a 200 Response.

* Now stop any one of the container.

* As the service scans the services every 1 minute(Configurable in Knobs).
  Initially the response to the impacted service is 500.
  After 1 minutes the endpoint will start returning 403 which is the default response mentioned in config.yml
  
## Logs
* *rgate.txt* is the file where all the logs are appended. 