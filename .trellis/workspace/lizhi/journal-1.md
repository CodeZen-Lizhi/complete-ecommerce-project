# Journal - lizhi (Part 1)

> AI development session journal
> Started: 2026-07-14

---


## Session 1: 完成阶段 2 Eino 练习与项目规范初始化

**Date**: 2026-07-14
**Task**: 完成阶段 2 Eino 练习与项目规范初始化
**Branch**: `master`

### Summary

完成基础 HTTP 与 Eino Generate 学习练习，补充中文 TODO 协作规范，初始化并完善 Trellis 后端开发规范，清理敏感 API Key 后通过 Go 测试、vet 与提交检查。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `47494d3` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 2: 完成 System Prompt 角色对照练习

**Date**: 2026-07-15
**Task**: 完成 System Prompt 角色对照练习
**Branch**: `master`

### Summary

完成基于 Eino AgenticModel 和 Responses API 的 System Prompt 角色对照练习，验证相同 User Prompt 在不同角色约束下的输出差异，并通过目标包测试和 go vet。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `f2bcf3f` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 3: 完成程序内多轮对话练习

**Date**: 2026-07-15
**Task**: 完成程序内多轮对话练习
**Branch**: `master`

### Summary

完成 Eino ResponsesModel 程序内两轮对话练习，正确保存完整 Assistant 消息并重组历史，为每轮调用设置独立超时，补齐错误处理和输出标签，同时清理源码中的真实格式 API Key。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `0033d2c` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 4: 完成 Redis 会话历史练习

**Date**: 2026-07-15
**Task**: 完成 Redis 会话历史练习
**Branch**: `master`

### Summary

完成基于 Eino ResponsesModel 和 Redis List 的 session 会话历史练习，按轮次区分完整模型上下文与待持久化消息，补齐 session 校验、JSON 编解码、错误上下文、中文函数注释规范及敏感 Key 清理。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `853f3d3` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 5: 完成全阶段编号化 AI 练习

**Date**: 2026-07-15
**Task**: 完成全阶段编号化 AI 练习
**Branch**: `master`

### Summary

生成 7 阶段 60 个编号练习，迁移阶段 2 既有内容，补齐目录契约测试、安全占位骨架、Compose 契约及项目规范，并完成全量 Go 验证与 review。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `44a427b` | (see git log) |
| `f32569d` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 6: 归档会话历史截断练习

**Date**: 2026-07-15
**Task**: 归档会话历史截断练习
**Branch**: `master`

### Summary

确认阶段 2 会话历史截断与用户隔离练习验收完成，归档遗漏任务；代码已包含在编号化 AI 练习提交中。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `44a427b` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 7: 完成阶段 2 会话、模板与流式练习

**Date**: 2026-07-18
**Task**: 完成阶段 2 会话、模板与流式练习
**Branch**: `master`

### Summary

完成阶段 2 练习 6-8：补齐会话历史重组与隔离、ChatTemplate 变量和顺序校验、ChatModel Stream 的 EOF/取消/空流处理；改用环境变量读取 API Key，并补充单元测试。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `9b5a932` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete


## Session 8: 阶段二 AI 练习提交收口

**Date**: 2026-07-18
**Task**: 阶段二 AI 练习提交收口
**Branch**: `master`

### Summary

补充会话角色不变量、失败路径测试与 API Key 文档，并完成双轮 Go 质量门禁。

### Main Changes

- Detailed change bullets were not supplied; see the summary above.

### Git Commits

| Hash | Message |
|------|---------|
| `8ffa28e` | (see git log) |

### Testing

- Validation was not recorded for this session.

### Status

[OK] **Completed**

### Next Steps

- None - task complete
