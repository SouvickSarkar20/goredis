package cmd

import (
	"goredis/resp"
	"goredis/store"
	"strings"
)

type HandlerFunc func(*resp.Writer, *store.Store, resp.Value) error

var routes = map[string]HandlerFunc{
	"PING":  handlePing,
	"SET":   handleSet,
	"GET":   handleGet,
	"DEL":   handleDel,
	"LPUSH": handleLPush,
	"LPOP":  handleLPop,
	"HGET":  handleHGet,
	"HSET":  handleHAdd,
	"HDEL":  handleHDel,
}

func Handle(writer *resp.Writer, db *store.Store, input resp.Value) error {
	command := strings.ToUpper(input.Array[0].Str)

	handler, exists := routes[command]
	if !exists {
		return writer.WriteError("ERR unkown command" + command)
	}

	return handler(writer, db, input)
}

func handlePing(w *resp.Writer, db *store.Store, args resp.Value) error {
	return w.WriteSimpleString("PONG")
}
