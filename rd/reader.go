package rd

import (
	"fmt"
	"os"
)

type Reader struct {
	file *os.File
}

func Get(f *os.File) *Reader {
	rd := Reader{
		file: f,
	}
	return &rd
}

func (r *Reader) Read() (byte, error) {
	bf := []byte{}
	bf = make([]byte, 1)
	_, err := r.file.Read(bf)
	if err != nil {
		fmt.Printf("Error reading a single byte from the file\n")
		return 0, err
	}
	return bf[0], nil
}

func (r *Reader) ReadBuffer(cp int) ([]byte, int, error) {
	bf := []byte{}
	bf = make([]byte, cp)
	count, err := r.file.Read(bf)
	if err != nil {
		fmt.Printf("Error in reading a buffer from the file\n")
		return []byte{}, 0, err
	}
	return bf, count, nil
}
