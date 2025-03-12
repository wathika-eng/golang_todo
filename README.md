# Crud TODO API

![Go](https://img.shields.io/badge/Go-1.24-blue)

## Description

Technologies used:

[] - Golang 1.24

[] - Gin // Router <https://gin-website.vercel.app/>

[] - Bun // ORM <https://bun.uptrace.dev/guide/>

[] - Postgres // Database

[] - Redis // Cache and Blacklist Tokens <https://redis.io/docs/latest/develop/clients/go/>
<!-- [] - Docker
[] - Docker Compose -->

# To get started

- Clone the repository

```bash
git clone https://github.com/wathika-eng/golang_todo --depth 1 && cd golang_todo 
```

- Install dependencies

```bash
go mod download
```

- Create a .env file in the root directory and add the following environment variables

```bash
cp .env.example .env
```

Get a free postgresql database from <https://neon.tech/> and add the database url to the .env file (CONNECTION_STRING=postgres://username:password@host:port/dbname)

- Run the application

```bash
go run main.go
```

Or with docker:

```bash
COMPOSE_BAKE=true docker compose up -d --build
```

Access the API via: <http://localhost:8000/api/users/test>

Or <http://localhost/api/users/test> if you are using docker

Live URL => <http://157.230.108.128:8000/api/users/test>

Frontend =>

Routes:

```json
[GIN-debug] GET    /api/users/test           
[GIN-debug] POST   /api/users/signup         
[GIN-debug] POST   /api/users/login          
[GIN-debug] POST   /api/users/refresh        

# protected routes
[GIN-debug] GET    /api/notes/test           
[GIN-debug] POST   /api/notes/create        
[GIN-debug] GET    /api/notes/               
[GIN-debug] GET    /api/notes/:id           
[GIN-debug] PATCH  /api/notes/:id           
[GIN-debug] DELETE /api/notes/:id           
[GIN-debug] POST   /api/notes/logout         
[GIN-debug] GET    /api/profile/           
```

## Resources

<https://12factor.net/>

<https://threedots.tech/post/repository-pattern-in-go/>

<https://jub0bs.com/posts/>
