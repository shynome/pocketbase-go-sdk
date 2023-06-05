# Changelog

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
