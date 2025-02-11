package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

func Inject(ctx context.Context, dockerFilePath string, cmd string, deps []string) error {
	useEntrypoint := strings.Contains(cmd, "docker") // In case of docker, we need to use the entrypoint to run the command

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
	var lines []string
	lastCmdIndex := -1
	lastEntrypointIndex := -1

	// First pass: find the last CMD and ENTRYPOINT
	for i, line := range strings.Split(dockerFileString, "\n") {
		if strings.Contains(line, "ENTRYPOINT") {
			lastEntrypointIndex = i
		}
		if strings.Contains(line, "CMD") {
			lastCmdIndex = i
		}
		lines = append(lines, line)
	}

	// Replace only the last occurrence
	lastIndex := lastCmdIndex
	if lastEntrypointIndex > lastCmdIndex {
		lastIndex = lastEntrypointIndex
		if useEntrypoint {
			entrypoint := strings.Split(cmd, ",")[:len(strings.Split(cmd, ","))-1]
			entryPointFromLastIndex := strings.Split(lines[lastIndex], "ENTRYPOINT [")[1]
			entryPointFromLastIndex = strings.Split(entryPointFromLastIndex, "]")[0]
			lines[lastIndex] = fmt.Sprintf("ENTRYPOINT [%s,%s]", strings.Join(entrypoint, ","), entryPointFromLastIndex)
		} else {
			lines[lastIndex] = fmt.Sprintf("ENTRYPOINT [\"supergateway\",\"--stdio\",\"%s\"]", cmd)
		}
	} else if lastIndex != -1 {
		if useEntrypoint {
			entrypoint := strings.Split(cmd, ",")[:len(strings.Split(cmd, ","))-1]
			entryPointFromLastIndex := strings.Split(lines[lastIndex], "CMD [")[1]
			entryPointFromLastIndex = strings.Split(entryPointFromLastIndex, "]")[0]
			lines[lastIndex] = fmt.Sprintf("CMD [%s,%s]", strings.Join(entrypoint, ","), entryPointFromLastIndex)
		} else {
			lines[lastIndex] = fmt.Sprintf("CMD [\"supergateway\",\"--stdio\",\"%s\"]", cmd)
		}
	}

	// Add a newline before the last line
	if len(lines) > 0 {
		for _, dep := range deps {
			lines = append(lines[:len(lines)-1], fmt.Sprintf("RUN %s", dep), lines[len(lines)-1])
		}
	}

	return os.WriteFile(dockerFilePath, []byte(strings.Join(lines, "\n")), 0644)
}
