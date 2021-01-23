# gw-pool
A micro service for real time data persistence

## Deployment

There are two ways to run this service, via `localhost` or using `docker-compose`. Even so all required commands are located into the Makefile

### Localhost
1. Make sure all environment variables are properly set into the `.env` file. Especially the __DATABASE_DNS__ points to `localhost`.
1. Both server and client needs an environment variable __DOTENV_PATH__ to know where to find the `.env` file. The Makefile sets it automatically, but be aware if using your own dotenv file.
1. Execute the server in a terminal using the `make run` command.
1. In a new terminal you can deploy and undeploy the database using the `make deploy` and `make undeploy` respectively. 

### Docker
1. Make sure all environment variables are properly set into the `.env` file. Especially the __DATABASE_DNS__ points to `postgres` instead of `localhost`.
1. Build the image running into the terminal the command `make build`.
1. Make sure the `docker-compose` file has the backend specifications uncommented.
1. Deploy the server and database using the `make deploy` command. The `make undeploy` will undeploy both containers as well.

## Setup

Once server and database are both running, it's all ready for sending requests to the server. The command `make ping` will send a dummy request for you. However, this will fail if the database do not have the expected table: `locations`.

To make it so, into the `sql` folder are provided both sql scrips to create the database and drop it. Also, the docker-compose deploys an instance of `pgadmin` in order to manage de database and execute such scrips.

When connecting to the database through `pgadmin` remember that all credentials are located into the `.env` file and the database's host is `postgres`. 