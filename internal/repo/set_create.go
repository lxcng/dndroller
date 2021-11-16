// Code generated by entc, DO NOT EDIT.

package repo

import (
	"context"
	"dndroller/internal/model"
	"dndroller/internal/repo/set"
	"dndroller/internal/repo/user"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SetCreate is the builder for creating a Set entity.
type SetCreate struct {
	config
	mutation *SetMutation
	hooks    []Hook
}

// SetData sets the "data" field.
func (sc *SetCreate) SetData(ms *model.DiceSet) *SetCreate {
	sc.mutation.SetData(ms)
	return sc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (sc *SetCreate) SetOwnerID(id int) *SetCreate {
	sc.mutation.SetOwnerID(id)
	return sc
}

// SetOwner sets the "owner" edge to the User entity.
func (sc *SetCreate) SetOwner(u *User) *SetCreate {
	return sc.SetOwnerID(u.ID)
}

// Mutation returns the SetMutation object of the builder.
func (sc *SetCreate) Mutation() *SetMutation {
	return sc.mutation
}

// Save creates the Set in the database.
func (sc *SetCreate) Save(ctx context.Context) (*Set, error) {
	var (
		err  error
		node *Set
	)
	if len(sc.hooks) == 0 {
		if err = sc.check(); err != nil {
			return nil, err
		}
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SetMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sc.check(); err != nil {
				return nil, err
			}
			sc.mutation = mutation
			if node, err = sc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			if sc.hooks[i] == nil {
				return nil, fmt.Errorf("repo: uninitialized hook (forgotten import repo/runtime?)")
			}
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SetCreate) SaveX(ctx context.Context) *Set {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SetCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SetCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SetCreate) check() error {
	if _, ok := sc.mutation.Data(); !ok {
		return &ValidationError{Name: "data", err: errors.New(`repo: missing required field "data"`)}
	}
	if _, ok := sc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New("repo: missing required edge \"owner\"")}
	}
	return nil
}

func (sc *SetCreate) sqlSave(ctx context.Context) (*Set, error) {
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (sc *SetCreate) createSpec() (*Set, *sqlgraph.CreateSpec) {
	var (
		_node = &Set{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: set.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: set.FieldID,
			},
		}
	)
	if value, ok := sc.mutation.Data(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: set.FieldData,
		})
		_node.Data = value
	}
	if nodes := sc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   set.OwnerTable,
			Columns: []string{set.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_set = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SetCreateBulk is the builder for creating many Set entities in bulk.
type SetCreateBulk struct {
	config
	builders []*SetCreate
}

// Save creates the Set entities in the database.
func (scb *SetCreateBulk) Save(ctx context.Context) ([]*Set, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Set, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SetMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SetCreateBulk) SaveX(ctx context.Context) []*Set {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SetCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SetCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
