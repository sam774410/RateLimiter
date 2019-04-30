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
* response headers中加入剩餘請求次數(X-RateLimit-Remaining) &
  重置時間(X-RateLimit-Reset)

### 啟動環境

1. cd至middle資料夾
2. 開啟local Redis DB (預設port: 6379)
4. **啟動 main 執行檔** or **go run main.go**
5. 開啟連結 127.0.0.1:2000/draw
6. 查看response headers

