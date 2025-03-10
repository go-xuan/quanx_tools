package encrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	fmt.Println(doCrypto("base64(hello)"))
	fmt.Println(doCrypto("md5(base64(hello))"))
}
