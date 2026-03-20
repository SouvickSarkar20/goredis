package cmd

import (
	"goredis/resp"
	"goredis/store"
)

// HSET user:1 name Souvick
// HSET user:1 age 21

func handleHAdd(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) < 3 {
		return w.WriteError("ERR : wrong no of arguments")
	}

	key := input.Array[1].Str
	field := input.Array[2].Str

	for i := 2; i < len(input.Array); i++ {
		value := input.Array[i].Str
		error := db.HSet(key, field, value)
		if error != nil {
			return w.WriteError(error.Error())
		}
	}

	return w.WriteSimpleString("OK")
}

func handleHGet(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) < 3 {
		return w.WriteError("ERR : wrong no of arguments")
	}

	key := input.Array[1].Str
	field := input.Array[2].Str

	value, exists, err := db.HGet(key, field)

	if !exists {
		return w.WriteError("The value does not exist")
	}

	if err != nil {
		return w.WriteError(err.Error())
	}

	return w.WriteBulkString(value)
}

// HDEL user:1 name
func handleHDel(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) < 2 {
		return w.WriteError("ERR : wrong no of arguments")
	}

	key := input.Array[1].Str
	field := input.Array[2].Str

	_, err := db.HDel(key, field)

	if err != nil {
		return w.WriteError(err.Error())
	}

	return w.WriteSimpleString("OK")

}
