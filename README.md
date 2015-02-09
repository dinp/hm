HM
==========

HM(Health Monitor)用于APP rs (docker container) 的7层健康检查。
如要开启APP rs的7层健康检查，在DashBoard的APP页面的health栏填写健康检查的URL (如/health)。

HM每隔check_interval (单位s)
- 从Server模块的DB中查询需要健康检查的APP列表 (即health字段不为空的APP)，
- 从Server模块的HTTP接口中查询APP及相应的rs列表，
然后curl相应的rs，如果在response_timeout (单位s) 内返回内容包括health_sign，则视该rs为健康；否则将杀掉该rs (docker container)，由Server模块进行调度创建新的rs。

## 配置项说明

- **debug**: true/false 只影响打印的log
- **check_interval**: 健康检查的周期，单位s
- **dockerPort**: Docker Daemon的侦听端口
- **response_timeout**: APP health接口的响应超时时间，单位s
- **health_sign**: APP health接口的返回内容，如'''ok'''
- **server_http_api**: DINP Server模块的http api接口
- **db**: DINP Server模块数据库的地址，以及超时时间

install

```
mkdir -p $GOPATH/src/github.com/dinp
cd $GOPATH/src/github.com/dinp; git clone https://github.com/dinp/hm.git
cd hm
go get ./...

# check cfg.json, depend docker daemon and server
hm -c cfg.json
```
