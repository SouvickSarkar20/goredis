package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Value represents any parsed RESP value.
type Value struct {
	Typ   string  // "SimpleString", "Error", "Integer", "BulkString", "Array"
	Str   string  // Holds the string data (for SimpleString, Error, BulkString)
	Num   int     // Holds the integer data (for Integer)
	Array []Value // Holds the array elements (for Array)
}

// Parser wraps a bufio.Reader to read from the TCP stream.
type Parser struct {
	reader *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

// ParseOne reads exactly one complete RESP value from the network stream.
func (p *Parser) ParseOne() (Value, error) {
	// Read a single byte to see what kind of RESP type we are dealing with.
	typ, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case '*':
		return p.parseArray()
	case '$':
		return p.parseBulkString()
	default:
		// GoRedis currently only expects clients to send Arrays of BulkStrings.
		return Value{}, fmt.Errorf("unknown RESP type byte: %c", typ)
	}
}

// readLine is a helper to read bytes until we hit '\r\n'.
// It returns the bytes BEFORE the '\r\n', and strips the '\r\n' off.
func (p *Parser) readLine() (line []byte, n int, err error) {
	// ReadBytes('\n') keeps reading until it sees a '\n' character.
	line, err = p.reader.ReadBytes('\n')
	if err != nil {
		return nil, 0, err
	}
	
	n = len(line)
	// We expect lines to end in '\r\n'.
	// So we need to trim the last 2 bytes off the string we return.
	if n >= 2 && line[n-2] == '\r' {
		return line[:n-2], n, nil
	}

	return line, n, nil
}

// parseArray reads a RESP Array (e.g. "*2\r\n").
func (p *Parser) parseArray() (Value, error) {
	// 1. Read the length line (e.g. "2")
	line, _, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	// 2. Convert the string "2" into an actual integer 2
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	// 3. Create a Go slice to hold that many elements
	v := Value{
		Typ:   "Array",
		Array: make([]Value, length),
	}

	// 4. Loop exactly 'length' times, calling ParseOne to read each element
	for i := 0; i < length; i++ {
		val, err := p.ParseOne()
		if err != nil {
			return Value{}, err
		}
		v.Array[i] = val
	}

	return v, nil
}

// parseBulkString reads a RESP Bulk String (e.g. "$5\r\nhello\r\n").
func (p *Parser) parseBulkString() (Value, error) {
	// 1. Read the length line (e.g. "5")
	line, _, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	// 2. Convert "5" into the integer 5
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	// 3. Now we know exactly how many bytes the string is!
	// We make a byte array of EXACTLY that size.
	bulkBytes := make([]byte, length)

	// io.ReadFull guarantees it will block until it fills the EXCAT size of 'bulkBytes'.
	_, err = io.ReadFull(p.reader, bulkBytes)
	if err != nil {
		return Value{}, err
	}

	// 4. RESP Bulk strings always end with a trailing "\r\n".
	// We need to read those 2 bytes and throw them away.
	p.readLine()

	return Value{
		Typ: "BulkString",
		Str: string(bulkBytes),
	}, nil
}
