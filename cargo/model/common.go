package model

import (
	"time"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
)

type Common struct {
	ID      int64     `json:"-"`
	UUID    []byte    `json:"uuid"`
	Added   time.Time `json:"added"`
	Updated time.Time `json:"updated"`
	Flag    uint8     `json:"flag"`
}

func NewUUID() ([]byte, error) {
	timestamp := time.Now().UnixNano()
	nano := make([]byte, 8)
	binary.LittleEndian.PutUint64(nano, uint64(timestamp))

	data := make([]byte, 8)
	seed := make([]byte, 0)

	_, err := rand.Read(data)
	if err != nil {
			return nil, err
	}

	seed = append(seed, nano...)
	seed = append(seed, data...)

	uuid := sha256.Sum256(seed)
	return uuid[:], nil
}

func Now() time.Time {
	return time.Now().UTC()
}

func NewCommon() (Common, error) {
	var self Common

	uuid, err := NewUUID()
	if err != nil {
		return self, err
	}

	self = Common{
		UUID:    uuid,
		Flag:    0,
	}

	return self, nil
}
