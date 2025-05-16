package encrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	fmt.Println(doEncrypt("base64(hello)"))
	fmt.Println(doEncrypt("md5(base64(hello))"))
}
