# Monitoring Service

> ⚠️ **Disclaimer**: This is a **fork** of [Monitoring Service](https://github.com/Dentanoid/Monitoring-Service), originally created and maintained by the [Dentanoid Organization](https://github.com/Dentanoid)

Welcome to the Monitoring Service! This service handles logging, statistics and monitoring.

## Getting started

This service is written in Go. [Check this link for more information about GO.](https://go.dev/)

To run this service you need to follow the steps described below:

### Installing GO using BREW (if you dont have GO)

If you do not have GO installed on your computer you can download both brew and GO with these commands:

#### Install brew
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
``````

#### Install GO with brew
```
brew install go
``````

### Add .env file
The .env file contains information about the MQTT broker and mongoDB. This informatin is best contained locally on your computer, to keep your connections private. You will have to insert a MONGO_URI for your database and a BROKER_URL.

For our instances of the service, we used a [HIVE](https://www.hivemq.com/mqtt/) private broker.

```
MONGO_URI = "YOUR_URI"

BROKER_URL = "YOUR_BROKER:PORT_NR"
```

### Run Monitoring service
In order to build and run the monitoring service you need to type these commands in to your terminal:


```
go build
go run main.go
```
Congratulations! You are now running the Monitoring service.

## Roadmap
This service will not get updated in the future, due to project being considered as closed when GU course DIT356 is finished.


## Authors and acknowledgment

- Lucas Holter
- Cornelia Olofsson Larsson 
- James Klouda 
- Jonatan Boman 
- Mohamad Khalil
- Joel Mattson 

## Project status
The service may recieve updates until 9th January 2024, and none after.
