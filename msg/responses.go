package msg

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Config struct {
	Ch         byte
	Flip       byte
	WifiSec    byte
	WifiName   string
	WifiPass   string
	Time       uint64
	SdcMounted byte
	SdcSize    uint64
	SdcFree    uint64
	Version    string
}

func ConfigFromBytes(data []byte) (*Config, error) {
	if len(data) != 140 {
		return nil, fmt.Errorf("invalid len: %d (must be 140)", len(data))
	}

	le := binary.LittleEndian
	c := &Config{}
	c.Ch = data[0]
	c.Flip = data[1]
	c.WifiSec = data[2]
	n := 3
	c.WifiName = z2s(string(data[n : n+32]))
	n += 32
	c.WifiPass = z2s(string(data[n : n+32]))
	n += 32
	c.Time = le.Uint64(data[n:])
	n += 8
	c.SdcMounted = data[n]
	n++
	c.SdcSize = le.Uint64(data[n:])
	n += 8
	c.SdcFree = le.Uint64(data[n:])
	n += 8
	c.Version = z2s(string(data[n:]))
	return c, nil
}

func (c *Config) ToBytes() []byte {
	res := make([]byte, 140)
	le := binary.LittleEndian
	res[0] = c.Ch
	res[1] = c.Flip
	res[2] = c.WifiSec
	n := 3
	copy(res[n:], c.WifiName)
	n += 32
	copy(res[n:], c.WifiPass)
	n += 32
	le.PutUint64(res[n:], c.Time)
	n += 8
	res[n] = c.SdcMounted
	n++
	le.PutUint64(res[n:], c.SdcSize)
	n += 8
	le.PutUint64(res[n:], c.SdcFree)
	n += 8
	copy(res[n:], c.Version)
	return res
}

func z2s(s string) string {
	return strings.TrimRight(s, "\x00")
}
