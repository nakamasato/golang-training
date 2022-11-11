package main

import (
	"context"
	"fmt"
	"log"

	"tmp/pragmatic-cases/ent/simple-example/ent"
	"tmp/pragmatic-cases/ent/simple-example/ent/category"
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

	// Create Item
	itemId, _ := UpsertItem(ctx, client, "item_id_1", "Item 1")
	fmt.Printf("id: %s\n", itemId)
	itemId2, _ := UpsertItem(ctx, client, "item_id_2", "Item 2")
	fmt.Printf("id: %s\n", itemId2)

	item, err := QueryItem(ctx, client, "Item 1")
	if err != nil {
		log.Fatal(err)
	}

	// Create Category
	if id, err := UpsertCategory(ctx, client); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("id: %s\n", id)
	}
	category, err := QueryCategory(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	// Update with Edge Item-To-Category
	client.Debug().Item.UpdateOneID("item_id_1").AddCategories(category).Save(ctx)
	client.Debug().Item.UpdateOneID("item_id_2").AddCategories(category).Save(ctx)
	client.Debug().Category.UpdateOneID("category_id_1").AddItems(item).Save(ctx)

	QueryCategoryForItem(ctx, category)
}

func QueryCategoryForItem(ctx context.Context, category *ent.Category) error {
	items, err := category.QueryItems().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user categories: %w", err)
	}
	// Query the inverse edge.
	for _, i := range items {
		category, err := i.QueryCategories().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying item %q category: %w", i.Name, err)
		}
		log.Printf("item %q category: %q\n", i.Name, category.Name)
	}
	return nil
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

func UpsertCategory(ctx context.Context, client *ent.Client) (string, error) {
	id, err := client.Debug().Category.
		Create().
		SetID("category_id_1"). // UserIDがUnique指定されている場合
		SetName("Category 1").
		OnConflict(
			sql.ConflictColumns(item.FieldID),
		).
		Update(func(u *ent.CategoryUpsert) {
			u.SetName("Category 1")
		}).
		ID(ctx)
	return id, err
}

func QueryCategory(ctx context.Context, client *ent.Client) (*ent.Category, error) {
	c, err := client.Debug().Category.
		Query().
		Where(category.Name("Category 1")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying item: %w", err)
	}
	log.Println("category returned: ", c)
	return c, nil
}
