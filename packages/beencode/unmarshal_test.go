package beencode

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnmarshalString(t *testing.T) {
	r := strings.NewReader("3:foo")
	res, _ := Unmarshal(r)
	assert.Equal(t, "foo", res)
}

func TestUnmarshalInt(t *testing.T) {
	r := strings.NewReader("i120e")
	res, _ := Unmarshal(r)
	assert.Equal(t, 120, res)
}

func TestUnmarshalList(t *testing.T) {
	r := strings.NewReader("l3:foo3:bare")
	res, _ := Unmarshal(r)
	assert.Equal(t, []interface{}{"foo", "bar"}, res)
}

func TestUnmarshalDictionary(t *testing.T) {
	r := strings.NewReader("d7:astring3:str4:binti120e5:clistl4:listee")
	res, _ := Unmarshal(r)
	assert.Equal(t, map[string]interface{}{"astring": "str", "bint": 120, "clist": []interface{}{"list"}}, res)
}
