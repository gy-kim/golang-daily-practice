package isp

import (
	"context"
	"errors"
)

func Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	// AS this operation make take too long, we need to be able to kill it
	stop := ctx.Done()
	result := make(chan []byte, 1)

	go func() {
		defer close(result)

		// put the encryption key from context
		keyRaw := ctx.Value("encrption-key")
		if keyRaw == nil {
			panic("encryption key not found in context")
		}
		key := keyRaw.([]byte)

		//perform encrypyion
		ciperText := performEncryption(key, data)

		// signal complete by sending the result
		result <- ciperText
	}()

	select {
	case ciperText := <-result:
		// happy path
		return ciperText, nil
	case <-stop:
		// cancelled
		return nil, errors.New("operation cancelled")
	}
}

func performEncryption(key []byte, data []byte) []byte {
	// TODO: implement
	return nil
}
