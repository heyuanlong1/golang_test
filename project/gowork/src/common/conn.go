package common

import (
	"fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/hoisie/redis"
    CONF "conf"
    "strconv"
)


var (
    client *redis.Client
)

func GetMysql()(*sql.DB,error) {
	var ss string
	ss = fmt.Sprintf("%s:%s@tcp(%s:%d)/ko_resource?charset=utf8",
		CONF.Mysql_user_name,
		CONF.Mysql_user_pass,
		CONF.Mysql_server_ip,
		CONF.Mysql_server_port,
		)
	db, err := sql.Open("mysql", ss)
	return db,err
}

func GetRedis()( *redis.Client,error) {
	client = &redis.Client{
        Addr:        CONF.Redis_server_ip + ":" + strconv.Itoa(CONF.Redis_server_port),
        Db:          0, // default db is 0
        Password:    CONF.Redis_server_pass,
        MaxPoolSize: 10000,
    }

	if err := client.Auth(client.Password); err != nil {      
        return nil,err
    }
    return client , nil
}
