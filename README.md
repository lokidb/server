# lokidb server
LokiDB server

---

### Table of Contents
- [lokidb server](#lokidb-server)
    - [Table of Contents](#table-of-contents)
    - [Featurs](#featurs)
    - [Installation](#installation)
      - [Docker](#docker)
    - [Usage](#usage)
        - [server](#server)
        - [cli](#cli)
    - [API](#api)


### Featurs
- gRPC support
- REST support
- slim docker image


### Installation
#### Docker
```
docker run -d -p 50051:50051 -p 8080:8080 --name lokidb yoyocode/lokidb
```

```
docker run -it yoyocode/lokidb /lokidb-cli -addr <grpc_host>:<grpc_port>
```

### Usage
##### server
```shell
$ lokidb --help
Usage of lokidb:
  -cache_size int
        Max number of keys to save in the RAM for fast reads (default 50000)
  -files_count int
        number of files to save the keys on (more files means better performance, max 1000) (default 10)
  -grpc_host string
        gRPC server host (default "127.0.0.1")
  -grpc_port int
        gRPC server port (default 50051)
  -rest_host string
        REST server host (default "127.0.0.1")
  -rest_port int
        REST server port (default 8080)
  -run_grpc
        Serve gRPC API (default true)
  -run_rest
        Serve REST API (default true)
  -storage_dir string
        The path where the data files will be created (default "./")
```

```shell
$ lokidb -storage_dir data/
2022/08/25 23:52:13 gRPC server listening at 127.0.0.1:50051
2022/08/25 23:52:13 REST server listening at 127.0.0.1:8080
```

##### cli
```shell
$ lokidb-cil --help
Usage of /tmp/go-build1997907887/b001/exe/main:
  -addr string
        the address to connect to (default "localhost:50051")
```

```shell
$ lokidb-cli
>>> help
commands:
    get 
    set 
    del 
    keys 
    flush 
    bye 
    help 

>>> keys
0) a
1) b
2) mom

>>> get asdlg


>>> get a
vavava

>>> exit
```

### API  
The system support REST API as well as a gRPC server.  
REST: [OpenAPI schema](/rest/spec.yaml).  
gRPC: [ProtoBuf schema](/grpc/spec.proto).  
The repo also contains an [insomnia workspace](https://insomnia.rest/) (a better looking postman) with all of the requests and enviroments.  