# pickcampus-backend

PickCampus（去哪儿上大学）后端服务。T1 阶段提供：邮箱注册/登录、JWT + Redis 会话、拿当前用户。

## 技术栈

- Gin（HTTP）
- GORM + MySQL（数据）
- go-redis/v8（会话）
- golang-jwt/v5（JWT）
- bcrypt（密码哈希）
- cleanenv（配置）+ Cobra（CLI）

## 分层

```
cmd/app/main.go            入口
internal/
  app/run.go               启动引导（读配置→初始化 DB/Redis→迁移→建 Gin→起服务）
  config/                  配置结构 + 全局 G
  bootstrap/               MySQL / Redis 客户端初始化（官方库，无重封装）
  common/                  错误码、JWT、Redis key、鉴权中间件、常量
  handler/                 HTTP 层（响应封装、路由、user handler）
  logic/                   业务逻辑（bcrypt、token、Redis 会话编排）
  types/                   请求/响应 DTO
models/
  base.go / user.go        GORM 模型（tbl_ 前缀、Unix 秒时间戳）
  init.go                  AllTables 自动迁移清单
  repo/                    数据访问接口
  factory/                 GORM 实现
```

## 配置

```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml，填写 MySQL / Redis / jwt_token_secret
```

`config.yaml` 含密钥，已在 `.gitignore` 中，不进 Git。所有字段可用同名环境变量覆盖。

## 起服务

前置：本地或远程可用的 MySQL 与 Redis（DB 若不存在会自动创建）。

```bash
go mod tidy
go run cmd/app/main.go -c configs/config.yaml
# 默认监听 0.0.0.0:8080
```

## 接口

响应统一格式：`{ "code": 0, "message": "...", "data": {...} }`（code 非 0 为错误）。
鉴权 Header：`Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 | 鉴权 |
|---|---|---|---|
| GET | `/ping` | 健康检查 | 否 |
| POST | `/api/v1/register` | 邮箱注册 | 否 |
| POST | `/api/v1/login` | 登录（返回 JWT） | 否 |
| GET | `/api/v1/user` | 拿当前用户 | 是 |
| POST | `/api/v1/logout` | 登出（吊销会话） | 是 |

### curl 示例

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"S3cur3Pass","username":"alice"}'

# 登录（拿 token）
curl -X POST http://localhost:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"S3cur3Pass"}'

# 拿当前用户
curl http://localhost:8080/api/v1/user \
  -H 'Authorization: Bearer <token>'

# 登出（登出后原 token 调 /user 会 401）
curl -X POST http://localhost:8080/api/v1/logout \
  -H 'Authorization: Bearer <token>'
```

## 测试

```bash
go test ./...    # 单测不依赖外部 DB/Redis
go vet ./...
```

## 会话与安全说明

- 登录时把 token 写入 Redis `token:{userID}`；每次请求鉴权都会从 Redis 读回并逐字比对（fail-secure）：不一致 / 不存在 / Redis 不可用一律 401。因此服务端可主动吊销（登出即删 key）。
- 同一用户只存一个 token（单设备语义），重新登录覆盖旧 token。
- 登录失败统一返回「邮箱或密码错误」，不区分用户不存在与密码错误，防账号枚举。
- 密码用 bcrypt 哈希存储，`password_hash` 字段 `json:"-"` 绝不外泄。
