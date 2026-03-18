package cmd

import (
	"goredis/resp"
	"goredis/store"
)

func handleDel(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'del' command")
	}

	key := args.Array[1].Str
	db.Delete(key)

	return w.WriteInteger(1)
}
