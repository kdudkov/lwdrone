package msg

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type VideoFrame struct {
	flag   uint32
	size   int64
	count  int64
	gphoto uint32

	data []byte
}

func NewVideoFrame(u *VideoFrameUnmunger, data []byte) *VideoFrame {
	n := 0
	le := binary.LittleEndian
	v := &VideoFrame{}

	v.flag = le.Uint32(data[n:])
	n += 4
	v.size = int64(le.Uint32(data[n:]))
	n += 4
	v.count = int64(le.Uint64(data[n:]))
	n += 8
	v.gphoto = le.Uint32(data[n:])
	n += 4
	v.data = data[32:]
	if len(v.data) != int(v.size) {
		panic("!")
	}

	u.Unmunge(v.data, v.size, v.count)

	return v
}

type Streamer struct {
	host           string
	port           int
	connectTimeout time.Duration
	hbTimeout      time.Duration
	ch             chan *VideoFrame
	frames         int64
}

func StartStreamer(host string, port int, startCmd *Command, perm bool) *Streamer {
	s := &Streamer{
		host:           host,
		port:           port,
		connectTimeout: time.Second * 3,
		hbTimeout:      time.Second,
		ch:             make(chan *VideoFrame, 100),
	}

	go s.loop(startCmd, perm)

	return s
}

func (s *Streamer) loop(startCmd *Command, perm bool) {
	defer close(s.ch)

	hbCmd := NewCommand(CmdHeartbeat, nil).ToByte()
	buf := make([]byte, hdrLen)

	var conn net.Conn

	for {
		var err error
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.host, s.port), s.connectTimeout)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		lastHb := time.Now()

		_, err = conn.Write(startCmd.ToByte())
		if err != nil {
			_ = conn.Close()
			conn = nil
			break
		}

		for {
			cmd, err := ReadFrameWithBuf(conn, buf)
			if err != nil {
				break
			}

			if cmd.GetCode() == CmdHeartbeat {
				continue
			}

			if cmd.GetCode() == CmdRetreplayend {
				break
			}

			if len(cmd.body) == 0 {
				break
			}

			fu := NewUnmunger(cmd.GetStreamType(), cmd.GetDec1(), cmd.GetDec2())

			s.frames++
			s.ch <- NewVideoFrame(fu, cmd.body)
			if time.Now().Sub(lastHb) > s.hbTimeout {
				lastHb = time.Now()
				_, _ = conn.Write(hbCmd)
			}
		}

		_ = conn.Close()

		if !perm {
			return
		}
	}
}
