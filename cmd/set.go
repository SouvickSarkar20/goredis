package cmd

import (
	"github.com/krishsinghhura/goredis/resp"
	"github.com/krishsinghhura/goredis/store"
	"strings"
)

func handleSAdd(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) != 3 {
		return w.WriteError("ERR wrong number of arguments for 'SADD'")
	}

	key := input.Array[1].Str
	value := input.Array[2].Str

	if err := db.SAdd(key, value); err != nil {
		return w.WriteError(err.Error())
	}

	return w.WriteInteger(1)
}

func handleSMembers(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'SMEMBERS'")
	}

	key := input.Array[1].Str

	members, err := db.SMembers(key)
	if err != nil {
		return w.WriteError(err.Error())
	}

	if err := w.WriteArray(len(members)); err != nil {
		return err
	}

	for _, m := range members {
		if err := w.WriteBulkString(m); err != nil {
			return err
		}
	}

	return nil
}

func handleSIsMember(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) != 3 {
		return w.WriteError("ERR wrong number of arguments for 'SISMEMBER'")
	}

	key := input.Array[1].Str
	value := input.Array[2].Str

	ok, err := db.SIsMember(key, value)
	if err != nil {
		return w.WriteError(err.Error())
	}

	if ok {
		return w.WriteInteger(1)
	}
	return w.WriteInteger(0)
}

func handleSRem(w *resp.Writer, db *store.Store, input resp.Value) error {
	if len(input.Array) != 3 {
		return w.WriteError("ERR wrong number of arguments for 'SREM'")
	}

	key := input.Array[1].Str
	value := input.Array[2].Str

	removed, err := db.SRem(key, value)
	if err != nil {
		if strings.HasPrefix(err.Error(), "WRONGTYPE") {
			return w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return w.WriteError(err.Error())
	}

	if removed {
		return w.WriteInteger(1)
	}
	return w.WriteInteger(0)
}
