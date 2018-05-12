# Search Microservice Demo

## Getting started
To run this all you need to do is:

Build the docker-compose file with `docker-compose build`

And then run it with `docker-compose up`

## Populating Data

http://localhost:8080/populate?number={number of items to be added into elasticsearch}

## Query data

http://localhost:8080/search?q={query}&from={first item index}&size={number of items}
