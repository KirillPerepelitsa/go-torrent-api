package torrent

import (
	"errors"
	"fmt"
	"reflect"
)

type file struct {
	info     info
	announce string
}

type info struct {
	pieceLength int
	length      int
	name        string
	pieces      []byte
}

func newFile(d map[string]interface{}) (f file, err error) {
	f = file{}

	if len(d) == 0 {
		return f, errors.New("empty map")
	}

	a, ok := d["announce"]
	if !ok {
		return f, errors.New("missed field `announce` in a .torrent file")
	}
	aStr, ok := a.(string)
	if !ok {
		return f, fmt.Errorf("wrong type `announce`, expected `string` received %s", reflect.TypeOf(a))
	}
	f.announce = aStr

	i, ok := d["info"]
	if !ok {
		return f, errors.New("missed field `info` in a .torrent file")
	}
	infoMap, ok := i.(map[string]interface{})
	if !ok {
		return f, fmt.Errorf("wrong type `info`, expected `map` received %s", reflect.TypeOf(i))
	}
	f.info = info{}

	pl, ok := infoMap["piece length"]
	if !ok {
		return f, errors.New("missed field `info.pieceLength` in a .torrent file")
	}
	plInt, ok := pl.(int)
	if !ok {
		return f, fmt.Errorf("wrong type `info.pieceLength`, expected `int` received %s", reflect.TypeOf(pl))
	}
	f.info.pieceLength = plInt

	p, ok := infoMap["pieces"]
	if !ok {
		return f, errors.New("missed field `info.pieces` in a .torrent file")
	}
	pStr, ok := p.(string)
	if !ok {
		return f, fmt.Errorf("wrong type `info.pieceLength`, expected `string` received %s", reflect.TypeOf(p))
	}
	f.info.pieces = []byte(pStr)

	n, ok := infoMap["name"]
	if !ok {
		return f, errors.New("missed field `info.name` in a .torrent file")
	}
	nStr, ok := n.(string)
	if !ok {
		return f, fmt.Errorf("wrong type `info.name`, expected `string` received %s", reflect.TypeOf(n))
	}
	f.info.name = nStr

	l, ok := infoMap["length"]
	if !ok {
		return f, errors.New("missed field `info.length` in a .torrent file")
	}
	lInt, ok := l.(int)
	if !ok {
		return f, fmt.Errorf("wrong type `info.length`, expected `int` received %s", reflect.TypeOf(l))
	}
	f.info.length = lInt

	return
}
