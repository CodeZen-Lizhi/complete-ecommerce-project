# 阶段 1 练习 7：后端分层与失败路径测试

## 练习目标

为 Service、Handler、Repository 和工具边界建立最小但有效的测试保护。

## 实际修改位置建议

- `service/`
- `handler/`
- `repository/`
- `util/`
- `middleware/`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：选择一个完整业务路径，列出成功、空值、非法输入、依赖失败和权限失败场景。
2. TODO 2：Service 使用 Fake Repository，验证业务规则和事务回滚。
3. TODO 3：Handler 使用 httptest 和 Gin，验证状态码、响应结构和 Abort 行为。
4. TODO 4：Repository 使用隔离数据库或受控测试环境验证 NotFound、软删除和零值更新。
5. TODO 5：工具函数使用表驱动测试，保证删除实现后测试会失败。
6. TODO 6：运行全仓测试、vet；涉及共享状态时追加 race。

## 验证方式

```bash
go test -timeout=60s ./...
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
