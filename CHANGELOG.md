# Changelog

## [2.0.0] - Unrelease

重构接口, 以便和 js-sdk 的接口保持一致.

暂时没空, 先用着 1.0 吧

## [1.3.5] - 2026-06-02

- 修复: 原来选项是直接附加到 topic 上的

## [1.3.4] - 2026-06-02

- 修复: 移除 debug 标识

## [1.3.3] - 2026-06-02

- 修复: `Subscription.Record` 对应的 json 字段是 `record` 而不是 `data`

## [1.3.2] - 2026-06-02

- 修复: 使用 `c.SubscribeToAll` 会导致 `PB_CONNECT` 事件也被发送回调, 转而使用 `c.SubscribeEvent(topic)`

## [1.3.1] - 2026-06-02

- 修复: pocketbase 使用 event type 来区分订阅 topic, 不能使用 `c.SubscribeMessages` 而是要用 `c.SubscribeToAll`

## [1.3.0] - 2026-06-02

- 添加: 实现 Subscribe

## [1.2.0] - 2026-06-02

- 添加 Send 方法以便发送 SSE 认证信息

## [1.1.0] - 2023-07-16

- 兼容 pb v0.23

## [1.0.0] - 2023-07-16

添加了测试, 自信地发布 1.0.0

## [0.0.5] - 2023-06-05

### Fix

- 只重试权限错误的请求, 即 403 状态码

## [0.0.4] - 2023-06-05

### Fix

- 修复 error message data 类型错误

## [0.0.3] - 2023-06-04

### Add

- 添加 `client.SetAuthStore`, 手动设置 AuthStore. 这点和官方 sdk 不一样

## [0.0.2] - 2023-06-03

### Add

- 添加 `crud.Service.FirstListItem`, 便于只获取一个元素
