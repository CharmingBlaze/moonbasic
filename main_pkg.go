package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/internal/pkgmgr"
)

func runPackageCLI(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic install|list|publish [args]")
		return 2
	}
	switch strings.ToLower(strings.TrimSpace(args[0])) {
	case "install":
		return cmdInstall(args[1:])
	case "list":
		return cmdList(args[1:])
	case "publish":
		return cmdPublish(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown package command %q\n", args[0])
		return 2
	}
}

func cmdInstall(argv []string) int {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic install <name|url|dir|zip>")
		return 2
	}
	if err := pkgmgr.Install(rest[0], nil); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func cmdList(argv []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	if len(fs.Args()) != 0 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic list")
		return 2
	}
	if err := pkgmgr.ListInstalled(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

func cmdPublish(argv []string) int {
	fs := flag.NewFlagSet("publish", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	upload := fs.Bool("upload", false, "upload zip to GitHub release (needs GITHUB_TOKEN, MOONBASIC_PUBLISH_REPO, MOONBASIC_PUBLISH_TAG)")
	out := fs.String("o", "", "output zip path (default name-version.zip in cwd)")
	if err := fs.Parse(argv); err != nil {
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: moonbasic publish [-o path.zip] [-upload] <package_dir>")
		return 2
	}
	dir := rest[0]
	zipPath := *out
	if zipPath == "" {
		man, err := pkgmgr.ParseManifestFile(filepath.Join(dir, "manifest.json"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
		zipPath = fmt.Sprintf("%s-%s.zip", man.Name, man.Version)
	}
	if err := pkgmgr.PublishPack(dir, zipPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	if *upload {
		if err := pkgmgr.UploadGitHubReleaseAsset(zipPath); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
	}
	return 0
}
