package cmd

import (
	"github.com/SouvickSarkar20/goredis/resp"
	"github.com/SouvickSarkar20/goredis/store"
)

// HSET user:1 name Souvick
// HSET user:1 age 21

func handleHAdd(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) < 3 {
		return w.WriteError("ERR : wrong no of arguments")
	}

	key := input.Array[1].Str

	for i := 2; i < len(input.Array); i += 2 {
		if i+1 >= len(input.Array) {
			return w.WriteError("ERR syntax error")
		}
		field := input.Array[i].Str
		value := input.Array[i+1].Str
		err := db.HSet(key, field, value)
		if err != nil {
			return w.WriteError(err.Error())
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

	if err != nil {
		return w.WriteError(err.Error())
	}

	if !exists {
		return w.WriteBulkString("") // Return empty string if field does not exist
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
