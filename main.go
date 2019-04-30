package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm runnung~")
	})
	r.GET("/draw", RateLimiter, Draw)

	err := r.Run(":2000")
	if err != nil {
		log.Fatal("server start error: ", err)
	}
}

var (
	//最大請求次數
	MAXREQ = 1000

	//Expire time 1 hr
	EXPTIME = "3600"
)

func Draw(c *gin.Context) {

	c.JSON(http.StatusOK, "draw it~")
}

func RateLimiter(c *gin.Context) {

	//conn redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("cannot connect to redis: ", err)
		return
	}
	defer conn.Close()

	log.Println("client IP: ", c.ClientIP())

	//check ip in redis or not
	v, err := redis.String(conn.Do("GET", c.ClientIP()))
	if err != nil {
		log.Println("command GET ", c.ClientIP(), " error: ", err)
	}

	if len(v) == 0 {
		// not in redis or expire -> go through

		//set value (init)
		_, err := redis.String(conn.Do("SET", c.ClientIP(), "1"))
		if err != nil {
			log.Println("command SET ", c.ClientIP(), " error: ", err)
		}

		//expire time 1hr
		_, err = redis.String(conn.Do("EXPIRE", c.ClientIP(), EXPTIME))
		if err != nil {
			log.Println("command EXPIRE ", c.ClientIP(), " error: ", err)
		}

		//set header
		c.Writer.Header().Set("X-RateLimit-Remaining", "999")
		c.Writer.Header().Set("X-RateLimit-Reset", EXPTIME)

		c.Next()
	} else {
		// IP already in redis
		// check value 1000->429 or return X-RateLimit-Remaining、 X-RateLimit-Reset

		//get value
		v, err := redis.String(conn.Do("GET", c.ClientIP()))
		if err != nil {
			log.Println("command GET ", c.ClientIP(), " error: ", err)
		}
		log.Println("current value: ", v)

		reqs, _ := strconv.Atoi(v)
		if reqs == MAXREQ {
			// return 429
			//ttl
			//time to live (sec)
			remainTime, err := redis.Int(conn.Do("TTL", c.ClientIP()))
			if err != nil {
				log.Println("command TTL error: ", err)
			}

			s := fmt.Sprintf("%d sec", remainTime)

			//set header
			//remainT := strconv.Itoa(remainTime)

			c.Writer.Header().Set("X-RateLimit-Remaining", "0")
			c.Writer.Header().Set("X-RateLimit-Reset", strconv.Itoa(remainTime))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"description": "too more requests, more than 1000 in an hour",
				"wait":        s,
			})

			c.Abort()
		} else {

			//increase value
			_, err := redis.String(conn.Do("INCR", c.ClientIP()))
			if err != nil {
				log.Println("command INCR error: ", err)
			}

			//time to live (sec)
			remainTime, err := redis.Int(conn.Do("TTL", c.ClientIP()))
			if err != nil {
				log.Println("command TTL error: ", err)
			}

			//remain request
			v, err = redis.String(conn.Do("GET", c.ClientIP()))
			if err != nil {
				log.Println("command GET ", c.ClientIP(), " error: ", err)
			}

			remainReq, _ := strconv.Atoi(v)

			log.Println("remain: ", remainTime, " req: ", 1000-remainReq)

			//set header
			c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(1000-remainReq))
			c.Writer.Header().Set("X-RateLimit-Reset", strconv.Itoa(remainTime))

			c.Next()
		}
	}

}
