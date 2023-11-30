package hashidv2

import (
	"github.com/speps/go-hashids/v2"

	"github.com/liushuangls/go-server-template/pkg/ecode"
)

const alphabet = "QP1WO2EI3RU4TY5AL6SK7DJ8FH9GZ0MXNCBV"
const minLength = 8

type Config struct {
	Type   int    `yaml:"Type"`
	Prefix string `yaml:"Prefix"`
	Salt   string `yaml:"Salt"`
}

type HashID struct {
	prefix string
	typ    int
	salt   string
	hashID *hashids.HashID
}

func New(c *Config) *HashID {
	hashID := &HashID{
		typ:    c.Type,
		prefix: c.Prefix,
		salt:   c.Salt,
	}
	hdata := hashids.NewData()
	hdata.Alphabet = alphabet
	hdata.MinLength = minLength
	hdata.Salt = hashID.salt
	hash, err := hashids.NewWithData(hdata)
	if err != nil {
		panic(err)
	}
	hashID.hashID = hash
	return hashID
}

func (h *HashID) Decode(idStr string) (int, error) {
	if idStr == "" || idStr == "0" {
		return 0, nil
	}
	if len(idStr) <= len(h.prefix) {
		return 0, ecode.InvalidHashID
	}
	if h.prefix != idStr[0:len(h.prefix)] {
		return 0, ecode.InvalidHashID
	}
	idStr = idStr[len(h.prefix):]
	ids, err := h.hashID.DecodeWithError(idStr)
	if err != nil {
		return 0, err
	}
	if len(ids) != 2 && ids[0] != h.typ {
		return 0, ecode.InvalidHashID
	}
	return ids[1], nil
}

func (h *HashID) DecodeNotE(idStr string) int {
	n, _ := h.Decode(idStr)
	return n
}

func (h *HashID) Encode(id int) (string, error) {
	if id == 0 {
		return "0", nil
	}
	idStr, err := h.hashID.Encode([]int{h.typ, id})
	if err != nil {
		return "", err
	}
	return h.prefix + idStr, nil
}

func (h *HashID) EncodeNotE(id int) string {
	s, _ := h.Encode(id)
	return s
}
