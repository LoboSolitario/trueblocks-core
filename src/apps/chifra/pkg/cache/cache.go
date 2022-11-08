package cache

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
)

var cacheByteOrder = binary.LittleEndian

func save(filePath string, content io.Reader) (err error) {
	cacheDir := config.GetRootConfig().Settings.CachePath
	fullPath := path.Join(cacheDir, filePath)
	file, err := os.Create(fullPath)
	if err != nil {
		return
	}

	_, err = io.Copy(file, content)
	return
}

func load(filePath string) (file io.Reader, err error) {
	cacheDir := config.GetRootConfig().Settings.CachePath
	fullPath := path.Join(cacheDir, filePath)
	file, err = os.Open(fullPath)
	return
}

type Serializer interface {
	Serialize(b *bytes.Buffer) error
	Deserialize(r io.Reader) error
}

// TODO: find better name
type writeFn = func(filePath string, content io.Reader) error

func setHandler(writeHandler writeFn, itemType CacheItem, fileName string, value any) error {
	cachePath := itemToDirectory[itemType]
	filePath := path.Join(cachePath, fileName)
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, cacheByteOrder, value)
	if err != nil {
		return err
	}
	return writeHandler(filePath, buf)
}

func Set(itemType CacheItem, fileName string, value any) error {
	return setHandler(save, itemType, fileName, value)
}

type readFn = func(filePath string) (io.Reader, error)

func getHandler[Value any](readHandler readFn, itemType CacheItem, fileName string, value *Value) error {
	cachePath := itemToDirectory[itemType]
	filePath := path.Join(cachePath, fileName)
	reader, err := readHandler(filePath)
	if err != nil {
		return err
	}
	return binary.Read(reader, cacheByteOrder, value)
}

func Get[Value any](itemType CacheItem, fileName string, value *Value) error {
	return getHandler(load, itemType, fileName, value)
}
