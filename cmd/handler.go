package cmd

import (
	"goredis/resp"
	"goredis/store"
	"strconv"
	"strings"
	"time"
)

func Handle(writer *resp.Writer, db *store.Store, input resp.Value) error {
	command := strings.ToUpper(input.Array[0].Str)

	switch command {
	case "PING":
		return handlePing(writer)
	case "SET":
		return handleSet(writer, db, input)
	case "GET":
		return handleGet(writer, db, input)
	case "DEL":
		return handleDel(writer, db, input)
	default:
		return writer.WriteError("ERR unknown command '" + command + "'")
	}
}

func handlePing(w *resp.Writer) error {
	return w.WriteSimpleString("PONG")
}

func handleSet(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) < 3 {
		return w.WriteError("ERR wrong number of arguments for 'set' command")
	}

	key := args.Array[1].Str
	value := args.Array[2].Str

	var duration time.Duration

	// Check for optional arguments like EX (Expiry in seconds)
	if len(args.Array) >= 5 {
		option := strings.ToUpper(args.Array[3].Str)
		if option == "EX" {
			seconds, err := strconv.Atoi(args.Array[4].Str)
			if err != nil {
				return w.WriteError("ERR value is not an integer or out of range")
			}
			duration = time.Duration(seconds) * time.Second
		}
	}

	db.Set(key, value, duration)

	return w.WriteSimpleString("OK")
}

func handleGet(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'get' command")
	}

	key := args.Array[1].Str

	val, exists := db.Get(key)

	if !exists {
		return w.WriteBulkStringNil()
	}

	return w.WriteBulkString(val)
}

func handleDel(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'del' command")
	}

	key := args.Array[1].Str
	db.Delete(key)

	return w.WriteInteger(1)
}
