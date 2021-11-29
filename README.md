# autokueng-api
The autokueng.ch website api

## PostgreSQL
```bash
docker run -it -p 5432:5432 -d -e POSTGRES_HOST_AUTH_METHOD=trust postgres
export DB_USER=postgres
export DB_PASS=
export DB_NAME=postgres
export DB_HOST=localhost
```