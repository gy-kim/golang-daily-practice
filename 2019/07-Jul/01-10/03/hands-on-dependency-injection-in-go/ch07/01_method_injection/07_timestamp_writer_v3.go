package method_injection

import (
	"fmt"
	"io"
	"os"
	"time"
)

func TimeStampWriterV3(writer io.Writer, message string) {
	if writer == nil {
		// define to Standard Out
		writer = os.Stdout
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(writer, "%s -> %s", timestamp, message)
}
