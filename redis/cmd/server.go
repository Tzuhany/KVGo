package main

import (
	"github.com/tidwall/redcon"
	"kvgo"
	kbgoredis "kvgo/redis"
	"log"
	"sync"
)

const addr = "127.0.0.1:6380"

type KVGoServer struct {
	dbs    map[int]*kbgoredis.RedisDataStructure
	server *redcon.Server
	mu     sync.RWMutex
}

func main() {
	// 打开 Redis 数据结构服务
	redisDataStructure, err := kbgoredis.NewRedisDataStructure(kvgo.DefaultOptions)
	if err != nil {
		panic(err)
	}

	// 初始化 KVGoServer
	kvgoServer := &KVGoServer{
		dbs: make(map[int]*kbgoredis.RedisDataStructure),
	}
	kvgoServer.dbs[0] = redisDataStructure

	// 初始化一个 Redis 服务端
	kvgoServer.server = redcon.NewServer(addr, execClientCommand, kvgoServer.accept, kvgoServer.close)
	kvgoServer.listen()
}

func (svr *KVGoServer) listen() {
	log.Println("kvgo server running, ready to accept connections.")
	_ = svr.server.ListenAndServe()
}

func (svr *KVGoServer) accept(conn redcon.Conn) bool {
	cli := new(KVGoClient)
	svr.mu.Lock()
	defer svr.mu.Unlock()
	cli.server = svr
	cli.db = svr.dbs[0]
	conn.SetContext(cli)
	return true
}

func (svr *KVGoServer) close(conn redcon.Conn, err error) {
	for _, db := range svr.dbs {
		_ = db.Close()
	}
}
