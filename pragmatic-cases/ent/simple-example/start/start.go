package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"tmp/pragmatic-cases/ent/simple-example/ent"
	"tmp/pragmatic-cases/ent/simple-example/ent/category"
	"tmp/pragmatic-cases/ent/simple-example/ent/item"
	"tmp/pragmatic-cases/ent/simple-example/ent/predicate"

	"entgo.io/ent/dialect"
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

	// Create Category
	UpsertCategory(ctx, client)

	// Get Category
	category, err := QueryCategory(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	client.Debug().Item.UpdateOneID("item_id_1").AddCategories(category).Save(ctx)
	client.Debug().Item.UpdateOneID("item_id_2").AddCategories(category).Save(ctx)

	QueryCategoryForItem(ctx, category)

	// Practice
	fmt.Println("Original Query:")
	items, _ := client.Debug().Item.Query().Where(item.Or(
		item.CreatedAt(time.Now()),
		item.And(item.CreatedAtEQ(time.Now()),
			item.IDLT("item_id_1")),
	)).All(ctx)
	for _, i := range items {
		fmt.Println(i)
	}

	fmt.Println("New Query:")
	var predicates []predicate.Item // type Item func(*sql.Selector)

	predicates = append(predicates,
		func(s *sql.Selector) {
			s.Where(sql.CompositeGT([]string{s.C(item.FieldCreatedAt), s.C(item.FieldID)}, time.Now(), "item_ids"))
		},
	)
	client.Debug().Item.
		Query().
		Where(predicates...).AllX(ctx)

	t1 := sql.Table("items")
	dialectBuilder := sql.Dialect(dialect.Postgres).
		Select().
		From(t1).
		Where(sql.CompositeGT(t1.Columns(item.FieldCreatedAt, item.FieldID), time.Now(), "item_id_1"))

	query, _ := dialectBuilder.Query()
	fmt.Println(query)
	// client.Debug().Item.Query()
	// items, _ = client.Debug().Item.Query().Where(
	// 	sql.CompositeGT(item.CreatedAt(time.Now()), item.IDLT("item_id_1"))).All(ctx)
	// for _, i := range items {
	// 	fmt.Println(i)
	// }
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
