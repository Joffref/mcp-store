package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/beamlit/mcp-store/internal/docker"
	"github.com/beamlit/mcp-store/internal/git"
	"github.com/beamlit/mcp-store/internal/smithery"
	"github.com/beamlit/mcp-store/internal/store"
	"github.com/spf13/cobra"
)

var (
	configPath string
	push       bool
	registry   string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import MCPs from a directory",
	Long:  `import is a CLI tool to import MCPs from a directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if configPath == "" {
			cmd.Help()
			return
		}

		store := store.Store{}
		err := store.Read(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}

		err = store.ValidateWithDefaultValues()
		if err != nil {
			log.Fatalf("Failed to validate config file: %v", err)
		}

		os.MkdirAll("tmp", 0755)
		defer os.RemoveAll("tmp")

		for name, repository := range store.Repositories {
			repoPath := fmt.Sprintf("tmp/%s/%s", strings.Replace(repository.Repository, "https://github.com/", "", 1), repository.Branch)
			_, err := git.CloneRepository(repoPath, repository.Branch, repository.Repository)
			if err != nil {
				log.Printf("Failed to clone repository: %v", err)
				continue
			}
			cfg, err := smithery.Parse(fmt.Sprintf("%s/%s", repoPath, repository.SmitheryPath))
			if err != nil {
				log.Printf("Failed to parse smithery file: %v", err)
				continue
			}

			imageName := fmt.Sprintf("%s/%s:latest", registry, name)
			smitheryDir := strings.TrimSuffix(repository.SmitheryPath, "smithery.yaml")

			err = docker.Inject(context.Background(), fmt.Sprintf("%s/%s/Dockerfile", repoPath, smitheryDir), fmt.Sprintf("\"npx\",\"-y\",\"supergateway\",\"--stdio\",\"%s\"", cfg.ParsedCommand.String()))
			if err != nil {
				log.Printf("Failed to inject command: %v", err)
				continue
			}

			err = docker.BuildImage(context.Background(), imageName, fmt.Sprintf("%s/Dockerfile", smitheryDir), repoPath)
			if err != nil {
				log.Printf("Failed to build image: %v", err)
				continue
			}
			if push {
				err = docker.PushImage(context.Background(), imageName)
				if err != nil {
					log.Printf("Failed to push image: %v", err)
				}
			}
			git.DeleteRepository(repoPath)
		}
	},
}

func init() {
	importCmd.Flags().StringVarP(&configPath, "config", "c", "", "The path to the config file")
	importCmd.Flags().BoolVarP(&push, "push", "p", false, "Push the images to the registry")
	importCmd.Flags().StringVarP(&registry, "registry", "r", "ghcr.io/beamlit/store", "The registry to push the images to")
	rootCmd.AddCommand(importCmd)
}
