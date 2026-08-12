package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func main() {
	org := flag.String("org", "", "organization name (required)")
	name := flag.String("name", "", "module name (defaults to the module directory name)")
	provider := flag.String("provider", "", "provider name (required)")
	version := flag.String("version", "", "module version, semver (required)")
	src := flag.String("path", "", "path to the module directory, relative to the repository root (required)")
	flag.Parse()

	var missing []string
	if *org == "" {
		missing = append(missing, "-org")
	}
	if *provider == "" {
		missing = append(missing, "-provider")
	}
	if *version == "" {
		missing = append(missing, "-version")
	}
	if *src == "" {
		missing = append(missing, "-path")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "error: missing required flag(s): %s\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(2)
	}

	if *name == "" {
		*name = filepath.Base(*src)
	}

	modulePath, err := resolvePath(*src)
	if err != nil {
		log.Fatalf("resolve path: %v", err)
	}

	client, err := tfe.NewClient(&tfe.Config{
		Address: "https://tf.pulumi.com",
		Token:   os.Getenv("PULUMI_ACCESS_TOKEN"),
	})
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	ctx := context.Background()

	id := tfe.RegistryModuleID{
		Organization: *org,
		Name:         *name,
		Provider:     *provider,
		Namespace:    *org,
		RegistryName: tfe.PrivateRegistry,
	}

	if _, err := client.RegistryModules.Read(ctx, id); err != nil {
		if !errors.Is(err, tfe.ErrResourceNotFound) {
			log.Fatalf("read module: %v", err)
		}
		log.Printf("creating module %s/%s/%s", *org, *name, *provider)
		if _, err := client.RegistryModules.Create(ctx, *org, tfe.RegistryModuleCreateOptions{
			Name:         tfe.String(*name),
			Provider:     tfe.String(*provider),
			RegistryName: tfe.PrivateRegistry,
		}); err != nil {
			log.Fatalf("create module: %v", err)
		}
	}

	rmv, err := client.RegistryModules.CreateVersion(ctx, id, tfe.RegistryModuleCreateVersionOptions{
		Version: tfe.String(*version),
	})
	if err != nil {
		log.Fatalf("create version: %v", err)
	}

	if err := client.RegistryModules.Upload(ctx, *rmv, modulePath); err != nil {
		log.Fatalf("upload module: %v", err)
	}

	fmt.Printf("uploaded %s/%s/%s@%s\n", *org, *name, *provider, *version)
}

// resolvePath resolves a module path against the repository root so that
// relative paths (e.g. ./modules/s3-bucket) work no matter which directory
// `go -C` runs the program in. Absolute paths are returned unchanged.
func resolvePath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, p), nil
}

// repoRoot walks up from the current working directory looking for the
// repository root, identified by a .git entry (a directory in a normal repo,
// or a file in a worktree).
func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not locate repository root (no .git found)")
		}
		dir = parent
	}
}
