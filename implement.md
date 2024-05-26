# Implement space-online

## questions

### authorize

restful authorize

- jwt ?  session ?  oath2.0 ?
- microservice > jwt + restful > microservice
- microservice > gRPC > microservice

path authorize

- implement relationship in code ? or in json file which map path to role? (x)

### commuication

- popular commuication, RESTful, websocket, socket, gRPC
- socket.io support server to multiple client, keep connected

### db

- db config in normal file or docker file to env?
- orm lib? or sql syntax
- if use orm, auto migrate? or manual migrate? advantages and disadvantages?
- if Product has id, title, desc, price, color, stock, size, accecption, it need to divide to multiple table?
- if Member has name, email, birthday, addresses, delivery, shipping, it need to divide to multiple table? 711 table, family table, home table
- table primary key use auto increase or random?
- how to design command flow, event loop? or directly execute
- how to design cache, oo expire time? priortiy? or working pool? every table has cache?
- column value case insentive

### todo

- logger
- event loop
- command pool
- cobra for command
- viper for env
- websocket
- gRPC
