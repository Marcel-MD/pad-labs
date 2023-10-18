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
To run the application type these commands one by one in this exact order in the k8s folder.

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

## Use the Application

0. Run the application as mentioned above.
1. Access Swagger UI at [localhost:3000/api](http://localhost:3000/api).
2. Register a user with the `POST /users/register` endpoint.
3. Login with the `POST /users/login` endpoint. (optional)
4. Copy the access token from the response body.
5. Set the access token in the `Authorize` button from the top right corner of the Swagger UI.
6. Create a product with the `POST /products` endpoint.
7. Order a product with the `POST /orders/:productId` endpoint.

You can find the internal endpoints specification in `/proto` folder.

## System Architecture Diagram

![Diagram](https://github.com/Marcel-MD/pad-labs/blob/main/diagram.png)
