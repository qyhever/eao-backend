---
applyTo: "internal/**/*.go"
description: "Use when editing Go handlers, services, repositories, or router wiring in eao-backend, especially when adding or changing API endpoints."
---

# eao-backend API Endpoint Rules

- 保持现有分层：controller 负责绑定参数和返回响应，service 负责业务逻辑，repository 定义接口，persistence 提供实现。
- 新增接口时，优先沿着 router -> controller -> service -> repository/persistence 全链路补齐，而不是把逻辑直接塞进 controller。
- 统一复用 internal/controller/response.go 中的 ResponseSuccess、ResponseFailed、ResponseFailedWithMsg。
- 参数校验优先放在 controller，使用现有的绑定方式和错误返回风格。
- 修改路由时同步检查 internal/api/router.go 的依赖构造方式，保持 NewXxxController(NewXxxService(NewXxxRepository())) 这类装配模式一致。
- 如果新增或修改 HTTP 接口，顺手更新 rest/index.http 中的示例请求，方便后续代理和开发者验证。
- 如果改动会影响本地示例数据或静态响应，检查 public/meta.json 和 public/post.json 是否也需要同步。

参考入口：
- AGENTS.md
- internal/api/router.go
- internal/controller/response.go
- rest/index.http