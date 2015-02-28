package engine

import (
	ciphers "crypto/cipher"
	"errors"
)

var ErrInvalidPadding = errors.New("pkcs7 padding invalid")
var ErrInvalidLength = errors.New("invalid length; it should be multiple of aes block size")

func encryptAesWithPkcs7Padding(cipherPtr *ciphers.Block, origMsgPtr *[]byte) *[]byte {
	cipher := *cipherPtr
	blockSize := cipher.BlockSize()

	msg := *(addPkcs7Padding(cipher.BlockSize(), origMsgPtr))
	encryptedMsg := make([]byte, len(msg))
	buf := make([]byte, blockSize)
	for i := 0; i < len(msg); i += blockSize {
		beg, end := i, i+blockSize
		cipher.Encrypt(buf, msg[beg:end])
		copy(encryptedMsg[beg:end], buf)
	}
	return &encryptedMsg
}

func decryptAesWithPkcs7Padding(cipherPtr *ciphers.Block, encryptedMsgPtr *[]byte) (*[]byte, error) {
	cipher := *cipherPtr
	encryptedMsg := *encryptedMsgPtr
	blockSize := cipher.BlockSize()

	if len(encryptedMsg)%blockSize != 0 {
		return nil, ErrInvalidLength
	}

	msg := make([]byte, len(encryptedMsg))
	buf := make([]byte, blockSize)
	for i := 0; i < len(encryptedMsg); i += blockSize {
		beg, end := i, i+blockSize
		cipher.Decrypt(buf, encryptedMsg[beg:end])
		copy(msg[beg:end], buf)
	}
	return removePkcs7Padding(cipher.BlockSize(), &msg)
}

func addPkcs7Padding(blockSize int, origMsgPtr *[]byte) *[]byte {
	origMsg := *origMsgPtr

	paddingLength := blockSize - (len(origMsg) % blockSize)

	msg := make([]byte, len(origMsg), len(origMsg)+paddingLength)
	copy(msg, origMsg)

	for i := 0; i < paddingLength; i++ {
		msg = append(msg, byte(paddingLength))
	}

	return &msg
}

func removePkcs7Padding(blockSize int, paddedMsgPtr *[]byte) (*[]byte, error) {
	paddedMsg := *paddedMsgPtr
	// validate padding
	paddingLength := int(paddedMsg[len(paddedMsg)-1])
	for _, padding := range paddedMsg[len(paddedMsg)-paddingLength : len(paddedMsg)] {
		if int(padding) != paddingLength {
			return nil, ErrInvalidPadding
		}
	}
	msg := paddedMsg[0 : len(paddedMsg)-paddingLength]

	return &msg, nil
}
