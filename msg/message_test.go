package msg

import (
	"encoding/hex"
	"fmt"
	"testing"
)

var smsg = []string{
	"6c657765695f636d640004000000000000000000000008000000000000000000000000000000000000000000000006afd76300000000",
	"6c657765695f636d6400060000000000000000000000000000000000000000000000000000000000000000000000",
	"6c657765695f636d6400250000000000000000000000000000000000000000000000000000000000000000000000",
}

var ans = "2403000000000000000000000000000000000000000000000000000000000000000000313233343536373800000000000000000000000000000000000000000000000008afd7630000000001000000ea0e0000000000be780e000000563330300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

func Test1(t *testing.T) {
	for _, s := range smsg {
		msg1 := str2bin(s)

		c, err := FromByte(msg1)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(c)
	}
}

func Test2(t *testing.T) {
	msg1 := str2bin(smsg[0])

	c, err := FromByte(msg1)
	if err != nil {
		t.Error(err)
	}

	msg2 := c.ToByte()

	for i, b := range msg1 {
		if msg2[i] != b {
			t.Error()
		}
	}
}

func Test3(t *testing.T) {
	c, err := ConfigFromBytes(str2bin(ans))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(c)
}

func Test4(t *testing.T) {
	msg1 := str2bin(ans)
	c, err := ConfigFromBytes(msg1)
	if err != nil {
		t.Error(err)
	}
	msg2 := c.ToBytes()

	for i, b := range msg1 {
		if msg2[i] != b {
			t.Errorf("pos %d %d != %d", i, b, msg2[i])
		}
	}
}

func str2bin(s string) []byte {
	res, _ := hex.DecodeString(s)
	return res
}
