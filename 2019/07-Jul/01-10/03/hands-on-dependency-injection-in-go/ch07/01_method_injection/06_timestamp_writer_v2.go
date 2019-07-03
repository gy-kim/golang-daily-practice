package method_injection

import (
	"errors"
	"fmt"
	"io"
	"time"
)

func TimeStampWriterV2(writer io.Writer, message string) error {
	if writer == nil {
		return errors.New("Writer cannot be nil")
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(writer, "%s -> %s", timestamp, message)

	return nil
}
