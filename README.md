# devops-test app
Simple sample app that uses a REST API for storing and retrieving values written in Golang.

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)

## General info

This application exposes the following HTTP-based APIs: 

1. Description: Saves/updates the given user’s name and date of birth in the database.
`Request: PUT /hello/<username> { “dateOfBirth”: “YYYY-MM-DD” }
Response: 204 No Content`
Note:
<username> must contain only letters.
YYYY-MM-DD must be a date before the today date.
  
2. Description: Returns hello birthday message for the given user
`Request: Get /hello/<username>
Response: 200 OK`

Response Examples:
A. If username’s birthday is in N days:
`{ “message”: “Hello, <username>! Your birthday is in N day(s)”}`

Response Examples:
B. If username’s birthday is today:
`{ “message”: “Hello, <username>! Happy birthday!” }`


## Technologies
  
This simple application is written in Go and expects a MySQL database for storing data.

## Setup
 
After cloning this repository, if you want to run it locally you can do it in two ways:
  
1. Through podman-compose:
  `$ podman-compose up -d`
  
  For this you would need to build it locally by running either `docker build -t jrosental/devops-test .` or `podman build -t jrosental/devops-test`.
  Also bear in mind that in the docker-compose.yaml file you would need to adjust the environment variable `DBHOST` to point to your local machine IP address.
  
2. Standalone by building the image locally and passing the address of another database instance through the `DBHOST` environment variable
  
### Environment variables:
* **DBHOST** in the format `<IP-address>:<port>`, e.g: _127.0.0.1:3306_
* **DBUSER** which is the user for connect to the database.
* **DBPASS** which is the password for the user to connect to the database.
