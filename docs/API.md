# Mobile API Doc

Last updated: 2026-03-25

## Status

- 当前代码里已经暴露的接口，本文档都已覆盖。
- 这是“当前可对接版本”文档，不包含设计稿里尚未实现的未来接口。
- 机器可读版本见 `docs/openapi.yaml`

## Base Info

- Base URL:
  - dev: `http://localhost:8080`
  - production: 由部署环境域名决定
- API prefix: `/api`
- Content-Type: `application/json`
- Auth scheme: `Authorization: Bearer <access_token>`
- Metrics endpoint: `/metrics`

## Unified Response

所有接口统一返回：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

说明：

- `code = 0` 表示成功
- `code != 0` 表示业务或参数错误
- `msg` 为错误描述
- `data` 成功时为业务数据，失败时为 `null`

## Auth

### POST `/api/auth/register`

说明：用户注册

请求体：

```json
{
  "email": "user@example.com",
  "password": "123456",
  "nickname": "dawn"
}
```

字段：

- `email`: 必填，邮箱格式
- `password`: 必填，最少 6 位
- `nickname`: 选填

成功响应 `data`：

```json
{
  "access_token": "string",
  "refresh_token": "string",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "dawn",
    "avatar_url": "",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

### POST `/api/auth/login`

说明：邮箱密码登录

请求体：

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

成功响应：同注册接口

### POST `/api/auth/refresh`

说明：刷新 token

请求体：

```json
{
  "refresh_token": "string"
}
```

成功响应：

```json
{
  "access_token": "string",
  "refresh_token": "string",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

## User

### GET `/api/user/profile`

说明：获取当前登录用户资料

需要登录：是

成功响应 `data`：

```json
{
  "id": 1,
  "email": "user@example.com",
  "nickname": "dawn",
  "avatar_url": "https://cdn.example.com/avatar.png",
  "created_at": "2026-03-23T10:00:00Z",
  "updated_at": "2026-03-23T10:00:00Z"
}
```

### PUT `/api/user/profile`

说明：更新当前登录用户资料

需要登录：是

请求体：

```json
{
  "nickname": "new-name",
  "avatar_url": "https://cdn.example.com/avatar.png"
}
```

成功响应 `data`：

```json
{
  "updated": true
}
```

## Address

### GET `/api/user/addresses`

说明：获取当前用户地址列表

需要登录：是

成功响应 `data`：

```json
[
  {
    "id": 1,
    "user_id": 1,
    "receiver_name": "张三",
    "phone": "13800138000",
    "province": "上海市",
    "city": "上海市",
    "district": "浦东新区",
    "detail": "世纪大道 100 号",
    "is_default": true,
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
]
```

### POST `/api/user/address`

说明：新增地址

需要登录：是

请求体：

```json
{
  "receiver_name": "张三",
  "phone": "13800138000",
  "province": "上海市",
  "city": "上海市",
  "district": "浦东新区",
  "detail": "世纪大道 100 号",
  "is_default": true
}
```

成功响应 `data`：

```json
{
  "created": true
}
```

### PUT `/api/user/address/:id`

说明：更新地址

需要登录：是

路径参数：

- `id`: 地址 ID

请求体：同新增地址

成功响应 `data`：

```json
{
  "updated": true
}
```

### DELETE `/api/user/address?id={id}`

说明：删除地址

需要登录：是

注意：当前实现使用 query 参数传 `id`，不是 path 参数。

成功响应 `data`：

```json
{
  "deleted": true
}
```

## Product

### GET `/api/products/search`

说明：商品搜索

查询参数：

- `keyword`: 选填，商品名模糊搜索
- `sort`: 选填
  - `price_asc`
  - `price_desc`
  - `sales`
  - 不传则默认按 `id desc`

成功响应 `data`：

```json
[
  {
    "id": 1001,
    "category_id": 10,
    "name": "无线耳机",
    "description": "降噪蓝牙耳机",
    "price": 39900,
    "stock": 50,
    "cover_url": "https://cdn.example.com/p1.png",
    "sales": 12,
    "tags": null
  }
]
```

说明：

- `price` 当前为整数，建议按“分”处理

### GET `/api/products/:id`

说明：商品详情

路径参数：

- `id`: 商品 ID

成功响应 `data`：单个商品对象，字段同搜索接口

## Cart

### GET `/api/cart`

说明：获取购物车列表

需要登录：是

成功响应 `data`：

```json
[
  {
    "id": 1,
    "user_id": 1,
    "product_id": 1001,
    "quantity": 2,
    "checked": true,
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
]
```

### POST `/api/cart`

说明：加入购物车

需要登录：是

请求体：

```json
{
  "product_id": 1001,
  "quantity": 2,
  "checked": true
}
```

成功响应 `data`：

```json
{
  "created": true
}
```

### PUT `/api/cart/:id`

说明：更新购物车项数量和勾选状态

需要登录：是

路径参数：

- `id`: 购物车项 ID

请求体：

```json
{
  "quantity": 3,
  "checked": true
}
```

成功响应 `data`：

```json
{
  "updated": true
}
```

### DELETE `/api/cart?id={id}`

说明：删除购物车项

需要登录：是

注意：当前实现使用 query 参数传 `id`

成功响应 `data`：

```json
{
  "deleted": true
}
```

## Order

### GET `/api/orders`

说明：获取我的订单列表

需要登录：是

成功响应 `data`：

```json
[
  {
    "id": 1,
    "user_id": 1,
    "address_id": 10,
    "order_no": "M202603231010101",
    "status": "pending",
    "total_amount": 79800,
    "idempotency_key": "checkout-20260323-1",
    "created_at": "2026-03-23T10:10:10Z",
    "updated_at": "2026-03-23T10:10:10Z",
    "items": null
  }
]
```

### GET `/api/orders/:id`

说明：获取订单详情

需要登录：是

路径参数：

- `id`: 订单 ID

成功响应 `data`：

```json
{
  "id": 1,
  "user_id": 1,
  "address_id": 10,
  "order_no": "M202603231010101",
  "status": "pending",
  "total_amount": 79800,
  "idempotency_key": "checkout-20260323-1",
  "created_at": "2026-03-23T10:10:10Z",
  "updated_at": "2026-03-23T10:10:10Z",
  "items": [
    {
      "id": 11,
      "order_id": 1,
      "product_id": 1001,
      "product_name": "无线耳机",
      "product_price": 39900,
      "quantity": 2
    }
  ]
}
```

### POST `/api/orders`

说明：创建订单

需要登录：是

请求体：

```json
{
  "address_id": 10,
  "idempotency_key": "checkout-20260323-1"
}
```

说明：

- 订单来源于“当前购物车里 `checked=true` 的商品”
- `idempotency_key` 必填，移动端每次点击提交订单时应生成唯一值

成功响应 `data`：

```json
{
  "id": 1,
  "user_id": 1,
  "address_id": 10,
  "order_no": "M202603231010101",
  "status": "pending",
  "total_amount": 79800,
  "idempotency_key": "checkout-20260323-1",
  "items": [
    {
      "id": 11,
      "order_id": 1,
      "product_id": 1001,
      "product_name": "无线耳机",
      "product_price": 39900,
      "quantity": 2
    }
  ]
}
```

## Home

### GET `/api/home/banners`

说明：首页 banner 列表

成功响应 `data`：

```json
[
  {
    "id": 1,
    "title": "春季活动",
    "image_url": "https://cdn.example.com/banner.png",
    "link_url": "/pages/activity?id=1"
  }
]
```

### GET `/api/home/categories`

说明：首页分类列表

成功响应 `data`：

```json
[
  {
    "id": 10,
    "name": "数码",
    "icon_url": "https://cdn.example.com/cat.png",
    "sort_order": 1
  }
]
```

### GET `/api/home/recommend`

说明：首页推荐商品

成功响应 `data`：商品数组，字段同商品搜索接口

### GET `/api/home/recommend/personalized`

说明：个性化推荐

需要登录：是

当前策略：

- 优先使用用户历史订单里的偏好类目做推荐
- 如果用户没有历史订单，则回退到通用热门推荐

成功响应 `data`：商品数组，字段同商品搜索接口

## Video

### GET `/api/videos/recommend`

说明：短视频推荐列表

成功响应 `data`：

```json
[
  {
    "id": 1,
    "title": "耳机开箱",
    "cover_url": "https://cdn.example.com/video-cover.png",
    "playback_url": "https://cdn.example.com/video.mp4",
    "product_id": 1001
  }
]

## Seckill

### POST `/api/seckill/orders`

说明：秒杀下单入口

需要登录：是

请求体：

```json
{
  "product_id": 1001,
  "quantity": 1,
  "address_id": 10,
  "idempotency_key": "seckill-20260325-1"
}
```

说明：

- `product_id`: 秒杀商品 ID
- `quantity`: 数量，当前建议传 `1`
- `address_id`: 收货地址 ID
- `idempotency_key`: 必填，防重复提交

成功响应 `data`：

```json
{
  "queued": false,
  "message_id": "SQ123456789",
  "order": {
    "id": 1,
    "user_id": 1,
    "address_id": 10,
    "order_no": "SK202603251200001",
    "status": "pending",
    "total_amount": 19900,
    "idempotency_key": "seckill-20260325-1",
    "items": [
      {
        "id": 11,
        "order_id": 1,
        "product_id": 1001,
        "product_name": "秒杀商品",
        "product_price": 19900,
        "quantity": 1
      }
    ]
  }
}
```

补充：

- 当 `KAFKA_ENABLED=false` 时，当前服务会同步创建订单，`queued=false`
- 当 `KAFKA_ENABLED=true` 时，消息会投递到 Kafka，再由后台消费者处理，此时可根据 `message_id` 和订单列表轮询结果
```

## Health Check

### GET `/ping`

说明：服务健康检查

成功响应 `data`：

```json
{
  "service": "shop-go",
  "env": "dev"
}
```

### GET `/metrics`

说明：Prometheus 指标接口

用途：

- 供监控系统抓取
- 不用于移动端业务调用

## Error Handling

移动端建议统一按以下逻辑处理：

- 先看 HTTP 状态码
- 再看响应体 `code`
- `code != 0` 时直接展示 `msg`
- `401` 时触发登录失效或刷新 token 逻辑

## Notes

- 当前删除地址、删除购物车项都用 query 参数传 `id`
- 当前金额字段为整数，建议前端按“分”展示
- 当前订单状态先只返回 `pending`
- 当前未生成 Swagger，本文件就是当前对接基准
