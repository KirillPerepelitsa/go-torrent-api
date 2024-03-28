package beencode

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalString(t *testing.T) {
	var buf bytes.Buffer
	s := "string"
	_, e := Marshal(&buf, s)
	assert.NoError(t, e)
	assert.Equal(t, fmt.Sprintf("%d:%s", len(s), s), buf.String())
}

func TestMarshalInt(t *testing.T) {
	var buf bytes.Buffer
	i := 120
	_, e := Marshal(&buf, i)
	assert.NoError(t, e)
	assert.Equal(t, fmt.Sprintf("i%de", i), buf.String())
}

func TestMarshalUint(t *testing.T) {
	var buf bytes.Buffer
	i := uint(120)
	_, e := Marshal(&buf, i)
	assert.NoError(t, e)
	assert.Equal(t, fmt.Sprintf("i%de", i), buf.String())
}

func TestMarshalFloatReturnError(t *testing.T) {
	var buf bytes.Buffer
	i := 1.5
	_, e := Marshal(&buf, i)
	assert.Equal(t, "", buf.String())
	assert.EqualError(t, errors.New("cannot marshal value of type float64"), e.Error())
}

func TestMarshalMap(t *testing.T) {
	var buf bytes.Buffer
	m := map[string]string{"key": "value", "key2": "value2"}
	_, e := Marshal(&buf, m)
	assert.NoError(t, e)
	assert.Equal(t, "d3:key5:value4:key26:value2e", buf.String())
}

func TestMarshalMapWithNonStringKey(t *testing.T) {
	var buf bytes.Buffer
	m := map[int]string{1: "value", 2: "value2"}
	_, e := Marshal(&buf, m)
	assert.EqualError(t, errors.New("map keys can be only string, int provided"), e.Error())
	assert.Equal(t, "", buf.String())
}

func TestMarshalList(t *testing.T) {
	var buf bytes.Buffer
	l := []string{"first", "second"}
	_, e := Marshal(&buf, l)
	assert.NoError(t, e)
	assert.Equal(t, "l5:first6:seconde", buf.String())
}

func TestMarshalSimpleStruct(t *testing.T) {
	var buf bytes.Buffer
	s := struct {
		key       string
		keySecond int
	}{
		key:       "string",
		keySecond: 120,
	}
	_, e := Marshal(&buf, s)
	assert.NoError(t, e)
	assert.Equal(t, "d3:key6:string9:keySecondi120ee", buf.String())
}

func TestMarshalCompositeStruct(t *testing.T) {
	var buf bytes.Buffer
	type Type struct {
		list []string
		key  string
	}
	type TypeOther struct {
		nested Type
		key    string
	}
	s := TypeOther{
		key: "string",
		nested: Type{
			list: []string{"one", "two"},
			key:  "string",
		},
	}
	_, e := Marshal(&buf, s)
	assert.NoError(t, e)
	assert.Equal(t, "d3:key6:string6:nestedd3:key6:string4:listl3:one3:twoeee", buf.String())
}
