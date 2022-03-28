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

## Redis cli

```
docker exec -it redis redis-cli
127.0.0.1:6379> set key value
OK
127.0.0.1:6379> get key
"value"
```
