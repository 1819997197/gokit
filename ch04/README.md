# ch04 API监控(Prometheus)

## 1.安装 Prometheus
```
// 
https://prometheus.io/download/ 下载Prometheus
tar xvfz prometheus-*.tar.gz
```

## 2.修改配置文件监控当前服务
```
cd prometheus-*
vim prometheus.yml

scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:50050']  # 需监视的资源
```

## 3.Run prometheus
```
./prometheus --config.file=prometheus.yml
```

## 4.Run the server
```
go build -o order
./order
```

## 5.Run the client
```
// get请求 {name}为参数
curl localhost:50050/count/{name}

// post请求
curl -X POST -d '{"s":"will"}' localhost:50050/uppercase
```

## 6.查看prometheus监控数据
```
http://localhost:9090/

// prometheus自身的监控指标
http://localhost:9090/metrics 
```

## 7.获取当前服务的监控指标
```
http://localhost:50050/metrics
```

## prometheus 采集的指标模型
```
由 metric 的名字和一系列的标签（键值对）唯一标识的,<metric name>{<label name>=<label value>, ...}
eg: api_http_requests_total{method="POST", handler="/messages"}
api_http_requests_total 则是由指标的Namespace_Subsystem_Name组成
```