package main

import (
	"fmt"
	"net"
	"time"
)

type Lwdrone struct {
	host       string
	cmdPort    int
	streamPort int
	timeout    time.Duration
}

func NewDrone() *Lwdrone {
	return &Lwdrone{
		host:       "192.168.0.1",
		cmdPort:    8060,
		streamPort: 9060,
		timeout:    time.Second * 2,
	}
}

func (l *Lwdrone) GetConfig() (c *Config, err error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", l.host, l.cmdPort), l.timeout)

	if err != nil {
		return
	}

	cmd := NewCommand(getconfig, nil)

	_, err = conn.Write(cmd.ToByte())
	if err != nil {
		return
	}

	buf := make([]byte, 165535)

	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	_, err = FromByte(buf[:n])
	if err != nil {
		return
	}

	n, err = conn.Read(buf)
	if err != nil {
		return
	}

	c, err = ConfigFromBytes(buf[:n])
	return
}

func main() {
	c, err := NewDrone().GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("version: %s\n", c.version)
	fmt.Printf("flash mounted: %d\n", c.sdcMounted)
	fmt.Printf("flash size: %d MiB\n", c.sdcSize/1024/1024)
	fmt.Printf("flash free: %d MiB (%.d%%)\n", c.sdcFree/1024/1024, 100.*c.sdcFree/c.sdcSize)
}
