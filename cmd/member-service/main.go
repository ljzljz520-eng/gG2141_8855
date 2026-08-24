package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"clubmembers/internal/catalog"
	"clubmembers/internal/members"
	"clubmembers/internal/store"
	"clubmembers/internal/workflows"
)

func main() {
	path := "members.db"
	if len(os.Args) > 1 && strings.TrimSpace(os.Args[1]) != "" {
		path = os.Args[1]
	}
	db, err := store.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open member store: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	service := workflows.NewService(db, catalog.DefaultDirectory())
	ctx := context.Background()
	if err := service.SeedDemo(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "seed demo: %v\n", err)
		os.Exit(1)
	}
	rows, err := service.List(ctx, workflows.MemberFilter{Status: members.StatusActive})
	if err != nil {
		fmt.Fprintf(os.Stderr, "list members: %v\n", err)
		os.Exit(1)
	}
	for _, row := range rows {
		fmt.Printf("%s | %s | %s | %s\n", row.ID, row.Name, row.Faculty, row.Phone)
	}
}
