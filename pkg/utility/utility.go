package utility

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
	"runtime"
)

const (
	HASH_SIZE = sha256.Size
)

var (
	_ORDER = binary.LittleEndian
)

var (
	HASH_ZERO = EncodeToString(_HASH_ZERO[:])
	HASH_NULL = EncodeToString(_HASH_NULL[:])
)

func Serialize(e any) []byte {
	data, err := json.MarshalIndent(e, "", "\t")
	LogError(err)
	return data
}

func Deserialize(e any, data []byte) {
	LogError(json.Unmarshal(data, e))
}

func EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeString(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	LogError(err)
	return data
}

func HashSum(e any) string {
	hash := sha256.Sum256(Serialize(e))
	return EncodeToString(hash[:])
}

func Load(name string, e any) error {
	data, err := os.ReadFile(name)
	if err != nil {
		LogError(err)
		return err
	}
	Deserialize(e, data)
	return nil
}

func Save(name string, e any) {
	LogError(os.WriteFile(name, Serialize(e), 0633))
}

func LogError(err error) error {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("error: %s: %d: %s", file, line, err.Error())
	}
	return err
}

var (
	_HASH_ZERO = [HASH_SIZE]byte{}
	_HASH_NULL = sha256.Sum256(nil)
)
