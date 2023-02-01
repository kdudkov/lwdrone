package msg

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Ch         byte
	Flip       byte
	WifiSec    byte
	WifiName   string
	WifiPass   string
	Time       time.Time
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
	c.Time = time.Unix(int64(le.Uint64(data[n:])), 0)
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
	le.PutUint64(res[n:], uint64(c.Time.Unix()))
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

type Picture struct {
	Size int
	Time time.Time
	X    int
	Path string

	Data []byte
}

func (p *Picture) String() string {
	return fmt.Sprintf("Picture %s %d %s", p.Path, p.Size, p.Data)
}

func PictureFromBytes(data []byte) (*Picture, error) {
	if len(data) < 128 {
		return nil, fmt.Errorf("too small data")
	}
	le := binary.LittleEndian
	n := 0
	p := &Picture{}
	p.Size = int(le.Uint32(data[n:]))
	n += 4
	p.Time = time.UnixMilli(int64(le.Uint32(data[n:])))
	n += 4
	p.X = int(le.Uint32(data[n:]))
	n += 4
	p.Path = z2s(string(data[n : n+100]))
	n += 100
	p.Data = make([]byte, len(data)-128)
	copy(p.Data, data[128:])

	return p, nil
}

func z2s(s string) string {
	return strings.TrimRight(s, "\x00")
}
