package util

import (
	"hash/fnv"
	"os"
)

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func IsDevelopment() bool {
	return os.Getenv("SSO_ENV") == "dev"
}
