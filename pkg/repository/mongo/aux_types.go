package mongo

import (
	"time"
	_ "time/tzdata" // for location loading

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"golang.org/x/text/language"
)

// =============================================================================
// location
// =====================================================================================

type location time.Location

func (l *location) UnmarshalBSON(data []byte) error {
	var name string
	if err := bson.Unmarshal(data, &name); err != nil {
		return err
	}

	tl, err := time.LoadLocation(name)
	if err != nil {
		return err
	}

	*l = location(*tl)
	return nil
}

func (l *location) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue((*time.Location)(l).String())
}

func (l *location) baseType() *time.Location {
	return (*time.Location)(l)
}

// =============================================================================
// languageTag
// =====================================================================================

type languageTag language.Tag

func (t *languageTag) UnmarshalBSON(data []byte) error {
	var tagstr string
	if err := bson.Unmarshal(data, &tagstr); err != nil {
		return err
	}

	lt, err := language.Parse(tagstr)
	if err != nil {
		return err
	}

	*t = languageTag(lt)
	return nil
}

func (t languageTag) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(language.Tag(t).String())
}

func (t languageTag) baseType() language.Tag {
	return language.Tag(t)
}
