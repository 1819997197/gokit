# ch02 日志中间件

## 1.Run the server
```
go build -o order
./order
```

## 2.Run the client
```
// get请求 {name}为参数
curl localhost:50050/count/{name}

// post请求
curl -X POST -d '{"s":"will"}' localhost:50050/uppercase
```