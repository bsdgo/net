package factory

import (
	"github.com/skycoin/skycoin/src/cipher"
)

func GenRegMsg() []byte {
	result := make([]byte, MSG_HEADER_END)
	result[MSG_OP_BEGIN] = OP_REG
	return result
}

func GenRegRespMsg(key cipher.PubKey) []byte {
	result := make([]byte, MSG_PUBLIC_KEY_END)
	result[MSG_OP_BEGIN] = OP_REG
	copy(result[MSG_PUBLIC_KEY_BEGIN:], key[:])
	return result
}

func GenSendMsg(from, to cipher.PubKey, msg []byte) []byte {
	result := make([]byte, MSG_TO_PUBLIC_KEY_END+len(msg))
	result[MSG_OP_BEGIN] = OP_SEND
	copy(result[MSG_PUBLIC_KEY_BEGIN:], from[:])
	copy(result[MSG_TO_PUBLIC_KEY_BEGIN:], to[:])
	copy(result[MSG_TO_PUBLIC_KEY_END:], msg)
	return result
}
