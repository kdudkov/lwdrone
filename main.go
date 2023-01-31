package main

import (
	"fmt"
	"net"
	"time"

	"github.com/kdudkov/lwdrone/msg"
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

func (l *Lwdrone) cmd(cmd *msg.Command) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", l.host, l.cmdPort), l.timeout)

	if err != nil {
		return nil, err
	}

	conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 200))
	_, err = conn.Write(cmd.ToByte())
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	buf := make([]byte, 165535)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	c2, err := msg.FromByte(buf[:n])
	if err != nil {
		return nil, err
	}

	fmt.Println(cmd.GetCode(), " -> ", c2.GetCode(), c2.GetSize())
	if c2.GetSize() > 0 && len(c2.GetBody()) == 0 {
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		n, err = conn.Read(buf)
		if err != nil {
			return nil, err
		}
		if n != c2.GetSize() {
			return nil, fmt.Errorf("invalid size: %d != %d", n, c2.GetSize())
		}

		return buf[:n], nil
	} else {
		return c2.GetBody(), nil
	}
}

func (l *Lwdrone) GetConfig() (*msg.Config, error) {
	cmd := msg.NewCommand(msg.CmdGetconfig, nil)

	res, err := l.cmd(cmd)
	if err != nil {
		return nil, err
	}
	return msg.ConfigFromBytes(res)
}

func (l *Lwdrone) GetRecPlan() error {
	cmd := msg.NewCommand(msg.CmdGetrecplan, nil)

	res, err := l.cmd(cmd)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func (l *Lwdrone) GetRecList() error {
	cmd := msg.NewCommand(msg.CmdGetreclist, nil)

	res, err := l.cmd(cmd)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func main() {
	drone := NewDrone()
	c, err := drone.GetConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("version: %s\n", c.Version)
	fmt.Printf("flash mounted: %d\n", c.SdcMounted)
	fmt.Printf("flash size: %d MiB\n", c.SdcSize/1024/1024)
	fmt.Printf("flash free: %d MiB (%.d%%)\n", c.SdcFree/1024/1024, 100.*c.SdcFree/c.SdcSize)

	drone.GetRecList()
}
