[TOC]

# Meeting

> 基于 golang + gin + gorm + webrtc 实现在线会议

核心扩展

+ https://github.com/pion/webrtc
```shell
go get -u github.com/pion/webrtc/v3 
```

## 1.系统模块简介

+ [x] 会议管理
  + [x] 会议列表
  + [x] 创建会议
  + [x] 会议编辑
  + [x] 会议删除
+ [x] 用户管理
  + [x] 登录
+ [ ] WebRTC
  + [x] data channels
  + [x] 屏幕共享
  + [ ] 一对一音视频通信

## 2.项目启动
1. 在项目根目录下执行：`go mod tidy`
2. 新建MySQL数据库：meeting
3. 在`internal/server/main.go`目录下，运行main.go

## 3.接口验证

打开postman，或者是apifox等接口工具

### 3.1 /ping 

输入url：http://localhost:8080/ping ，请求方法为get

![](https://nateshao-blog.oss-cn-shenzhen.aliyuncs.com/wximage-20240829094253361.png)

接口响应pong，就代表项目启动成功

### 3.2 /user/login 用户登录

在上面步骤完成之后，在`meeting`数据库下会生成三张表

user_basic.sql

```sql
CREATE TABLE `user_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(100) NOT NULL,
  `password` varchar(36) NOT NULL,
  `sdp` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_basic_username` (`username`),
  KEY `idx_user_basic_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

room_user.sql

```sql
CREATE TABLE `room_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `rid` int NOT NULL,
  `uid` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

room_basic.sql

```sql
CREATE TABLE `room_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `identity` varchar(36) NOT NULL,
  `name` varchar(100) NOT NULL,
  `begin_at` datetime DEFAULT NULL,
  `end_at` datetime DEFAULT NULL,
  `create_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_room_basic_identity` (`identity`),
  KEY `idx_room_basic_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

**md5生成密码**

新建my_test.go测试用例

```go
package test

import (
	"meeting/internal/helper"
	"testing"
)

func TestName(t *testing.T) {
	println(helper.GetMd5("123456"))
}
```

执行测试用例，生成**密钥**：e10adc3949ba59abbe56e057f20f883e

```shell
=== RUN   TestName
e10adc3949ba59abbe56e057f20f883e
--- PASS: TestName (0.00s)
PASS

Process finished with the exit code 0
```

生成的密钥放入user_basic表的password字段中。

```sql
INSERT INTO `meeting`.`user_basic`(`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `sdp`) VALUES (1, NULL, NULL, NULL, 'root', '8017d0408f41b75489701e3fb1c3e773', NULL);
```

**postman接口验证**

1. 接口url：http://localhost:8080/user/login
2. 请求方法：POST
3. 入参：

```go
{
    "username": "root",
    "password": "1233456"
}
```

4. 接口响应：

```go
{
    "code": 200,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6InJvb3QifQ.kiFY5je3Hy_l7r41ku-eQ5PHZnsvBJ5StGP2d7bfpB0"
    }
}
```

### 3.3 /auth/meeting/list 会议列表

![](https://nateshao-blog.oss-cn-shenzhen.aliyuncs.com/wximage-20240829171446454.png)

1. 在postman接口Headers添加Authorization：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6InJvb3QifQ.kiFY5je3Hy_l7r41ku-eQ5PHZnsvBJ5StGP2d7bfpB0

2. 输入接口url：http://localhost:8080/auth/meeting/list?name=pom&page=1&size=20

3. 请求方法：GET

4. 接口响应

   ```json
   {
       "code": 200,
       "data": {
           "count": 3,
           "list": [
               {
                   "identity": "4e51c192-7b13-415f-8d88-341115c959ca",
                   "name": "pom",
                   "begin_at": "1970-01-21T07:06:06+08:00",
                   "end_at": "1970-01-21T07:06:06+08:00"
               },
               {
                   "identity": "752a4f40-451d-4b51-86ec-dc9301826085",
                   "name": "pom",
                   "begin_at": "1970-01-21T07:06:06+08:00",
                   "end_at": "1970-01-21T07:06:06+08:00"
               },
               {
                   "identity": "98a35566-7f37-4232-8cd4-d65605c3d502",
                   "name": "pom",
                   "begin_at": "1970-01-21T07:06:06+08:00",
                   "end_at": "1970-01-21T07:06:06+08:00"
               }
           ]
       }
   }
   ```

### 3.4 /auth/meeting/create 创建会议 

![](https://nateshao-blog.oss-cn-shenzhen.aliyuncs.com/wximage-20240829174946412.png)

1. 接口url：http://localhost:8080/auth/meeting/create

2. 请求方法：POST

3. 接口入参：

   ```json
   {
       "name": "qianyu",
       "begin_at": 1724766216,
       "end_at": 1724766216
   }
   ```

4. 接口响应：

   ```json
   {
       "code": 200,
       "msg": "ok"
   }
   ```

### 3.5 编辑会议

![](https://nateshao-blog.oss-cn-shenzhen.aliyuncs.com/wximage-20240829181327394.png)

1. 接口url：http://localhost:8080/auth/meeting/edit

2. 请求方法：PUT

3. 接口入参：

   ```json
   {
       "identity": "76bcb1b6-b142-40d3-b71a-4a0310c41246",
       "name": "千羽1024"
   }
   ```

4. 接口响应：

   ```json
   {
       "code": 200,
       "msg": "ok"
   }
   ```

### 3.6 删除会议

![](https://nateshao-blog.oss-cn-shenzhen.aliyuncs.com/wximage-20240829181444236.png)

1. 接口url：http://localhost:8080/auth/meeting/delete

2. 请求方法：DELETE

3. 接口入参：

   ```json
   {
       "identity": "76bcb1b6-b142-40d3-b71a-4a0310c41246",
       "name": "千羽1024"
   }
   ```

4. 接口响应：

   ```json
   {
       "code": 200,
       "msg": "ok"
   }
   ```

