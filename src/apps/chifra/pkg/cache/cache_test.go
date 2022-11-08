package cache

import (
	"bytes"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type TestSerializer struct {
	Addr common.Address
	Apps uint64
}

// func (ts *TestSerializer) Serialize(b *bytes.Buffer) (err error) {
// 	_, err = b.Write(ts.Addr.Bytes())
// 	if err != nil {
// 		return
// 	}
// 	_, err = b.Write(ts.Apps)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

// func (ts *TestSerializer) Deserialize(r io.Reader) (err error) {
// 	err = binary.Read(r, binary.LittleEndian, ts)
// }

func Test_setHandler(t *testing.T) {
	value := &TestSerializer{
		Addr: common.HexToAddress("0xf503017D7baF7FBC0fff7492b751025c6A78179b"),
		Apps: 20,
	}
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	saveFn := func(filePath string, content io.Reader) (err error) {
		_, err = io.Copy(w, content)
		return
	}
	err := setHandler(saveFn, ItemABI, "file", value)
	if err != nil {
		t.Fatal("error while saving:", err)
	}

	// Now try to read the data
	read := &TestSerializer{}
	loadFn := func(_ string) (io.Reader, error) {
		return w, nil
	}
	err = getHandler(loadFn, ItemABI, "", read)
	if err != nil {
		t.Fatal("error while reading", err)
	}

	if read.Addr.Hex() != value.Addr.Hex() {
		t.Fatal("value mismatch Addr", read.Addr.Hex(), value.Addr.Hex())
	}

	if read.Apps != value.Apps {
		t.Fatal("value mismatch Apps", read.Apps, value.Apps)
	}
}
