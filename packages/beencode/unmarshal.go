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
		length, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, err
		}
		// discarding 2 bytes, length + :, for e.g 3:foo to foo
		_, err = r.Discard(2)
		if err != nil {
			return nil, err
		}
		buf := make([]byte, length)
		_, err = r.Read(buf)
		if err != nil {
			return nil, err
		}

		return string(buf), nil
	}
}
