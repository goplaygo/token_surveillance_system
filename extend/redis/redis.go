/*
@Time : 2020/7/7 1:30 下午
@Author : L
@File : redis.go
@Software: GoLand
*/
package redis

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	"token_surveillance_system/extend/conf"
)

var redisConn *redis.Pool

func GetRedisConn() *redis.Pool{
	return redisConn
}

func Setup()error {
	redisConn = &redis.Pool{
		MaxActive: conf.RedisConf.MaxActive,
		MaxIdle: conf.RedisConf.MaxIdle,
		IdleTimeout: conf.RedisConf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c,err:=redis.Dial("tcp",conf.RedisConf.Host+":"+strconv.Itoa(conf.RedisConf.Port))
			if err != nil {
				return nil, err
			}
			//验证密码
			if conf.RedisConf.Password !=""{
				if _,err:=c.Do("AUTH",conf.RedisConf.Password);err!=nil{
					c.Close()
					return nil,err
				}
			}
			return c,nil
		},
		TestOnBorrow:func(c redis.Conn,t time.Time)error{
			if time.Since(t) < time.Minute {
				return nil
			}
			_,err := c.Do("PING")
			return err
		},
	}
	return nil
}

//set
func Set(key string,data string,seconds int) error{
	conn:=GetRedisConn().Get()
	defer conn.Close()
	_,err:=conn.Do("SET",key,data)
		if err != nil {
			return err
		}

	_,err=conn.Do("EXPIRE",key,seconds)
	if err != nil {
		return err
	}
	return nil
}

//exists方法
func Exists(key string) bool {
	conn := GetRedisConn().Get()
	defer conn.Close()
	exists,err:=redis.Bool(conn.Do("EXISTS",key))
	if err != nil {
		return false
	}
	return exists
}

//Get
func Get(key string)(string, error ){
	conn := GetRedisConn().Get()
	defer conn.Close()
	reply,err:=redis.String(conn.Do("GET",key))
	if err!=nil && err!=redis.ErrNil{
		return "",err
	}
	if err == redis.ErrNil {
		return "",nil
	}
	return reply,nil
}

func Del(key string)(bool,error){
	conn:=GetRedisConn().Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL",key))
}

func DeLike(key string)error{
	conn:=GetRedisConn().Get()
	defer conn.Close()
	keys,err:=redis.Strings(conn.Do("KEYS","*"+key+"*"))
	if err != nil {
		return err
	}
	for _,key:=range keys {
		_,err:= Del(key)
		if err!= nil {
			return err
		}
	}
	return nil
}
