package beencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func Unmarshal(reader io.Reader) (interface{}, error) {
	r, ok := reader.(*bufio.Reader)
	if !ok {
		r = bufio.NewReader(reader)
	}
	return decode(r)
}

func decode(r *bufio.Reader) (interface{}, error) {
	char, err := r.ReadByte()
	//s := r.Size()
	//buf := make([]byte, s)
	//_, err = r.Read(buf)
	//println(err)
	//println(string(buf))
	//return nil, errors.New("FOO")
	if err != nil {
		return nil, err
	}
	switch {
	case char == 'l':
		var l []interface{}
		for {
			b, e := r.ReadByte()
			if e != nil {
				return nil, e
			}
			if b == 'e' {
				return l, nil
			}
			err := r.UnreadByte()
			if err != nil {
				return nil, err
			}
			item, err := decode(r)
			if err != nil {
				return nil, err
			}
			l = append(l, item)
		}
	case char == 'i':
		i, err := readUntil(r, 'e')
		if err != nil {
			return nil, err
		}
		return strconv.Atoi(string(i))
	case char == 'd':
		d := make(map[string]interface{})
		for {
			c, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if c == 'e' {
				return d, nil
			}
			err = r.UnreadByte()
			if err != nil {
				return nil, err
			}
			k, err := decode(r)
			if err != nil {
				return nil, err
			}
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("cant cast dictionary key to string, %s", k)
			}
			v, err := decode(r)
			d[key] = v
		}
	case char >= '0' && char <= '9':
		err := r.UnreadByte()
		if err != nil {
			return nil, err
		}
		var strLen int64
		if s, err := readUntil(r, ':'); err == nil {
			strLen, err = strconv.ParseInt(string(s), 10, 64)
		}
		if err != nil {
			return nil, err
		}

		// Can we peek that much data out of r?
		if peekBuf, peekErr := r.Peek(int(strLen)); peekErr == nil {
			_, err = r.Discard(int(strLen))
			if err != nil {
				return nil, err
			}
			return string(peekBuf), nil
		}

		buf := make([]byte, strLen)
		_, err = readFull(r, buf)
		if err != nil {
			return nil, err
		}
		return string(buf), nil

	default:
		return r, fmt.Errorf("wrong starting character %s", string(char))
	}
}

func readFull(r *bufio.Reader, buf []byte) (n int, err error) {
	return readAtLeast(r, buf, len(buf))
}

func readAtLeast(r *bufio.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func readUntil(r *bufio.Reader, delim byte) (data []byte, err error) {
	data, err = r.ReadSlice(delim)
	if err != nil {
		return nil, err
	}
	lenData := len(data)
	if lenData > 0 {
		data = data[:lenData-1]
	} else {
		err = errors.New("bad r.ReadSlice() length")
	}
	return
}

func strLen(r *bufio.Reader) (num int, l int, err error) {
	var buf []byte
	for {
		char, err := r.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		switch char {
		case ':':
			l = len(string(buf))
			num, err = strconv.Atoi(string(buf))
			if err != nil {
				return 0, 0, err
			}
			return num, l, nil
		default:
			buf = append(buf, char)
		}
	}
}
