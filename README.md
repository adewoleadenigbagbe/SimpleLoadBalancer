# A Simple Load Balancer

A Load balancer sits in front of a group of servers and routes client requests across all of the servers that are capable of fulfilling those requests. The intention is to minimise response time and maximise utilisation whilst ensuring that no server is overloaded.There are types of load balancer - Network Loadbalancer (Layer 4) and Application Loadbalancer (Layer 7). This project focuses on the Application load balancer which will route HTTP requests from clients to a pool of HTTP servers.

## Load Balancing Algorithm used in project
* Round Robin
* Random Weighted Round Robin
* Smooth Weighted Round Robin
* Ip Hashing (Ketama Ring)
* Least Connection
* Least ResponseTime
 
Check this link to read more about the algorithm type [LoadBalancingAlgorithm](https://aws.amazon.com/what-is/load-balancing/#:~:text=A%20load%20balancing%20algorithm%20is,fall%20into%20two%20main%20categories.)


## Project Setup
There are two folders in the project:

* LoadBalancer
* Server

The loadBalancer folder contains the services as well as the corresponding algorithm. you can configure the backend url, default algorithm type and weights in the **config.json** file, cd the loadbalancer folder and run the main.go file in cmd line

The server contains a rest api endpoint with crud operations to test the load balancing service
* Add Product
* Get Products
* Get Product by Id

The database in use is SQLite , so you make sure you have the SQLite3 service installed on your environment, any version preferably

Cd to the server folder , change **config.json** of the server folder based on the address and port. You need to run multiple instance of the server with address and port matching the **config.json** the loadbalancer folder.Use a client to test the endpoint and see the implementation of how the load balancing works


## Dependencies
SQLite3 Database




