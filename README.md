# Serverless Go Project (POC)

A serveless go project that provides two services:
- Finds the n'th element of Fibonacci sequence
- Checks if a number is prime

## FaaS built Go & Dockerized

In order to launch the project, you need to install docker on your machine. 
Informations can be found on:

https://docs.docker.com/install/

Once you have installed docker, open the terminal and run the following commands
on the root of the project:

	docker-compose build --no-cache
    
and then:

    docker-compose up --force-recreate
   
After you see "Gateway On" on your terminal, you can access the server on your machine, at your localhost in port 8080.

Provide service and parameter as: 
    
    /lambda/{service}?a-param={parameter}

#### Services:
- **fibonacci**: *finds the "a" element of fibonacci sequence*
- **prime**: *checks if "a" is prime*