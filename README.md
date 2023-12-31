# Laboratory Works on Distributed Systems Programming

## Description

E-commerce platform that allows users to buy and sell products. This idea's implementation through distributed systems is suitable because it can be easily separated in distinct components that can be developed independently. Similar platforms that use microservices architecture are Amazon, eBay, AliExpress, etc.

## Run Application with Docker

More information about [Docker](https://www.docker.com/).  
To run the application type this command in the root folder.

```bash
$ docker compose up
```

You might have to run this command twice if it doesn't work the first time :)

## Run Application with Kubernetes

More information about [Kubernetes](https://kubernetes.io/). Before running next commands move to `k8s` directory.

```bash
$ cd ./k8s
```

First of all you need to start Redis cluster which will take some time.

```bash
$ kubectl apply -f ./redis-cluster.yaml
```

After that you can use these two commands to check if the cluster has been initialized correctly.

```bash
$ kubectl exec redis-cluster-0 -- redis-cli cluster nodes
$ kubectl exec redis-cluster-0 -- redis-cli --cluster check localhost 6379
```

To start Postgres cluster we'll be using [Kubegres](https://www.kubegres.io/doc/getting-started.html).  
First we have to install Kubegres Custom Resource Definitions.

```bash
$ kubectl apply -f https://raw.githubusercontent.com/reactive-tech/kubegres/v1.17/kubegres.yaml
$ kubectl get all -n kubegres-system
```

Then we can start Postgres cluster which may take some time.

```bash
$ kubectl apply -f postgres-cluster.yaml
```

To run everything else type these commands one by one in this exact order.

```bash
$ kubectl apply -f ./rabbitmq.yaml
$ kubectl apply -f ./postgres.yaml
$ kubectl apply -f ./services.yaml
$ kubectl apply -f ./gateway.yaml
```

For performance monitoring we'll be using [Prometheus Operator](https://prometheus-operator.dev/).  
The first step is to install the operator's Custom Resource Definitions as well as the operator itself with the required RBAC resources.

```bash
LATEST=$(curl -s https://api.github.com/repos/prometheus-operator/prometheus-operator/releases/latest | jq -cr .tag_name)
curl -sL https://github.com/prometheus-operator/prometheus-operator/releases/download/${LATEST}/bundle.yaml | kubectl create -f -
```

To start Prometheus run the following command.

```bash
$ kubectl apply -f ./prometheus-operator.yaml
```

## Cleanup Kubernetes Resources

To delete created resources type these commands.

```bash
$ kubectl delete -f ./prometheus-operator.yaml
$ kubectl delete -f ./gateway.yaml
$ kubectl delete -f ./services.yaml
$ kubectl delete -f ./postgres.yaml
$ kubectl delete -f ./rabbitmq.yaml
$ kubectl delete -f ./redis-cluster.yaml
$ kubectl delete kubegres warehouse
```

To also delete Kubegres Custom Resource Definitions.

```bash
$ kubectl delete -f https://raw.githubusercontent.com/reactive-tech/kubegres/v1.17/kubegres.yaml
```

I don't know how to remove Prometheus Operator.

## Use the Application

0. Run the application as mentioned above.
1. Access Swagger UI at [localhost:3010/api](http://localhost:3010/api).
2. Register a user with the `POST /users/register` endpoint.
3. Login with the `POST /users/login` endpoint. (optional)
4. Copy the access token from the response body.
5. Set the access token in the `Authorize` button from the top right corner of the Swagger UI.
6. Create a product with the `POST /products` endpoint.
7. Order a product with the `POST /orders/:productId` endpoint.

You can find the internal endpoints specification in `/proto` folder.

## Monitor the Application

0. Access Prometheus at [localhost:9090](http://localhost:9090).
1. Access Grafana at [localhost:3000](http://localhost:3000).
2. Login using credentials: `admin` and `password`.
3. Open the `Dashboards` tab from the left side menu.
4. Select `Microservices Metrics` dashboard.

## System Architecture Diagram

![Diagram](https://github.com/Marcel-MD/pad-labs/blob/main/diagram.png)

## Redis Cluster Diagram

<img src="https://github.com/Marcel-MD/pad-labs/blob/main/redis-cluster.png" alt="Redis Cluster Diagram" width="600"/>

## Postgres Cluster Diagram

<img src="https://github.com/Marcel-MD/pad-labs/blob/main/postgres-cluster.png" alt="Postgres Cluster Diagram" width="400"/>
