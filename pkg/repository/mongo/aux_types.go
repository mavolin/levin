package mongo

import (
	"time"
	_ "time/tzdata" // for location loading

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Location time.Location

func (l *Location) UnmarshalBSON(data []byte) error {
	var name string
	if err := bson.Unmarshal(data, &name); err != nil {
		return err
	}

	tl, err := time.LoadLocation(name)
	if err != nil {
		return err
	}

	*l = Location(*tl)
	return nil
}

func (l *Location) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue((*time.Location)(l).String())
}

func (l *Location) Location() *time.Location {
	return (*time.Location)(l)
}
