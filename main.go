package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// 连接池大小
var MAX_POOL_SIZE = 20
var redisPoll chan redis.Conn

func putRedis(conn redis.Conn) {
	// 基于函数和接口间互不信任原则，这里再判断一次，养成这个好习惯哦
	if redisPoll == nil {
		redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
	}
	if len(redisPoll) >= MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPoll <- conn
}
func InitRedis(network, address string) redis.Conn {
	// 缓冲机制，相当于消息队列
	if len(redisPoll) == 0 {
		// 如果长度为0，就定义一个redis.Conn类型长度为MAX_POOL_SIZE的channel
		redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
		go func() {
			for i := 0; i < MAX_POOL_SIZE/2; i++ {
				c, err := redis.Dial(network, address)
				if err != nil {
					panic(err)
				}
				putRedis(c)
			}
		}()
	}
	return <-redisPoll
}
func main() {
	c := InitRedis("tcp", "192.168.1.213:6379")
	dbkey := "netgame:info"
	if ok, err := redis.Bool(c.Do("LPUSH", dbkey, "yanetao")); ok {
	} else {
		log.Print(err)
	}
	fmt.Println("push sucessful ")

	values, _ := redis.Values(c.Do("lrange", "netgame:info", "0", "100"))
	var v1 string
	redis.Scan(values, &v1)
	fmt.Println(v1)

	_, err := c.Do("SET", "name", "red")
	if err != nil {
		fmt.Println(err)
		return
	}

	v, err := redis.String(c.Do("GET", "name"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
