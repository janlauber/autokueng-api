# autokueng-api
The autokueng.ch website api

## PostgreSQL
Run Test PostgreSQL docker instance:
```bash
docker run -it -p 5432:5432 -d -e POSTGRES_HOST_AUTH_METHOD=trust postgres
```
### Env Variables
Set the following env variables:
```bash
export DB_USERNAME=postgres
```
```bash
export DB_PASSWORD=
```
```bash
export DB_NAME=postgres
```
```bash
export DB_HOST=localhost
```

## API Reference
The Autokueng API is organized around REST. 
Our API has predictable resource-oriented URLs, accepts form-encoded request bodies, returns JSON-encoded responses, ans uses standard HTTP response codes, authentication, and verbs.

### Authentication
The authentication is done via JWT.

| URL | Method | Data | Description |
| ---- | ------ | ----------- | ----------- |
| `/api/v1/auth/login` | POST | `username`(string) `password`(string) | Login to the API |
| `/api/v1/auth/logout` | POST | `jwt`(JWT) | Logout from the API |
| `/api/v1/auth/register` | POST | `username`(string) `password`(string) | Register a new user (disabled by default) |