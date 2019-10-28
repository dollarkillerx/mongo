/**
 * @Author: DollarKiller
 * @Description: utils
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 08:54 2019-10-28
 */
package mongo

import (
	"crypto/sha1"
	"encoding/hex"
)

// 获取sha1
func Sha1Encode(str string) string {
	data := []byte(str)
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}
