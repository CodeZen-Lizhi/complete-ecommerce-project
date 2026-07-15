# 阶段 7 练习 6：Docker Compose 本地栈

## 练习目标

使用 Docker Compose 启动应用、MySQL、Redis 和向量库。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义不包含真实密钥的环境变量契约。
2. TODO 2：为 MySQL、Redis 和向量库配置持久卷与健康检查。
3. TODO 3：为应用配置依赖健康条件和本地端口。
4. TODO 4：验证启动、健康、重启和数据持久化。
5. TODO 5：提供停止与清理命令，默认不删除持久数据。

## 开始练习

```bash
export MYSQL_ROOT_PASSWORD='local-exercise-only'
docker compose -f examples/ai/phase7/06_docker_compose_stack/compose.yaml --profile exercise config
docker compose -f examples/ai/phase7/06_docker_compose_stack/compose.yaml --profile exercise up -d
```

基础设施 Profile 只启动 MySQL、Redis 和 Qdrant。完成应用 Dockerfile 后，再设置 `AI_APP_IMAGE` 并追加 `--profile app` 启动应用；源码中的模型 Key 仍保持占位符。

## 验证方式

```bash
gofmt -w examples/ai/phase7/06_docker_compose_stack/*.go
go test -timeout=60s ./examples/ai/phase7/06_docker_compose_stack
go vet ./examples/ai/phase7/06_docker_compose_stack
MYSQL_ROOT_PASSWORD=local-exercise-only docker compose \
  -f examples/ai/phase7/06_docker_compose_stack/compose.yaml \
  --profile exercise config --quiet
```

## 完成标准

- `docker compose config` 通过，模板中没有真实凭证。
- MySQL、Redis 和 Qdrant 使用持久卷并能通过健康检查。
- `app` 只从环境变量读取依赖地址和密钥，并等待依赖健康。
- `docker compose down` 默认保留数据；删除卷需要显式 `down -v`。

## 暂不实现

- Kubernetes、生产镜像发布和真实云存储。
