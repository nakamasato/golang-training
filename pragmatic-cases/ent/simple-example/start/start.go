package main

import (
	"context"
	"fmt"
	"log"

	"tmp/pragmatic-cases/ent/simple-example/ent"
	"tmp/pragmatic-cases/ent/simple-example/ent/item"

	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=ent_simple_example password=postgres sslmode=disable") // hardcoding
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	ctx := context.Background()
	if id, err := UpsertItem(ctx, client); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("id: %s\n", id)
	}
	if _, err = QueryItem(ctx, client); err != nil {
		log.Fatal(err)
	}
}

func UpsertItem(ctx context.Context, client *ent.Client) (string, error) {
	id, err := client.Debug().Item.
		Create().
		SetID("A0001"). // UserIDがUnique指定されている場合
		SetName("a8m").
		SetStatus(1).
		OnConflict(
			sql.ConflictColumns(item.FieldID),
		).
		Update(func(u *ent.ItemUpsert) {
			u.SetName("a8m")
			u.SetStatus(1)
		}).
		ID(ctx)
	return id, err
}

func QueryItem(ctx context.Context, client *ent.Client) (*ent.Item, error) {
	i, err := client.Debug().Item.
		Query().
		Where(item.Name("a8m")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying item: %w", err)
	}
	log.Println("item returned: ", i)
	return i, nil
}
