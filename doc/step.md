## [Frontend]
https://space-online-shop.web.app/welcome

## [Packages]
1. code with golang
1. config with viper
1. log with zap
1. webserver with gin
1. db with mysql
1. webpack with dockerfile
1. deploy with cloud?

## [Implement]
1. member system
    1. register
    1. login
    1. logout

## [Docker install & uninstall mysql]

Create networks, volumes. Start container

```sql
docker-compose.exe -f cfg\mysql\mysql.yml up -d
```

Remove container, networks, volumes

```sql
docker-compose.exe -f cfg\mysql\mysql.yml down --volumes
```
