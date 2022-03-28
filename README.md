# GO POSTGRES EXAMPLE

## PostgeSQL cli

```
docker exec -it postgres sh
psql -U user testdb
\l
\d
```

## Seed fake data

```
docker exec -it go-redis-example-go-1 go run cmd/seeder/main.go
```
