package datasource

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"matrix.works/gateway/conf"
)

var redisInstances map[string]*RedisConn
var redisLock sync.Mutex

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (this *RedisConn) GetConnection() redis.Conn {
	return this.pool.Get()
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
			log.Fatal("[redis_helper.Do] error ", err, e)
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

func GetRedisDefaultInstance() *RedisConn {
	return GetRedisInstance("default")
}

func GetRedisInstance(instanceName string) *RedisConn {
	if redisInstances == nil {
		redisInstances = make(map[string]*RedisConn)
	}

	if redisInstance, exist := redisInstances[instanceName]; exist {
		return redisInstance
	} else {
		redisLock.Lock()
		defer redisLock.Unlock()
		redisInstance = NewRedisInstance(instanceName)
		redisInstances[instanceName] = redisInstance
		return redisInstance
	}
}

func NewRedisInstance(instanceName string) *RedisConn {
	if _, exist := conf.Cfg.Redises[instanceName]; !exist {
		log.Fatalf("[NewRedisInstance] Error: The configure of Redis(%s) does not exist", instanceName)
		return nil
	}
	redisConfig := conf.Cfg.Redises[instanceName]
	host := redisConfig.Host
	port := redisConfig.Port
	db := redisConfig.Database
	if host == "" || port == 0 {
		log.Fatalf("[NewRedisInstance] Error: Missed required parameter(s) of Redis(%s) configure", instanceName)
		return nil
	}

	redisAddress := fmt.Sprintf("%s:%d", host, port)
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			setDb := redis.DialDatabase(db)
			c, err := redis.Dial("tcp", redisAddress, setDb)
			if err != nil {
				log.Fatalf("[NewRedisInstance] Error: Dial up error: %s", err.Error())
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
	conn.ShowDebug(true)
	log.Printf("[NewRedisInstance] Connected to %s server: %s/%d", instanceName, redisAddress, db)

	return conn
}
