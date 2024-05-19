# Logger Service

This allows different services to log different events from across the microservices cluster. This however does not dirrectly interface with the direct internet connection. Rather, they are used internally to log events in the application.

## Stacks and Frameworks

1. Chi Framework
2. MongoDB
3. json/RPC/gRPC - either of these is used as the data transfer protocol

## Connection Settings for the database

The database needs setup credentials as shown below:

```.env

MONGODB_PASSWORD=admin
MONGODB_USERNAME=password
MONGODB_URI="mongodb://admin:password@db-mongo:27017"
```

> Note: the `db-mongo` in the `MONGODB_URI` above is thesame as the name of the mongoDB instance in the docker compose `yml` file. In development, we can as well use `localhost`

## Suggestions/Improvement

1. Refactor the `helpers.go` file into a package which can eaily be used across all the microservices as needed

2. Refactor the data model methods so that you don't have to recreate the same logic for every entity
