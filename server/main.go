package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/shachar1236/Baasa/api"
	"github.com/shachar1236/Baasa/dashboard"
	"github.com/shachar1236/Baasa/database"
)

func run(ctx context.Context, w io.Writer, args []string) error {
    database.Init(ctx)

    err_channel := make(chan error)
    go dashboard.RunDashboard(ctx, err_channel)
    api.RunApi(ctx, err_channel)

    return <-err_channel
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
