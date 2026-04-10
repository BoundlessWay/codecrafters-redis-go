package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseArrayContent(reader *bufio.Reader) ([]string, error) {
	count, err := readInt(reader)
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		return []string{}, nil
	}

	args := make([]string, count)
	for i := 0; i < count; i++ {
		arg, err := readBulkString(reader)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	return args, nil
}

func readBulkString(reader *bufio.Reader) (string, error) {
	prefix, err := reader.ReadByte()
	if err != nil || prefix != '$' {
		return "", fmt.Errorf("expected '$'")
	}

	size, err := readInt(reader)
	if err != nil {
		return "", err
	}

	if size == -1 {
		return "", nil
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(reader, data); err != nil {
		return "", err
	}

	reader.ReadString('\n')
	return string(data), nil
}

func readInt(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(line))
}
