# 阶段 1 练习 6：Worker Pool、取消与错误收集

## 练习目标

用固定数量 Worker 控制并发，支持 Context 取消、任务结果收集和 goroutine 正常退出。

## 实际修改位置建议

- `建议在独立 package 中实现，再由学习入口调用`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：定义任务、结果和 Worker Pool 配置，拒绝非正 worker 数和无界队列。
2. TODO 2：启动固定数量 goroutine，从任务 channel 读取并执行。
3. TODO 3：在发送任务、执行任务和发送结果时响应 Context 取消。
4. TODO 4：关闭顺序由单一拥有者负责，避免重复 close、发送到已关闭 channel 和 goroutine 泄漏。
5. TODO 5：收集全部结果或首个错误，并明确取消剩余任务的策略。
6. TODO 6：使用 go test -race 覆盖成功、任务失败、提前取消和空任务。

## 验证方式

```bash
go test -race -timeout=60s ./...
go vet ./...
git diff --check
```

## 完成标准

- 改动遵循 Handler → Service → Repository 依赖方向。
- 错误和取消能够清楚传播，没有静默成功或吞错。
- 测试覆盖本练习的成功路径和主要失败路径。
- 没有修改无关模块，也没有提交真实密钥或本地配置。

## 暂不实现

- 与本练习无关的全项目重构。
- 生产环境迁移、真实数据清理或大规模性能压测。
