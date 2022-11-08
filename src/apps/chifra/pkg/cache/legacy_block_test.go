package cache

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestReadLegacyBlock(t *testing.T) {
	f, err := os.Open("./legacy_block.bin")
	defer f.Close()
	buf := bufio.NewReader(f)
	if err != nil {
		t.Fatal("cannot open file")
	}
	_, err = ReadLegacyBlock(buf)
	if err != nil {
		t.Fatal(err)
	}

}

func Test_readCppString(t *testing.T) {
	source := bytes.NewBuffer([]byte{})
	err := binary.Write(source, binary.LittleEndian, uint64(6))
	if err != nil {
		t.Fatal("while preparing test data:", err)
	}
	err = binary.Write(source, binary.LittleEndian, []byte("CBlock"))
	if err != nil {
		t.Fatal("while preparing test data:", err)
	}

	result := cppString{}
	buf := bufio.NewReader(source)
	err = readCppString(buf, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.size != 6 {
		t.Fatal("wrong size:", result.size)
	}
	if string(result.content) != "CBlock" {
		t.Fatal("wrong content:", result.content)
	}
}

// func Test_ReadTransaction(t *testing.T) {
// 	f, err := os.Open("./legacy_block.bin")
// 	buf := bufio.NewReader(f)
// 	if err != nil {
// 		t.Fatal("cannot open file")
// 	}
// 	tr, err := ReadTransaction(buf)
// 	log.Println(tr)
// 	t.Fatal("stop")
// }
