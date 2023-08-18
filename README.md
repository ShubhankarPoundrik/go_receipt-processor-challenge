# Run the following commands

## To build the docker image

`docker build -t receipt-points-app .`

## To run unit tests

`docker run -it receipt-points-app go test -v`

## To start the server

`docker run -p 4000:4000 receipt-points-app`
