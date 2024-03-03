## A simple 'software subscription service' - web application in Golang

### Modules in use

- [http session management](https://pkg.go.dev/github.com/alexedwards/scs)
- [http session redis store](https://pkg.go.dev/github.com/alexedwards/scs/redisstore)
- [postgres driver](https://pkg.go.dev/github.com/jackc/pgx)
- [web framework Go-chi](https://pkg.go.dev/github.com/go-chi/chi/v5)

### Requirements for env-file

- `.env` should contain - put correct values instead of **_'!!!'_**:

```
POSTGRES_USER=!!!
POSTGRES_PASSWORD=!!!
POSTGRES_DB=!!!
BINARY_NAME=myapp
DSN="host=localhost port=5432 user=!!! password=!!! dbname=!!! sslmode=disable timezone=UTC connect_timeout=5"
REDIS="127.0.0.1:6379"
```
