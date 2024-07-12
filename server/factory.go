package main

import (
	"context"

	"github.com/shachar1236/Baasa/access_rules"
	"github.com/shachar1236/Baasa/database"
	"github.com/shachar1236/Baasa/database/sqlite"
	querylang "github.com/shachar1236/Baasa/query_lang"
)

func GetDatabase(ctx context.Context, args []string) database.Database {
    db, err := sqlite.New(ctx)
    if err != nil {
        panic(err)
    }

    return &db
}

func GetAccessRules(ctx context.Context, args []string, db database.Database) access_rules.AccessRules {
    return access_rules.New(db)
}

func GetQueryLangAnalyzer(ctx context.Context, args []string, db database.Database) querylang.Analyzer {
    analyzer := querylang.New(db)
    return analyzer
}
