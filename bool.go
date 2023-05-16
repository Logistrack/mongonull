package mongonull

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/guregu/null.v4"
)

var nullBoolType = reflect.TypeOf(null.Bool{})

type nullBoolCodec struct{}

func (n nullBoolCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != nullBoolType {
		return bsoncodec.ValueEncoderError{
			Name:     "NullBoolEncodeValue",
			Types:    []reflect.Type{nullBoolType},
			Received: val,
		}
	}

	c := val.Interface().(null.Bool)

	if c.Valid {
		return vw.WriteBoolean(c.Bool)
	}

	return vw.WriteNull()
}

func (n nullBoolCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{
			Name:     "NullBoolDecodeValue",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var result null.Bool
	switch vr.Type() {
	case bsontype.Boolean:
		s, err := vr.ReadBoolean()
		if err != nil {
			return err
		}
		result = null.BoolFrom(s)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}
		result = null.Bool{}
	default:
		return fmt.Errorf("received invalid BSON type to decode into null.Bool: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(result))
	return nil
}
