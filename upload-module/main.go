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
	"time"

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
	log.Printf("resolved module path: %s", modulePath)

	token := os.Getenv("PULUMI_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("PULUMI_ACCESS_TOKEN is not set")
	}
	log.Printf("connecting to https://tf.pulumi.com as org %s", *org)

	client, err := tfe.NewClient(&tfe.Config{
		Address: "https://tf.pulumi.com",
		Token:   token,
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
	log.Printf("module ID: org=%s name=%s provider=%s namespace=%s registry=%s",
		id.Organization, id.Name, id.Provider, id.Namespace, id.RegistryName)

	existing, err := client.RegistryModules.Read(ctx, id)
	if err != nil {
		if !errors.Is(err, tfe.ErrResourceNotFound) {
			log.Fatalf("read module: %v", err)
		}
		log.Printf("module not found, creating %s/%s/%s", *org, *name, *provider)
		created, err := client.RegistryModules.Create(ctx, *org, tfe.RegistryModuleCreateOptions{
			Name:         tfe.String(*name),
			Provider:     tfe.String(*provider),
			RegistryName: tfe.PrivateRegistry,
		})
		if err != nil {
			log.Fatalf("create module: %v", err)
		}
		log.Printf("module created: id=%s status=%s", created.ID, created.Status)
	} else {
		log.Printf("module already exists: id=%s status=%s", existing.ID, existing.Status)
	}

	log.Printf("creating version %s", *version)
	rmv, err := client.RegistryModules.CreateVersion(ctx, id, tfe.RegistryModuleCreateVersionOptions{
		Version: tfe.String(*version),
	})
	if err != nil {
		log.Fatalf("create version: %v", err)
	}
	log.Printf("version created: id=%s status=%s upload-url=%s", rmv.ID, rmv.Status, rmv.Links["upload"])

	log.Printf("uploading from %s", modulePath)
	if err := client.RegistryModules.Upload(ctx, *rmv, modulePath); err != nil {
		log.Fatalf("upload: %v", err)
	}
	log.Printf("upload request sent, polling version status...")

	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		mod, err := client.RegistryModules.Read(ctx, id)
		if err != nil {
			log.Printf("poll %d: read error: %v", i+1, err)
			continue
		}
		for _, v := range mod.VersionStatuses {
			if v.Version == *version {
				log.Printf("poll %d: version %s status=%s", i+1, v.Version, v.Status)
			}
		}
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
