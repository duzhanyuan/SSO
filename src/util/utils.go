package util

import (
	"crypto/rand"
	"crypto/rc4"
	//"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"hash/fnv"
	"os"
	"time"
)

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func GenRandomBytes(len int) []byte {
	rb := make([]byte, len)
	rand.Read(rb)
	return rb
}
func GenTimestamp(key []byte) string {
	timestamp := time.Now().Unix()
	timestamp_bytes := Int64ToBytes(timestamp)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(timestamp_bytes, timestamp_bytes)
	timestamp_str := hex.EncodeToString(timestamp_bytes)
	return timestamp_str
}

func CheckTimestamp(timestamp_str string, key []byte) bool {
	c, _ := rc4.NewCipher(key)
	timestamp_bytes, _ := hex.DecodeString(timestamp_str)
	c.XORKeyStream(timestamp_bytes, timestamp_bytes)
	if len(timestamp_bytes) != 8 {
		return false
	}
	timestamp := BytesToInt64(timestamp_bytes)
	current := time.Now().Unix()
	if current > timestamp+60*5 || current+60*5 < timestamp {
		return false
	}
	return true
}

func EncryptString(str string, key []byte) string {
	bytes := []byte(str)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(bytes, bytes)
	return hex.EncodeToString(bytes)
}

func DecryptString(str string, key []byte) string {
	bytes, _ := hex.DecodeString(str)
	c, _ := rc4.NewCipher(key)
	c.XORKeyStream(bytes, bytes)
	return string(bytes)
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func IsDevelopment() bool {
	return os.Getenv("SSO_ENV") == "dev"
}
