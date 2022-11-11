// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tmp/pragmatic-cases/ent/simple-example/ent/category"
	"tmp/pragmatic-cases/ent/simple-example/ent/item"
	"tmp/pragmatic-cases/ent/simple-example/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ItemUpdate is the builder for updating Item entities.
type ItemUpdate struct {
	config
	hooks    []Hook
	mutation *ItemMutation
}

// Where appends a list predicates to the ItemUpdate builder.
func (iu *ItemUpdate) Where(ps ...predicate.Item) *ItemUpdate {
	iu.mutation.Where(ps...)
	return iu
}

// SetName sets the "name" field.
func (iu *ItemUpdate) SetName(s string) *ItemUpdate {
	iu.mutation.SetName(s)
	return iu
}

// SetStatus sets the "status" field.
func (iu *ItemUpdate) SetStatus(i int) *ItemUpdate {
	iu.mutation.ResetStatus()
	iu.mutation.SetStatus(i)
	return iu
}

// AddStatus adds i to the "status" field.
func (iu *ItemUpdate) AddStatus(i int) *ItemUpdate {
	iu.mutation.AddStatus(i)
	return iu
}

// SetCreatedAt sets the "created_at" field.
func (iu *ItemUpdate) SetCreatedAt(t time.Time) *ItemUpdate {
	iu.mutation.SetCreatedAt(t)
	return iu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (iu *ItemUpdate) SetNillableCreatedAt(t *time.Time) *ItemUpdate {
	if t != nil {
		iu.SetCreatedAt(*t)
	}
	return iu
}

// AddCategoryIDs adds the "categories" edge to the Category entity by IDs.
func (iu *ItemUpdate) AddCategoryIDs(ids ...string) *ItemUpdate {
	iu.mutation.AddCategoryIDs(ids...)
	return iu
}

// AddCategories adds the "categories" edges to the Category entity.
func (iu *ItemUpdate) AddCategories(c ...*Category) *ItemUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.AddCategoryIDs(ids...)
}

// Mutation returns the ItemMutation object of the builder.
func (iu *ItemUpdate) Mutation() *ItemMutation {
	return iu.mutation
}

// ClearCategories clears all "categories" edges to the Category entity.
func (iu *ItemUpdate) ClearCategories() *ItemUpdate {
	iu.mutation.ClearCategories()
	return iu
}

// RemoveCategoryIDs removes the "categories" edge to Category entities by IDs.
func (iu *ItemUpdate) RemoveCategoryIDs(ids ...string) *ItemUpdate {
	iu.mutation.RemoveCategoryIDs(ids...)
	return iu
}

// RemoveCategories removes "categories" edges to Category entities.
func (iu *ItemUpdate) RemoveCategories(c ...*Category) *ItemUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iu.RemoveCategoryIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (iu *ItemUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(iu.hooks) == 0 {
		affected, err = iu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ItemMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			iu.mutation = mutation
			affected, err = iu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(iu.hooks) - 1; i >= 0; i-- {
			if iu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, iu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (iu *ItemUpdate) SaveX(ctx context.Context) int {
	affected, err := iu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (iu *ItemUpdate) Exec(ctx context.Context) error {
	_, err := iu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iu *ItemUpdate) ExecX(ctx context.Context) {
	if err := iu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iu *ItemUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   item.Table,
			Columns: item.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: item.FieldID,
			},
		},
	}
	if ps := iu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iu.mutation.Name(); ok {
		_spec.SetField(item.FieldName, field.TypeString, value)
	}
	if value, ok := iu.mutation.Status(); ok {
		_spec.SetField(item.FieldStatus, field.TypeInt, value)
	}
	if value, ok := iu.mutation.AddedStatus(); ok {
		_spec.AddField(item.FieldStatus, field.TypeInt, value)
	}
	if value, ok := iu.mutation.CreatedAt(); ok {
		_spec.SetField(item.FieldCreatedAt, field.TypeTime, value)
	}
	if iu.mutation.CategoriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.RemovedCategoriesIDs(); len(nodes) > 0 && !iu.mutation.CategoriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iu.mutation.CategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, iu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{item.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// ItemUpdateOne is the builder for updating a single Item entity.
type ItemUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ItemMutation
}

// SetName sets the "name" field.
func (iuo *ItemUpdateOne) SetName(s string) *ItemUpdateOne {
	iuo.mutation.SetName(s)
	return iuo
}

// SetStatus sets the "status" field.
func (iuo *ItemUpdateOne) SetStatus(i int) *ItemUpdateOne {
	iuo.mutation.ResetStatus()
	iuo.mutation.SetStatus(i)
	return iuo
}

// AddStatus adds i to the "status" field.
func (iuo *ItemUpdateOne) AddStatus(i int) *ItemUpdateOne {
	iuo.mutation.AddStatus(i)
	return iuo
}

// SetCreatedAt sets the "created_at" field.
func (iuo *ItemUpdateOne) SetCreatedAt(t time.Time) *ItemUpdateOne {
	iuo.mutation.SetCreatedAt(t)
	return iuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (iuo *ItemUpdateOne) SetNillableCreatedAt(t *time.Time) *ItemUpdateOne {
	if t != nil {
		iuo.SetCreatedAt(*t)
	}
	return iuo
}

// AddCategoryIDs adds the "categories" edge to the Category entity by IDs.
func (iuo *ItemUpdateOne) AddCategoryIDs(ids ...string) *ItemUpdateOne {
	iuo.mutation.AddCategoryIDs(ids...)
	return iuo
}

// AddCategories adds the "categories" edges to the Category entity.
func (iuo *ItemUpdateOne) AddCategories(c ...*Category) *ItemUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.AddCategoryIDs(ids...)
}

// Mutation returns the ItemMutation object of the builder.
func (iuo *ItemUpdateOne) Mutation() *ItemMutation {
	return iuo.mutation
}

// ClearCategories clears all "categories" edges to the Category entity.
func (iuo *ItemUpdateOne) ClearCategories() *ItemUpdateOne {
	iuo.mutation.ClearCategories()
	return iuo
}

// RemoveCategoryIDs removes the "categories" edge to Category entities by IDs.
func (iuo *ItemUpdateOne) RemoveCategoryIDs(ids ...string) *ItemUpdateOne {
	iuo.mutation.RemoveCategoryIDs(ids...)
	return iuo
}

// RemoveCategories removes "categories" edges to Category entities.
func (iuo *ItemUpdateOne) RemoveCategories(c ...*Category) *ItemUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return iuo.RemoveCategoryIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (iuo *ItemUpdateOne) Select(field string, fields ...string) *ItemUpdateOne {
	iuo.fields = append([]string{field}, fields...)
	return iuo
}

// Save executes the query and returns the updated Item entity.
func (iuo *ItemUpdateOne) Save(ctx context.Context) (*Item, error) {
	var (
		err  error
		node *Item
	)
	if len(iuo.hooks) == 0 {
		node, err = iuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ItemMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			iuo.mutation = mutation
			node, err = iuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(iuo.hooks) - 1; i >= 0; i-- {
			if iuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = iuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, iuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Item)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from ItemMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (iuo *ItemUpdateOne) SaveX(ctx context.Context) *Item {
	node, err := iuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (iuo *ItemUpdateOne) Exec(ctx context.Context) error {
	_, err := iuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (iuo *ItemUpdateOne) ExecX(ctx context.Context) {
	if err := iuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (iuo *ItemUpdateOne) sqlSave(ctx context.Context) (_node *Item, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   item.Table,
			Columns: item.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: item.FieldID,
			},
		},
	}
	id, ok := iuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Item.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := iuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, item.FieldID)
		for _, f := range fields {
			if !item.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != item.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := iuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := iuo.mutation.Name(); ok {
		_spec.SetField(item.FieldName, field.TypeString, value)
	}
	if value, ok := iuo.mutation.Status(); ok {
		_spec.SetField(item.FieldStatus, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.AddedStatus(); ok {
		_spec.AddField(item.FieldStatus, field.TypeInt, value)
	}
	if value, ok := iuo.mutation.CreatedAt(); ok {
		_spec.SetField(item.FieldCreatedAt, field.TypeTime, value)
	}
	if iuo.mutation.CategoriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.RemovedCategoriesIDs(); len(nodes) > 0 && !iuo.mutation.CategoriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := iuo.mutation.CategoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   item.CategoriesTable,
			Columns: item.CategoriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: category.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Item{config: iuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, iuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{item.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
