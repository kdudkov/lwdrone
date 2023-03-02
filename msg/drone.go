package msg

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const delta = -time.Hour * 5

const (
	CamUp = iota
	CamUpMirror
	CamDownMirror
	CamDown
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
		streamPort: 7060,
		timeout:    time.Second * 2,
	}
}

func (l *Lwdrone) sendCommand(cmd *Command) (*Command, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", l.host, l.cmdPort), l.timeout)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 200))
	_, err = conn.Write(cmd.ToByte())
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

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
	binary.LittleEndian.PutUint64(data, uint64(time.Now().Add(delta).Unix()))
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

func (l *Lwdrone) StartStream(hires bool, fl *os.File) error {
	cmd := NewCommand(CmdStartstream, nil)
	if hires {
		cmd.SetArg(1)
	}

	s, err := StartStreamer(l.host, l.streamPort, cmd)
	if err != nil {
		return err
	}

	for b := range s.ch {
		fl.Write(b.data)
	}

	return nil
}
