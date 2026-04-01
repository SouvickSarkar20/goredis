package cmd

import (
	"github.com/krishsinghhura/go-redis/resp"
	"github.com/krishsinghhura/go-redis/store"
)

func handleLPush(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) < 3 {
		return w.WriteError("ERR wrong number of arguments for 'lpush' command")
	}

	key := args.Array[1].Str

	for i := 2; i < len(args.Array); i++ {
		value := args.Array[i].Str
		err := db.LPush(key, value)
		if err != nil {
			return w.WriteError(err.Error())
		}
	}

	return w.WriteSimpleString("OK")
}

func handleLPop(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'lpop' command")
	}

	key := args.Array[1].Str

	val, exists, err := db.LPop(key)
	if err != nil {
		return w.WriteError(err.Error())
	}

	if !exists {
		return w.WriteBulkStringNil()
	}

	return w.WriteBulkString(val)
}
