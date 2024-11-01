# Todo List API
In this project, create a blogging app with authentication and authorization  

## roadmap.sh
This project is based on https://roadmap.sh/projects/blogging-platform-api  

## install echo
```go get github.com/labstack/echo/v4```

## install postgres
```go get github.com/jackc/pgx/v5```
```go get github.com/jackc/pgx/v5/pgxpool```

## install validator
```go get github.com/go-playground/validator/v10```

## install testify
```go get github.com/stretchr/testify```

## test
```go test -v tests/unit_tests/features/users/login/services/login_service_test.go```

## add evironment variables
```
export ECHO_HOST=:8080
export POSTGRES_HOST=localhost:5432
export POSTGRES_USERNAME=postgres
export POSTGRES_PASSWORD=12345
export POSTGRES_DATABASE=blog
export POSTGRES_MAX_CONNECTION=10
export POSTGRES_MAX_IDLETIME=10
export POSTGRES_MAX_LIFETIME=10
export COOKIE_SECURE=false
```

## run project
To run this project, just download the project, go to downloaded project and run it by typing ```go run main.go``` and press enter
access it through browser with ```http://localhost:8080/posts```

```happy koding and thank you :D```