package main

import (
	"context"
	"fmt"
	"log"

	flags "github.com/jessevdk/go-flags"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Options struct {
	IdToCopyFrom string `long:"id-to-copy-from" description:"Id to copy from. e.g. " required:"true"`
	NewFileName  string `long:"new-file-name" description:"New filename." required:"true"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// Google Application Default Credentials are used for authentication.
	driveService, err := drive.NewService(ctx, option.WithScopes(drive.DriveScriptsScope))
	if err != nil {
		log.Fatalf("failed to init driveService: %v", err)
	}

	// Copy files
	fileService := drive.NewFilesService(driveService)

	file, err := fileService.Copy(opts.IdToCopyFrom, &drive.File{
		Name: opts.NewFileName,
	}).Do()
	if err != nil {
		log.Fatalf("failed to copy %s: %v", opts.IdToCopyFrom, err)
	}
	fmt.Printf("created %s", file.Id)
}
