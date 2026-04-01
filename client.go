package goredis

import (
	"fmt"
	"net"

	"github.com/krishsinghhura/go-redis/resp"
)

// Client represents a GoRedis client.
type Client struct {
	conn   net.Conn
	parser *resp.Parser
	writer *resp.Writer
}

// NewClient creates a new GoRedis client connecting to the given address.
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	return &Client{
		conn:   conn,
		parser: resp.NewParser(conn),
		writer: resp.NewWriter(conn),
	}, nil
}

// Close closes the connection to the server.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Set sets a key to a string value.
func (c *Client) Set(key string, value string) error {
	cmd := resp.Value{
		Typ: "Array",
		Array: []resp.Value{
			{Typ: "BulkString", Str: "SET"},
			{Typ: "BulkString", Str: key},
			{Typ: "BulkString", Str: value},
		},
	}

	err := c.writer.Write(cmd)
	if err != nil {
		return err
	}

	res, err := c.parser.ParseOne()
	if err != nil {
		return err
	}

	if res.Typ == "Error" {
		return fmt.Errorf("server error: %s", res.Str)
	}

	return nil
}

// Get retrieves a string value by key.
func (c *Client) Get(key string) (string, error) {
	cmd := resp.Value{
		Typ: "Array",
		Array: []resp.Value{
			{Typ: "BulkString", Str: "GET"},
			{Typ: "BulkString", Str: key},
		},
	}

	err := c.writer.Write(cmd)
	if err != nil {
		return "", err
	}

	res, err := c.parser.ParseOne()
	if err != nil {
		return "", err
	}

	if res.Typ == "Error" {
		return "", fmt.Errorf("server error: %s", res.Str)
	}

	// For BulkString, empty Str means nil/not found if it's the RESP convention,
	// but our parser handles nil by returning an empty Str or a specialized value.
	return res.Str, nil
}

// Del deletes a key. Returns the number of keys deleted (0 or 1 for standard Redis).
func (c *Client) Del(key string) (int, error) {
	cmd := resp.Value{
		Typ: "Array",
		Array: []resp.Value{
			{Typ: "BulkString", Str: "DEL"},
			{Typ: "BulkString", Str: key},
		},
	}

	err := c.writer.Write(cmd)
	if err != nil {
		return 0, err
	}

	res, err := c.parser.ParseOne()
	if err != nil {
		return 0, err
	}

	if res.Typ == "Error" {
		return 0, fmt.Errorf("server error: %s", res.Str)
	}

	return res.Num, nil
}

// HSet sets a field in a hash to a value.
func (c *Client) HSet(key, field, value string) error {
	cmd := resp.Value{
		Typ: "Array",
		Array: []resp.Value{
			{Typ: "BulkString", Str: "HSET"},
			{Typ: "BulkString", Str: key},
			{Typ: "BulkString", Str: field},
			{Typ: "BulkString", Str: value},
		},
	}

	err := c.writer.Write(cmd)
	if err != nil {
		return err
	}

	res, err := c.parser.ParseOne()
	if err != nil {
		return err
	}

	if res.Typ == "Error" {
		return fmt.Errorf("server error: %s", res.Str)
	}

	return nil
}

// HGet retrieves a field's value from a hash.
func (c *Client) HGet(key, field string) (string, error) {
	cmd := resp.Value{
		Typ: "Array",
		Array: []resp.Value{
			{Typ: "BulkString", Str: "HGET"},
			{Typ: "BulkString", Str: key},
			{Typ: "BulkString", Str: field},
		},
	}

	err := c.writer.Write(cmd)
	if err != nil {
		return "", err
	}

	res, err := c.parser.ParseOne()
	if err != nil {
		return "", err
	}

	if res.Typ == "Error" {
		return "", fmt.Errorf("server error: %s", res.Str)
	}

	return res.Str, nil
}
