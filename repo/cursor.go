package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Cursor interface {
	Next() bool
	Close() error
	Decode(interface{}) error
}

type cursor struct {
	cursor *mongo.Cursor
}

func NewCursor(c *mongo.Cursor) Cursor {
	return &cursor{cursor: c}
}

func (c *cursor) Next() bool {
	return c.cursor.Next(context.TODO())
}

func (c *cursor) Close() error {
	return c.cursor.Close(context.TODO())
}

func (c *cursor) Decode(val interface{}) error {
	return c.cursor.Decode(val)
}
