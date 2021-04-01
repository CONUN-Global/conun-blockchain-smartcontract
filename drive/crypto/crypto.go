package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

/**
EncodeToSha256
@param String
*/
func EncodeToSha256(ipfsHash string) string {

	hash := sha1.New()
	hash.Write([]byte(ipfsHash))
	return "con" + hex.EncodeToString(hash.Sum(nil))

}
