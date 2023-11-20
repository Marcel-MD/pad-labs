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

More information about [Kubernetes](https://kubernetes.io/).  
First of all you need to start redis cluster which will take some time.

```bash
$ kubectl apply -f ./redis-cluster.yaml
```

Once all the pods are initialized, you can see that Pod "redis-cluster-0" became the cluster master with the other nodes as slaves.

```bash
$ kubectl exec redis-cluster-0 -- redis-cli cluster nodes
```

You can also check the slot configuration.

```bash
$ kubectl exec redis-cluster-0 -- redis-cli --cluster check localhost 6379
```

To run everything else type these commands one by one in this exact order.

```bash
$ kubectl apply -f ./rabbitmq.yaml
$ kubectl apply -f ./postgres.yaml
$ kubectl apply -f ./services.yaml
$ kubectl apply -f ./gateway.yaml
```

To delete created resources type these commands.

```bash
$ kubectl delete -f ./gateway.yaml
$ kubectl delete -f ./services.yaml
$ kubectl delete -f ./postgres.yaml
$ kubectl delete -f ./rabbitmq.yaml
```

To delete redis cluster resources type these commands.

```bash
$ kubectl delete service,statefulsets redis-cluster
$ kubectl delete configmaps redis-cluster-config
$ kubectl delete poddisruptionbudgets.policy redis-cluster-pd
$ kubectl delete pod redis-cluster-0 redis-cluster-1 redis-cluster-2 redis-cluster-3 redis-cluster-4 redis-cluster-5
```

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
