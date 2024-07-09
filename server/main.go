package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/shachar1236/Baasa/api"
	"github.com/shachar1236/Baasa/dashboard"
)

func run(ctx context.Context, w io.Writer, args []string) error {
    db := GetDatabase(ctx, args)
    access_rules := GetAccessRules(ctx, args, db)

    err_channel := make(chan error)
    go dashboard.RunDashboard(ctx, err_channel, db)
    api.RunApi(ctx, err_channel, db, &access_rules)

    return <-err_channel
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
