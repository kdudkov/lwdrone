package msg

import (
	"encoding/binary"
	"fmt"
	"io"
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
}

func StartStreamer(host string, port int, cmd *Command) (*Streamer, error) {
	s := &Streamer{
		host:        host,
		port:        port,
		connTimeout: time.Second * 15,
		ch:          make(chan *VideoFrame, 100),
	}

	var err error

	s.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.host, s.port), s.connTimeout)

	if err == nil {
		go s.start(cmd)
	}
	return s, err
}

func (s *Streamer) Stop() {

}

func (s *Streamer) start(cmd *Command) {
	hbCmd := NewCommand(CmdHeartbeat, nil).ToByte()

	lastHb := time.Now()
	buf := make([]byte, hdrLen)

	s.conn.Write(cmd.ToByte())

	for {
		_, err := io.ReadFull(s.conn, buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		cmdOut, err := FromByte(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		var dataBuf []byte

		if cmdOut.GetSize() > 0 {
			dataBuf = make([]byte, cmdOut.GetSize())
			_, err := io.ReadFull(s.conn, dataBuf)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if cmdOut.GetCode() == CmdHeartbeat {
			continue
		}

		if cmdOut.GetCode() == CmdRetreplayend {
			return
		}

		if len(dataBuf) == 0 {
			return
		}

		fu := NewUnmunger(cmdOut.GetStreamType(), cmdOut.GetDec1(), cmdOut.GetDec2())

		s.ch <- NewVideoFrame(fu, dataBuf)
		if time.Now().Sub(lastHb) > time.Second {
			lastHb = time.Now()
			s.conn.Write(hbCmd)
		}
	}

}
