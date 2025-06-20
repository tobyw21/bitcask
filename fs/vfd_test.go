package fs

import (
	"strconv"
	"testing"
)

func TestVfdOpenClose(t *testing.T) {

	vfdmgr := NewVfdMgr()

	vid := make([]int8, 0, 6)
	for i := range 6 {
		filePath := "../examples/" + strconv.Itoa(i+1)

		t.Logf(`Opening: %s`, filePath)
		vfdid, err := vfdmgr.VfdOpen(filePath)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		vid = append(vid, vfdid)

		t.Logf(`Opened File: %s, with Vfd ID: %d`, filePath, vfdid)

	}

	t.Logf(`vfd id list: %v`, vid)

	// here vfdid = 1 will be the least used after push front
	for _, i := range vid[1:] {
		t.Logf(`Closing Vfd ID: %d`, i)
		err := vfdmgr.VfdClose(i)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	}

}

func TestVfdRWNormal(t *testing.T) {

	vfdmgr := NewVfdMgr()

	vfdid, err := vfdmgr.VfdOpen("../examples/rw")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// data := []byte("Hello world")
	data := []byte{0x00, 0x01, 0xFE, 0xFF}

	vfdmgr.VfdWrite(vfdid, data, 0)

	buffer := make([]byte, 12)

	nread, err := vfdmgr.VfdRead(vfdid, &buffer, 0)
	t.Logf(`read to buffer %d, %v`, nread, string(buffer))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = vfdmgr.VfdClose(vfdid)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}
