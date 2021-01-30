# profilesvc 服务注册与发现

## 1.Run the consul
```
./consul agent -dev
```

## 2.查看consul中已注册的某个服务
```
curl -s 127.0.0.1:8500/v1/catalog/service/{service_name}
```

## 3.Run the  service
```bash
// 服务实例一
$ go run main.go -service.port 8080
// 服务实例二
$ go run main.go -service.port 8081
```

## 4.Create a Profile:

```bash
$ curl -d '{"id":"1234","Name":"Go Kit"}' -H "Content-Type: application/json" -X POST http://localhost:8080/profiles/
{}
```

## 5.Get the profile you just created

```bash
$ curl localhost:8080/profiles/1234
{"profile":{"id":"1234","name":"Go Kit"}}
```

## 6.client目录代码实现的客户端(负载均衡)
```
go run client/client.go
```