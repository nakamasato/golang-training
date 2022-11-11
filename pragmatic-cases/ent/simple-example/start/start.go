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
	// Init Client
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=ent_simple_example password=postgres sslmode=disable") // hardcoding
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Create Item
	ctx := context.Background()
	itemId, _ := UpsertItem(ctx, client, "item_id_1", "Item 1")
	fmt.Printf("id: %s\n", itemId)
	itemId2, _ := UpsertItem(ctx, client, "item_id_2", "Item 2")
	fmt.Printf("id: %s\n", itemId2)

	// Get Item
	_, err = QueryItem(ctx, client, "Item 1")
	if err != nil {
		log.Fatal(err)
	}
}

func UpsertItem(ctx context.Context, client *ent.Client, itemId, itemName string) (string, error) {
	id, err := client.Debug().Item.
		Create().
		SetID(itemId).
		SetName(itemName).
		SetStatus(1).
		OnConflict(
			sql.ConflictColumns(item.FieldID),
		).
		Update(func(u *ent.ItemUpsert) {
			u.SetName(itemName)
			u.SetStatus(1)
		}).
		ID(ctx)
	return id, err
}

func QueryItem(ctx context.Context, client *ent.Client, name string) (*ent.Item, error) {
	i, err := client.Debug().Item.
		Query().
		Where(item.Name(name)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying item: %w", err)
	}
	log.Println("item returned: ", i)
	return i, nil
}
