# CQRS Demo with Redpanda

A CQRS pattern implementation go with Redpanda as eventing persistent store. The demo uses PostgreSQL for writes and MongoDB for reads.

## Setup

### Environment

Create a file called `.envrc` and add the following variables to it,

```shell
export RPK_BROKERS="127.0.0.1:19092"
export CONSUMER_GROUP_ID=cqrs-todo
export TOPICS=todo-go
export PGHOST=localhost
export PGPORT=5432
export PGUSER=demo
export PGPASSWORD=superS3cret!
export APP_PORT=8085
export ATLAS_PORT=27778
export ATLAS_USER=demo
export ATLAS_PASSWORD=superS3cret!
export ATLAS_DATABASE=go-todo-cqrs
```

## Redpanda, Postgres and MongoDB

Run the following command to bring up Redpanda(Single Node), MongoDB Atlas and Postgres DB,

```shell
docker compose up -d
```

> **IMPORTANT**: Wait for all the services to be up

### Verify services

#### Redpanda

Check if Redpanda is up,

```shell
rpk topic list
```

The command should be successful but should should return a response like,

```text
NAME  PARTITIONS  REPLICAS
```

Let us create a topic that we will use in this demo,

```shell
rpk topic create todo-go
```

```text
TOPIC         STATUS
todo-go  OK
```

#### MongoDB

On a new terminal run the following command,

```shell
mongosh "mongodb://$ATLAS_USER:$ATLAS_PASSWORD@localhost:$ATLAS_PORT/?directConnection=true"
```

After successfully connecting to MongoDB run the command,

```shell
show databases
```

```text
admin   228.00 KiB
config  228.00 KiB
local   444.00 KiB
```

> **NOTE**: MongoDB will create the database `$ATLAS_DATABASE` on first insert of a document

#### PostgreSQL

The PostgreSQL can be accessed using the admin tool called `adminer` which was deployed along with the other services. To access the adminer too use the user <http://localhost:9090>.

The credentials to login:

-   System: PostgreSQL
-   Server: db
-   Username: $PGUSER
-   Password: $PGPassword
-   Database: $PGDatabase

You should see an empty database, the application will lazy create the database resources.

## Start Demo API

The CQRS demo API is simple Todo API. To start API, on a new terminal run the command,

`go run cmd/server/main.go`

The API should be available on port `8085`

### Add Todo

```shell
http -v POST :$APP_PORT <<EOF
{
    "title": "title 1",
    "description": "title one",
    "category": "learning",
    "status": false
}
EOF
```

Refresh the `adminer` window on your browser to see the inserted data.

### Update Todo

```shell
http -v PATCH :$APP_PORT/1 <<EOF
{
    "ID": 1,
    "title": "title 1",
    "description": "title one",
    "category": "learning",
    "status": true
}
EOF
```

### Delete Todo

```shell
http -v DELETE :$APP_PORT/1
```

**WIP**: Listing of data
