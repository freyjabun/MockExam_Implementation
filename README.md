# DISYS-Auction

Submission of solution for Increment service made by Nadja Brix Koch, nako.

## User manual

### Starting the server

1. Open a terminal and `cd` to the project's directory.

2. Run the following commands:

    `$ docker build -t increment --no-cache .`

    `$ docker-compose build`

    `$ docker compose up`

### Connecting as a client

1. To connect to the auction service, open a terminal and `cd` to the project's directory.

2. Run command

    `$ go run client/client.go`

3. Enjoy the incrementor.
