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
To run the application type these commands in the k8s folder.

```bash
$ kubectl apply -f ./infra.yaml
$ kubectl apply -f ./services.yaml
$ kubectl apply -f ./gateway.yaml
```

To delete created resources type these commands in the k8s folder.

```bash
$ kubectl delete -f ./gateway.yaml
$ kubectl delete -f ./services.yaml
$ kubectl delete -f ./infra.yaml
```

## Access Application

To access Swagger UI go to [localhost:3000/api](http://localhost:3000/api).

## System Architecture Diagram

![Diagram](https://github.com/Marcel-MD/pad-labs/blob/main/diagram.png)
