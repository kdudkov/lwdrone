package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Config struct {
	ch         byte
	flip       byte
	wifiSec    byte
	wifiName   string
	wifiPass   string
	time       uint64
	sdcMounted byte
	sdcSize    uint64
	sdcFree    uint64
	version    string
}

func ConfigFromBytes(data []byte) (*Config, error) {
	if len(data) != 140 {
		return nil, fmt.Errorf("invalid len: %d (must be 140)", len(data))
	}

	le := binary.LittleEndian
	c := &Config{}
	c.ch = data[0]
	c.flip = data[1]
	c.wifiSec = data[2]
	n := 3
	c.wifiName = z2s(string(data[n : n+32]))
	n += 32
	c.wifiPass = z2s(string(data[n : n+32]))
	n += 32
	c.time = le.Uint64(data[n:])
	n += 8
	c.sdcMounted = data[n]
	n++
	c.sdcSize = le.Uint64(data[n:])
	n += 8
	c.sdcFree = le.Uint64(data[n:])
	n += 8
	c.version = z2s(string(data[n:]))
	return c, nil
}

func (c *Config) ToBytes() []byte {
	res := make([]byte, 140)
	le := binary.LittleEndian
	res[0] = c.ch
	res[1] = c.flip
	res[2] = c.wifiSec
	n := 3
	copy(res[n:], c.wifiName)
	n += 32
	copy(res[n:], c.wifiPass)
	n += 32
	le.PutUint64(res[n:], c.time)
	n += 8
	res[n] = c.sdcMounted
	n++
	le.PutUint64(res[n:], c.sdcSize)
	n += 8
	le.PutUint64(res[n:], c.sdcFree)
	n += 8
	copy(res[n:], c.version)
	return res
}

func z2s(s string) string {
	return strings.TrimRight(s, "\x00")
}
