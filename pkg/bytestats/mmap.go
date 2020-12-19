package bytestats

import (
	"os"

	"github.com/alexeymaximov/go-bio/mmap"
)

func memmap(filename string) []byte {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	size := int(stat.Size())

	buf := make([]byte, size)
	readBytes := 0
	open := func() (*mmap.Mapping, error) {
		return mmap.OpenFile(filename, os.FileMode(0600), uintptr(size), 0, func(m *mmap.Mapping) error {
			readBytes++
			_, err := m.WriteAt(buf, 0)
			return err
		})
	}
	m, err := open()
	if err != nil {
		panic(err)
	}
	return m.Memory()
}
