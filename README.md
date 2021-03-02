## Demo  项目

### 接口

1、 添加用户

```shell
curl --location --request POST 'localhost:8080/v1/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "nickName": "pdf13",
    "role": 1
}'
```

2、 删除用户

```shell
curl --location --request DELETE 'localhost:8080/v1/user/:userId'
```

3、 获取用户详情

```shell
curl --location --request GET 'localhost:8080/v1/user/:userId'
```

4、 获取用户列表

```shell
curl --location --request GET 'localhost:8080/v1/user?skip=10&limit=10011'
```

### Redis 存储说明

每个用户都以KV存储，用户信息使用JSON序列化为字符串，所以获取用户信息直接可以拿用户Id去获取。

同时，由于有列表需求，所以把用户Id放一个 ZSET 中，即方便 RANGE ，也方便删除


### 单元测试说明

启动单测应该先启动一个sidecar redis，并且设置环境变量 `GO_TEST_REDIS_URL` 为 sidecar redis的链接地址。

然后运行：

```shell
go test cmd/main_test.go
```
