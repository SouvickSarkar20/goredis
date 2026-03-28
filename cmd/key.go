package cmd

import (
	"goredis/resp"
	"goredis/store"
)

func handleDel(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) < 2 {
		return w.WriteError("ERR wrong number of arguments for 'del' command")
	}

	count := 0
	for i := 1; i < len(args.Array); i++ {
		key := args.Array[i].Str
		if db.Delete(key) {
			count++
		}
	}

	return w.WriteInteger(int64(count))
}
