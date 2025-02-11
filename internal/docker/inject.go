package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

func Inject(ctx context.Context, dockerFilePath string, cmd string) error {
	dockerFile, err := os.Open(dockerFilePath)
	if err != nil {
		return err
	}
	defer dockerFile.Close()

	dockerFileBytes, err := io.ReadAll(dockerFile)
	if err != nil {
		return err
	}

	dockerFileString := string(dockerFileBytes)
	for _, line := range strings.Split(dockerFileString, "\n") {
		if strings.Contains(line, "ENTRYPOINT") {
			dockerFileString = strings.Replace(dockerFileString, line, fmt.Sprintf("ENTRYPOINT [%s]", cmd), 1)
		}
		if strings.Contains(line, "CMD") {
			dockerFileString = strings.Replace(dockerFileString, line, fmt.Sprintf("CMD [%s]", cmd), 1)
		}
	}

	return os.WriteFile(dockerFilePath, []byte(dockerFileString), 0644)
}
