# Airport API

<!-- My thought process and decisions goes here -->

---
_For tasks, checkout [tasks.md](tasks.md)_

---
Here I have dockerized the Go Backend Application
- Used minikube for the local k8s configuration
- have written terraform script to make  a bucket on s3
- Used ingress controller virtual server to split traffic between /airports and /airports_v2 API
- Wrote service and deployment files as well. Check the `k8s` folder. All Yaml files are thers

