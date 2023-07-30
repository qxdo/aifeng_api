package utils

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisClient struct {
	client *redis.Pool
}

var Redis *RedisClient

// InitClient 初始化redis对象
func InitClient(host string, port int, password string, db int) error {
	pwd := redis.DialPassword(password)
	dbIndex := redis.DialDatabase(db)
	client := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), pwd, dbIndex)
		},
		MaxIdle:         50,
		MaxActive:       150,
		IdleTimeout:     30 * time.Second,
		Wait:            true,
		MaxConnLifetime: 10 * time.Second,
	}
	Redis = &RedisClient{
		client: client,
	}
	return nil
}

func (r *RedisClient) GetString(key string) (string, error) {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	res, err := redis.String(client.Do("get", key))
	return res, err
}

func (r *RedisClient) SetString(key, value string) error {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	_, err := redis.String(client.Do("set", key, value))
	return err
}

func (r *RedisClient) ClearList(key string) error {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	_, err := client.Do("ltrim", key, 1, 0)
	return err
}

func (r *RedisClient) GetKeys(pattern string) ([]string, error) {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	res, err := redis.Strings(client.Do("keys", pattern))
	return res, err
}

func (r *RedisClient) GetList(key string) ([]string, error) {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	length, err := redis.Int(client.Do("llen", key))
	if err != nil {
		return nil, err
	}
	res, err := redis.Strings(client.Do("lrange", key, 0, length))
	return res, err
}

func (r *RedisClient) RemoveItem(key, value string) error {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	_, err := client.Do("lrem", key, 0, value)
	return err
}

func (r *RedisClient) DelItem(key string) error {
	client := r.client.Get()
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("redis client close error: ", err.Error())
		}
	}()
	_, err := client.Do("del", key)
	return err
}
