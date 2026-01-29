# Swagger Documentation Guide

This project uses [gin-swagger](https://github.com/swaggo/gin-swagger) and [swag](https://github.com/swaggo/swag) to generate and serve API documentation.

## Accessing the Documentation
When the server is running, you can access the interactive API documentation at:

**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

## How to Regenerate Documentation
If you modify the API controllers or models, you need to regenerate the Swagger documentation.

Prerequisites:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Command to regenerate:
```bash
swag init -d cmd/server,internal/controller,internal/model -g main.go --parseDependency --parseInternal
```

## Annotations
The documentation is generated from comments in the code.
- `main.go`: General API information (Title, Version, Host).
- Controllers (`internal/controller/*.go`): Endpoint definitions (@Summary, @Param, @Success).

Example Annotation:
```go
// Login godoc
// @Summary Login user
// @Description Login user and generate JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) { ... }
```
