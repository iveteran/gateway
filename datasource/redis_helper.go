package datasource

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"matrix.works/fmx-gateway/conf"
)

var redisInstance *RedisConn
var redisLock sync.Mutex

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (this *RedisConn) Do(
	commandName string,
	args ...interface{},
) (reply interface{}, err error) {

	conn := this.pool.Get()
	defer conn.Close() // 将连接放回连接池

	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)

	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Fatal("redis_helper.Do error ", err, e)
		}
	}

	t2 := time.Now().UnixNano()

	if this.showDebug {
		fmt.Printf(
			"[redis] [info] [%dus] cmd=%s, args=%v, reply=%s, err=%s\n",
			(t2-t1)/1000, commandName, args, reply, err,
		)
	}

	return reply, err

}

func (this *RedisConn) ShowDebug(show bool) {
	this.showDebug = show
}

func GetRedisInstance() *RedisConn {

	if redisInstance == nil {
		redisLock.Lock()
		defer redisLock.Unlock()
		redisInstance = NewRedisInstance()
	}
	return redisInstance

}

func NewRedisInstance() *RedisConn {
	redisAddress := fmt.Sprintf("%s:%d", conf.Cfg.Redises["default"].Host, conf.Cfg.Redises["default"].Port)
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				log.Fatal("redis_helper.NewRedisInstance error ", err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000, // 最多连接数
		MaxActive:       10000, // 最多活跃数
		IdleTimeout:     0,     // 超时时间
		Wait:            false, // 连接等待
		MaxConnLifetime: 0,     //最大连接时间，0 一直连接
	}

	conn := &RedisConn{pool: &pool}
	//conn.ShowDebug(true)
	log.Printf("[NewRedisInstance] %s", redisAddress)

	return conn
}
