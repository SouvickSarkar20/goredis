package cmd

import (
	"strings"

	"goredis/resp"
	"goredis/store"
)

// Handle routes the parsed RESP array to the correct command function.
func Handle(writer *resp.Writer, db *store.Store, input resp.Value) error {
	// The client might send "set", "SET", or "sEt".
	// Real Redis is case-insensitive for command names.
	// We extract the first string in the array (the command name) and uppercase it.
	command := strings.ToUpper(input.Array[0].Str)

	// Switch on the command name to route it.
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
		// If the command is unknown, we send a RESP error (-ERR)
		return writer.WriteError("ERR unknown command '" + command + "'")
	}
}

// handlePing responds with "PONG"
func handlePing(w *resp.Writer) error {
	return w.WriteSimpleString("PONG")
}

// handleSet implements the SET command: SET key value
func handleSet(w *resp.Writer, db *store.Store, args resp.Value) error {
	// 1. Argument Validation
	// SET requires exactly 3 tokens: ["SET", "key", "value"]
	if len(args.Array) != 3 {
		return w.WriteError("ERR wrong number of arguments for 'set' command")
	}

	key := args.Array[1].Str
	value := args.Array[2].Str

	// 2. Perform the database write
	db.Set(key, value)

	// 3. Return success to the client
	return w.WriteSimpleString("OK")
}

// handleGet implements the GET command: GET key
func handleGet(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'get' command")
	}

	key := args.Array[1].Str

	// Perform the database read
	val, exists := db.Get(key)
	
	if !exists {
		// In Redis, if a key doesn't exist, it returns a "Nil Bulk String".
		// In RESP, a nil string is represented as "$-1\r\n".
		// We haven't built WriteNull() in the writer yet, so let's just write the raw bytes for now!
		return w.WriteBulkStringNil()
	}

	// If it exists, return it as a Bulk String
	return w.WriteBulkString(val)
}

// handleDel implements the DEL command: DEL key
func handleDel(w *resp.Writer, db *store.Store, args resp.Value) error {
	if len(args.Array) != 2 {
		return w.WriteError("ERR wrong number of arguments for 'del' command")
	}

	key := args.Array[1].Str
	db.Delete(key)

	// DEL historically returns the number of keys deleted (an integer),
	// but we haven't implemented integer writing yet.
	// For now, let's just send back an Integer "1".
	return w.WriteInteger(1)
}
