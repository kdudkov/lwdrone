package msg

import (
	"bufio"
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
	host        string
	port        int
	connTimeout time.Duration
	ch          chan *VideoFrame
	conn        net.Conn
	frames      int64
}

func StartStreamer(host string, port int, cmd *Command) (*Streamer, error) {
	s := &Streamer{
		host:        host,
		port:        port,
		connTimeout: time.Second * 5,
		ch:          make(chan *VideoFrame, 100),
	}

	var err error

	s.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.host, s.port), s.connTimeout)

	if err == nil {
		go s.start(cmd)
	} else {
		close(s.ch)
	}
	return s, err
}

func (s *Streamer) Stop() {

}

func (s *Streamer) start(startCmd *Command) {
	hbCmd := NewCommand(CmdHeartbeat, nil).ToByte()

	reader := bufio.NewReader(s.conn)
	writer := bufio.NewWriter(s.conn)

	defer close(s.ch)

	lastHb := time.Now()
	buf := make([]byte, hdrLen)

	_, err := writer.Write(startCmd.ToByte())
	if err != nil {
		println(err)
		return
	}

	for {
		cmd, err := ReadFrameWithBuf(reader, buf)
		if err != nil {
			return
		}

		if cmd.GetCode() == CmdHeartbeat {
			continue
		}

		if cmd.GetCode() == CmdRetreplayend {
			return
		}

		if len(cmd.body) == 0 {
			return
		}

		fu := NewUnmunger(cmd.GetStreamType(), cmd.GetDec1(), cmd.GetDec2())

		s.frames++
		s.ch <- NewVideoFrame(fu, cmd.body)
		if time.Now().Sub(lastHb) > time.Second {
			lastHb = time.Now()
			writer.Write(hbCmd)
		}
	}
}
