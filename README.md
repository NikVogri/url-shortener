# Go: Url Shortener

Fun little project to teach myself the Go language.

## Endpoints

`POST /add`
Add a new URL and return an ID pointing towards it.

Request body:

```
{
    "url": string,
    "duration": int // time in MS
}
```

Response body:

```

{
    "url": string // redirect ID
}
```

`GET /:id`
If valid ID, redirects to the URL and increments URL record click count.

## Usage
Simply start the database & API using docker-compose, run:
```
$ git clone git@github.com:NikVogri/url-shortener.git .
$ docker-compose up -d
```

To stop the execution, run:
```
$ docker-compose down
```
## Configuration

To configure API port or Postgres credentials, visit the `docker-compose.yml` file. 

Make sure to also update the `DB_CONN_STR` env. variable when changing PG credentials and `PORT` constant when changing the port


## Debugging

Use the Docker cli or Docker Desktop to view the the logs.

Delete `db-data` directory to reset the database or connect to the database via cli or PgAdmin and manage it from there.