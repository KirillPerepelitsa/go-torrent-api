package torrent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileNewOk(t *testing.T) {
	f, err := newFile(getMap())
	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(file{announce: "test-announce", info: info{
		pieceLength: 262144,
		length:      1,
		name:        "name-test",
		pieces:      []byte("hex>2F 41 1E 3E DC 1E EC 28 A7 96 7C 65 CC F1 2A 24 0E 46 36 2C 1E 8E 76 39 E3 FA 32 89 43 2A A4 A7 1C D5 8E 1F B3 B2 29 B0 6C 4B C6 62 F9</hex>"),
	}}, f)
}

func TestFileNewWithMissedAnnounce(t *testing.T) {
	m := getMap()
	delete(m, "announce")
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `announce` in a .torrent file")
}

func TestFileNewWithWrongAnnounce(t *testing.T) {
	m := getMap()
	m["announce"] = 1
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `announce`, expected `string` received int")
}

func TestFileNewWithEmptyMap(t *testing.T) {
	m := map[string]interface{}{}
	_, err := newFile(m)
	assert.EqualError(t, err, "empty map")
}

func TestFileNewWithMissedInfo(t *testing.T) {
	m := map[string]interface{}{"announce": "test-announce"}
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `info` in a .torrent file")
}

func TestFileNewWithWrongInfo(t *testing.T) {
	m := getMap()
	m["info"] = "wrong"
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `info`, expected `map` received string")
}

func TestFileNewWithMissedInfoLength(t *testing.T) {
	m := getMap()
	delete(m["info"].(map[string]interface{}), "length")
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `info.length` in a .torrent file")
}

func TestFileNewWithWrongInfoLength(t *testing.T) {
	m := getMap()
	m["info"].(map[string]interface{})["length"] = "1"
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `info.length`, expected `int` received string")
}

func TestFileNewWithMissedInfoName(t *testing.T) {
	m := getMap()
	delete(m["info"].(map[string]interface{}), "name")
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `info.name` in a .torrent file")
}

func TestFileNewWithWrongInfoName(t *testing.T) {
	m := getMap()
	m["info"].(map[string]interface{})["name"] = 1
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `info.name`, expected `string` received int")
}

func TestFileNewWithMissedInfoPieceLength(t *testing.T) {
	m := getMap()
	delete(m["info"].(map[string]interface{}), "piece length")
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `info.pieceLength` in a .torrent file")
}

func TestFileNewWithWrongPieceLength(t *testing.T) {
	m := getMap()
	m["info"].(map[string]interface{})["piece length"] = "262144"
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `info.pieceLength`, expected `int` received string")
}

func TestFileNewWithMissedPieces(t *testing.T) {
	m := getMap()
	delete(m["info"].(map[string]interface{}), "pieces")
	_, err := newFile(m)
	assert.EqualError(t, err, "missed field `info.pieces` in a .torrent file")
}

func TestFileNewWithWrongPieces(t *testing.T) {
	m := getMap()
	m["info"].(map[string]interface{})["pieces"] = 1
	_, err := newFile(m)
	assert.EqualError(t, err, "wrong type `info.pieceLength`, expected `string` received int")
}

func getMap() map[string]interface{} {
	return map[string]interface{}{"announce": "test-announce", "info": map[string]interface{}{
		"length":       1,
		"name":         "name-test",
		"piece length": 262144,
		"pieces":       "<hex>2F 41 1E 3E DC 1E EC 28 A7 96 7C 65 CC F1 2A 24 0E 46 36 2C 1E 8E 76 39 E3 FA 32 89 43 2A A4 A7 1C D5 8E 1F B3 B2 29 B0 6C 4B C6 62 F9</hex>",
	}}
}
