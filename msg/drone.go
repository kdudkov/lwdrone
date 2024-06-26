package msg

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	CamUp = iota
	CamUpMirror
	CamDownMirror
	CamDown
)

type Lwdrone struct {
	host        string
	cmdPort     int
	streamPort  int
	dialTimeout time.Duration
	connTimeout time.Duration
}

func NewDrone() *Lwdrone {
	return &Lwdrone{
		host:        "192.168.0.1",
		cmdPort:     8060,
		streamPort:  7060,
		dialTimeout: time.Second * 2,
		connTimeout: time.Second * 2,
	}
}

func (l *Lwdrone) sendCommand(cmd *Command) (*Command, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", l.host, l.cmdPort), l.dialTimeout)

	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(l.connTimeout))

	_, err = conn.Write(cmd.ToByte())
	if err != nil {
		return nil, err
	}

	return ReadFrame(conn)
}

func (l *Lwdrone) GetConfig() (*Config, error) {
	cmd := NewCommand(CmdGetconfig, nil)

	res, err := l.sendCommand(cmd)
	if err != nil {
		return nil, err
	}
	return ConfigFromBytes(res.GetBody())
}

func (l *Lwdrone) TakePicture() (*Picture, error) {
	cmd := NewCommand(CmdTakepic, nil)

	res, err := l.sendCommand(cmd)
	if err != nil {
		return nil, err
	}

	p, err := PictureFromBytes(res.body)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (l *Lwdrone) TakePicture2(save bool) (*Picture, error) {
	cmd := NewCommand(CmdTakepic2, nil)
	if save {
		cmd.SetArg(1)
	}
	res, err := l.sendCommand(cmd)
	if err != nil {
		return nil, err
	}

	p, err := PictureFromBytes(res.body)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (l *Lwdrone) GetPicList() error {
	cmd := NewCommand(CmdGetpiclist2, nil)
	cmd.SetArg(100)

	res, err := l.sendCommand(cmd)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func (l *Lwdrone) GetTime() (time.Time, error) {
	cmd := NewCommand(CmdGettime, nil)

	res, err := l.sendCommand(cmd)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return time.UnixMilli(int64(binary.LittleEndian.Uint64(res.body))), nil
}

func (l *Lwdrone) SetTime() error {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(time.Now().Unix()))
	cmd := NewCommand(CmdSettime, data)
	res, err := l.sendCommand(cmd)
	if err != nil {
		return err
	}
	if res.GetArg() != 0 {
		return fmt.Errorf("arg")
	}
	return nil
}

func (l *Lwdrone) SetCamFlip(cam int) error {
	cmd := NewCommand(CmdSetcamflip, nil)
	cmd.SetArg(cam)
	res, err := l.sendCommand(cmd)
	if err != nil {
		return err
	}
	if res.GetArg() != 0 {
		return fmt.Errorf("arg")
	}
	return nil
}

func (l *Lwdrone) GetCamFlip(cam int) (int, error) {
	cmd := NewCommand(CmdGetcamflip, nil)
	res, err := l.sendCommand(cmd)
	if err != nil {
		return 0, err
	}
	if res.GetArg() != 0 {
		return 0, fmt.Errorf("arg")
	}
	return res.GetArg(), nil
}

func (l *Lwdrone) StartStream(hires bool, fl io.Writer, perm bool) {
	cmd := NewCommand(CmdStartstream, nil)
	if hires {
		cmd.SetArg(1)
	}

	s := StartStreamer(l.host, l.streamPort, cmd, perm)

	for b := range s.ch {
		fl.Write(b.data)
	}
}
