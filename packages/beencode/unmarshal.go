package beencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
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
	if err != nil {
		return nil, err
	}
	switch char {
	case 'l':
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
	case 'i':
		i, err := r.ReadBytes('e')
		if err != nil {
			return nil, err
		}
		return strconv.Atoi(strings.Trim(string(i), "e"))
	case 'd':
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
	default:
		err := r.UnreadByte()
		if err != nil {
			return nil, err
		}

		num, _, err := strLen(r)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		buf := make([]byte, num)
		_, err = r.Read(buf)
		if err != nil {
			return nil, err
		}
		return string(buf), nil
	}
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
