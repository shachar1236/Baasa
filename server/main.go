package main

import (
	"context"

	"github.com/shachar1236/Baasa/api"
	"github.com/shachar1236/Baasa/dashboard"
	"github.com/shachar1236/Baasa/database"
)


func main() {
    background_ctx := context.Background()
    ctx, cancel := context.WithCancel(background_ctx)
    defer cancel()

    database.Init(ctx)
    go dashboard.RunDashboard(ctx)
    api.RunApi(ctx)
}
