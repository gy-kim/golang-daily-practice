package method_injection

import (
	"fmt"
	"io"
	"time"
)

func TimeStampWriterV1(writer io.Writer, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(writer, "%s -> %s", timestamp, message)
}
