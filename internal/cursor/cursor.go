package cursor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type Direction string

const (
	DirectionAfter  = Direction("after")
	DirectionBefore = Direction("before")
)

type Cursor struct {
	Direction Direction
	Reference string
}

func New(direction Direction, reference string) *Cursor {
	return &Cursor{direction, reference}
}

func (c *Cursor) UnmarshalText(text []byte) error {
	cursor, err := Parse(string(text))
	if err != nil {
		return err
	}

	c.Direction = cursor.Direction
	c.Reference = cursor.Reference
	return nil
}

func (c Cursor) MarshalText() (text []byte, err error) {
	return []byte(c.marshalString()), nil
}

func (c Cursor) String() string {
	return c.string()
}

func (c Cursor) string() string {
	return fmt.Sprintf("%s_%s", c.Direction, c.Reference)
}

func (c Cursor) marshalString() string {
	return base64.StdEncoding.EncodeToString([]byte(c.string()))
}

func Parse(cursor string) (*Cursor, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	split := strings.Split(string(decoded), "_")
	if len(split) < 1 {
		return nil, nil
	}

	if split[0] == string(DirectionBefore) {
		return New(DirectionBefore, split[1]), nil
	}
	if split[0] == string(DirectionAfter) {
		return New(DirectionAfter, split[1]), nil
	}
	return nil, errors.New("invalid cursor direction")
}
