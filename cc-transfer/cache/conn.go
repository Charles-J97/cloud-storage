package cache

import (
	"cc-transfer/config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

var(
	pool *redis.Pool
	redisHost = config.RedisHost
	//redisPwd = config.RedisPwd
)

//创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 50,
		MaxActive: 30,
		IdleTimeout: 300 * time.Second,
		//Dial方法用于创建连接
		Dial: func() (redis.Conn, error) {
			//1 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println("Fail to connect to redis, error: " + err.Error())
				return nil, err
			}
			////2 访问认证
			//_, err = c.Do("AUTH", redisPwd)
			//if err != nil {
			//	c.Close()
			//	fmt.Print(err.Error())
			//	return nil, err
			//}
			return c, nil
		},
		//定时测试redis连接是否正常
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
	_, err := pool.Get().Do("PING")
	if err != nil {
		fmt.Println("Fail to connect to redis, error: " + err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Redis is connected successfully!")
	}
}

//pool在外部是访问不到的，故需要将其暴露给外部
func RedisPool() *redis.Pool {
	return pool
}
