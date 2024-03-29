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

func TestMarshalDictWithDashFields(t *testing.T) {
	r := strings.NewReader("d11:string-dashi120ee")
	res, e := Unmarshal(r)
	assert.NoError(t, e)
	assert.Equal(t, map[string]interface{}{"string-dash": 120}, res)
}

func TestMarshalLongString(t *testing.T) {
	r := strings.NewReader("146:foobarfoobarfoobarfoobarfoobarfoobar-foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar-foobarfoobarfoobarfoobarfoobarfoobar")
	res, e := Unmarshal(r)
	assert.NoError(t, e)
	assert.Equal(t, "foobarfoobarfoobarfoobarfoobarfoobar-foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar-foobarfoobarfoobarfoobarfoobarfoobar", res)

	//reader := bufio.NewReader(r)
	//buf := make([]byte, 195)
	//reader.Read(buf)
	//println(string(buf))
}
