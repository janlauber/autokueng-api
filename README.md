# Autoküng API
The autokueng.ch website api
- [frontend](https://github.com/janlauber/autokueng-frontend)
- [data api](https://github.com/janlauber/autokueng-data)

## PostgreSQL
Run Test PostgreSQL docker instance:
```bash
docker run -it -p 5432:5432 -d -e POSTGRES_HOST_AUTH_METHOD=trust postgres
```
### Env Variables
Set the following env variables:  
*database*
```bash
export DB_USERNAME=postgres
export DB_PASSWORD=
export DB_NAME=postgres
export DB_HOST=localhost
```
*jwt secret*
```bash
# must be the same as in autokueng data api
export JWT_SECRET_KEY=<secret>
```
*user recovery/administration mode*
if you want to reset the password of a user or register a new one, set the following env variable:
```bash
export USER_ADMIN=enabled
```

## API Reference
The Autokueng API is organized around REST. 
Our API has predictable resource-oriented URLs, accepts form-encoded request bodies, returns JSON-encoded responses, ans uses standard HTTP response codes, authentication, and verbs.

### Authentication
The authentication is done via JWT.

| URL | Method | Data | Description |
| ---- | ------ | ----------- | ----------- |
| `/api/v1/login` | POST | `username`(string) `password`(string) | Login to the API |
| `/api/v1/logout` | POST | `jwt`(JWT) | Logout from the API |
| `/api/v1/user` | GET | `username` `Id` | Checks if user is logged in and sends back username and userId, Is used to check if user is authenticated |

### User Administration
The user administration is done via the `USER_ADMIN` env variable.

| URL | Method | Data | Description |
| ---- | ------ | ----------- | ----------- |
| `/api/v1/admin/login` | POST | `username`(string) `password`(string) | Login to the API |
| `/api/v1/admin/logout` | POST | `jwt`(JWT) | Logout from the API |
| `/api/v1/admin/register` | POST | `username`(string) `password`(string) | Register a new user (disabled by default) |