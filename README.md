# README

<h2>Start PostgreSQLon Docker 🐋</h2>

```bash
docker run --name e-commerce -e POSTGRES_USER=banana -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:alpine
```

<h2>Execute a container and CREATE a new database</h2>

```bash
docker exec -it e-commerce bash
psql -U banana
CREATE DATABASE ecommerce;
\l
```

<h2>Migrate command</h2>

install migrate cli golang

```bash
brew install golang-migrate

```

```bash
# run migration file
make up.migrate

# down migration file
make down.migrate
```

<h2>Run Project</h2>

```bash
make dev
```
