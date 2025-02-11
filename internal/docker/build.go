package docker

import (
	"context"
	"os"
	"os/exec"
)

func BuildImage(ctx context.Context, imageName string, dockerfile string, dockerContext string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(dockerContext)
	cmd := exec.Command("docker", "build", "-t", imageName, "-f", dockerfile, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	os.Chdir(currentDir)
	return nil
}
