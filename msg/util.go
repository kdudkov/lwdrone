package msg

import (
	"fmt"
	"io"
)

func ReadFrame(reader io.Reader) (*Command, error) {
	buf := make([]byte, hdrLen)
	return ReadFrameWithBuf(reader, buf)
}

func ReadFrameWithBuf(reader io.Reader, hdrBuf []byte) (*Command, error) {
	if len(hdrBuf) != hdrLen {
		return nil, fmt.Errorf("invalid buffer")
	}

	_, err := io.ReadFull(reader, hdrBuf)
	if err != nil {
		return nil, err
	}

	cmd, err := FromByte(hdrBuf)
	if err != nil {
		return nil, err
	}

	if cmd.GetSize() > 0 {
		cmd.body = make([]byte, cmd.GetSize())
		_, err := io.ReadFull(reader, cmd.body)
		if err != nil {
			return cmd, err
		}
	}

	return cmd, nil
}
