package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/beamlit/mcp-store/internal/docker"
	"github.com/beamlit/mcp-store/internal/git"
	"github.com/beamlit/mcp-store/internal/smithery"
	"github.com/beamlit/mcp-store/internal/store"
	"github.com/spf13/cobra"
)

const (
	tmpDir        = "tmp"
	githubPrefix  = "https://github.com/"
	dockerfileDir = "Dockerfile"
)

var (
	configPath string
	push       bool
	registry   string
	mcp        string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import MCPs from a directory",
	Long:  `import is a CLI tool to import MCPs from a directory`,
	Run:   runImport,
}

func init() {
	importCmd.Flags().StringVarP(&configPath, "config", "c", "", "The path to the config file")
	importCmd.Flags().BoolVarP(&push, "push", "p", false, "Push the images to the registry")
	importCmd.Flags().StringVarP(&registry, "registry", "r", "ghcr.io/beamlit/store", "The registry to push the images to")
	importCmd.Flags().StringVarP(&mcp, "mcp", "m", "", "The MCP to import, if not provided, all MCPs will be imported")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	if configPath == "" {
		cmd.Help()
		return
	}

	store := store.Store{}
	handleError("read config file", store.Read(configPath))
	handleError("validate config file", store.ValidateWithDefaultValues())

	setupTempDirectory()
	defer os.RemoveAll(tmpDir)

	for name, repository := range store.Repositories {
		if mcp != "" && mcp != name {
			continue
		}

		if err := processRepository(name, repository); err != nil {
			log.Printf("Failed to process repository %s: %v", name, err)
		}
	}
}

func processRepository(name string, repository *store.Repository) error {
	repoPath := fmt.Sprintf("%s/%s/%s", tmpDir, strings.TrimPrefix(repository.Repository, githubPrefix), repository.Branch)
	defer git.DeleteRepository(repoPath)

	if _, err := git.CloneRepository(repoPath, repository.Branch, repository.Repository); err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}

	cfg, err := smithery.Parse(filepath.Join(repoPath, repository.SmitheryPath), repository.Overrider)
	if err != nil {
		return fmt.Errorf("parse smithery file: %w", err)
	}

	imageName := fmt.Sprintf("%s/%s:latest", strings.ToLower(registry), strings.ToLower(name))
	smitheryDir := strings.TrimSuffix(repository.SmitheryPath, "/smithery.yaml")
	deps := manageDeps(repository)

	if err := buildAndPushImage(&cfg, repoPath, smitheryDir, imageName, deps); err != nil {
		return fmt.Errorf("build and push image: %w", err)
	}

	return nil
}

func buildAndPushImage(cfg *smithery.SmitheryConfig, repoPath, smitheryDir, imageName string, deps []string) error {
	dockerfilePath := filepath.Join(repoPath, smitheryDir, dockerfileDir)
	if err := docker.Inject(context.Background(), dockerfilePath, cfg.ParsedCommand.String(), deps); err != nil {
		return fmt.Errorf("inject command: %w", err)
	}

	buildContext := "."
	if cfg.Build != nil && cfg.Build.DockerBuildPath != nil {
		buildContext = *cfg.Build.DockerBuildPath
	}

	if err := docker.BuildImage(context.Background(), imageName, filepath.Join(smitheryDir, dockerfileDir),
		filepath.Join(repoPath, smitheryDir), buildContext); err != nil {
		return fmt.Errorf("build image: %w", err)
	}

	if push {
		if err := docker.PushImage(context.Background(), imageName); err != nil {
			return fmt.Errorf("push image: %w", err)
		}
	}

	return nil
}

func setupTempDirectory() {
	os.RemoveAll(tmpDir)
	handleError("create temp directory", os.MkdirAll(tmpDir, 0755))
}

func manageDeps(repository *store.Repository) []string {
	switch repository.PackageManager {
	case store.PackageManagerNPM:
		return []string{}
	case store.PackageManagerAPK:
		return []string{
			"apk add --no-cache node npm",
		}
	case store.PackageManagerAPT:
		return []string{
			"apt-get update",
			"apt-get install -y nodejs npm",
		}
	default:
		log.Fatalf("Unsupported package manager: %s", repository.PackageManager)
		return []string{}
	}
}
