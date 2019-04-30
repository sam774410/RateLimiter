Middleware
===

### 環境建置
* MacOS Mojave: 10.14
* Golang: 1.11.5
* Framework: Gin
* DB: Redis: 5.0.3

### 建置說明
* 限制同一IP來源請求每小時最多1000次
* 一小時內超過1000次請求回傳429(Too many requests)
* 每次請求中 response headers 加入剩餘請求次數(X-RateLimit-Remaining) &
  重置時間(X-RateLimit-Reset)

### 啟動環境

1. 激活網站 [web](https://ratelimiter-redis.herokuapp.com/)
2. 查看網站 [web](https://ratelimiter-redis.herokuapp.com/draw) 及   response headers