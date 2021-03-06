package tool

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var RedisConn *redis.Pool

type RedisConf struct {
	Host        string        `yaml:"host"`
	Password    string        `yaml:"password"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxActive   int           `yaml:"maxActive"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

type RedisDB int

// Setup Initialize the Redis instance
func EnableRedis(conf RedisConf) error {
	RedisConn = &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeout,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", conf.Host)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// Set a key/value
func (r RedisDB) Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if time != 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}

	}

	return nil
}

// Get get a key
func (r RedisDB) Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r RedisDB) SetString(key string, value string, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if time != 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}

	}
	return nil
}
func (r RedisDB) GetString(key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}
func (r RedisDB) SetBool(key string, value bool, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if time != 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}

	}
	return nil
}
func (r RedisDB) GetBool(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("GET", key))
}

func (r RedisDB) GetInt(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GET", key))
}
func (r RedisDB) SetInt(key string, data int, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}
	if time != 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists check a key
func (r RedisDB) Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func (r RedisDB) Incre(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Delete delete a kye
func (r RedisDB) Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func (r RedisDB) LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = r.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r RedisDB) ZAdd(key string, value, score int) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, value)
	return err
}

func (r RedisDB) ZRank(key string, value int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	rank, err := redis.Int(conn.Do("ZRANK", key, strconv.Itoa(value)))
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (r RedisDB) ZRange(key string, offset int) ([]int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	ranks, err := redis.Ints(conn.Do("ZRANGE", key, 0, offset))
	if err != nil {
		return nil, err
	}

	return ranks, nil
}

// 增加计数
func (r RedisDB) INCRBY(key string, increment int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	newNum, err := redis.Int(conn.Do("INCRBY", key, increment))
	if err != nil {
		return 0, err
	}
	return newNum, nil
}

// 获取计数
func (r RedisDB) GetCounter(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	count, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 初始化计数器
func (r RedisDB) InitCounter(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, 0)
	if err != nil {
		return err
	}
	return nil
}

var ErrCounterHasBeenZero = errors.New("the counter has been 0")

// 减少计数
func (r RedisDB) DECRBY(key string, decrement int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	count, err := r.GetCounter(key)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, ErrCounterHasBeenZero
	}
	newNum, err := redis.Int(conn.Do("DECRBY", key, decrement))
	if err != nil {
		return 0, err
	}
	return newNum, nil
}
