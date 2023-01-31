package main

import (
	"encoding/binary"
	"fmt"
)

const (
	heartbeat    = 1
	startstream  = 2 // stream
	stopstream   = 3 // stream
	settime      = 4
	gettime      = 5
	getrecplan   = 6
	getreclist   = 8
	startreplay  = 9  // stream
	stopreplay   = 16 // stream
	setrecplan   = 17
	getfile      = 18 // stream
	takepic      = 19
	delfile      = 20
	reformatsd   = 21
	setwifiname  = 22
	setwifipass  = 23
	setwifichan  = 24
	restartwifi  = 25
	setwifidefs  = 32
	getcamflip   = 33
	setcamflip   = 34
	getbaudrate  = 35
	setbaudrate  = 36
	getconfig    = 37
	setconfig    = 38
	getpiclist   = 39
	get1080p     = 40
	set1080p     = 41
	getpiclist2  = 42
	takepic2     = 43
	getrectime   = 48
	setrectime   = 49
	retstream    = 257
	retreplay    = 259
	retreplayend = 261
	retgetfile   = 262
)

var msg2text = map[uint32]string{
	1:   "heartbeat",
	2:   "startstream",
	3:   "stopstream",
	4:   "settime",
	5:   "gettime",
	6:   "getrecplan",
	8:   "getreclist",
	9:   "startreplay",
	16:  "stopreplay",
	17:  "setrecplan",
	18:  "getfile",
	19:  "takepic",
	20:  "delfile",
	21:  "reformatsd",
	22:  "setwifiname",
	23:  "setwifipass",
	24:  "setwifichan",
	25:  "restartwifi",
	32:  "setwifidefs",
	33:  "getcamflip",
	34:  "setcamflip",
	35:  "getbaudrate",
	36:  "setbaudrate",
	37:  "getconfig",
	38:  "setconfig",
	39:  "getpiclist",
	40:  "get1080p",
	41:  "set1080p",
	42:  "getpiclist2",
	43:  "takepic2",
	48:  "getrectime",
	49:  "setrectime",
	257: "retstream",
	259: "retreplay",
	261: "retreplayend",
	262: "retgetfile",
}

var magic = []byte("lewei_cmd\x00")

type Command struct {
	code   uint32
	header [8]uint32
	body   []byte
}

func NewCommand(code int, data []byte) *Command {
	c := &Command{code: uint32(code), body: data}
	c.header[2] = uint32(len(data))
	return c
}

func (c *Command) ToByte() []byte {
	res := make([]byte, len(magic)+4*9+len(c.body))
	copy(res, magic)
	le := binary.LittleEndian
	le.PutUint32(res[len(magic):], c.code)
	for i, v := range c.header {
		le.PutUint32(res[len(magic)+4+i*4:], v)
	}
	copy(res[len(magic)+4*9:], c.body)

	return res
}

func FromByte(data []byte) (*Command, error) {
	for i, b := range magic {
		if data[i] != b {
			return nil, fmt.Errorf("no magic")
		}
	}

	c := &Command{}
	le := binary.LittleEndian
	c.code = le.Uint32(data[len(magic):])
	for i, _ := range c.header {
		c.header[i] = le.Uint32(data[len(magic)+(i+1)*4:])
	}

	c.body = make([]byte, len(data)-len(magic)-4*9)
	copy(c.body, data[len(magic)+4*9:])

	return c, nil
}

func (c *Command) String() string {
	return fmt.Sprintf("cmd %s, data: %v", msg2text[c.code], c.body)
}
