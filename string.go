package mongonull

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/guregu/null.v4"
)

var nullStringType = reflect.TypeOf(null.String{})

type nullStringCodec struct{}

func (n nullStringCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != nullStringType {
		return bsoncodec.ValueEncoderError{
			Name:     "NullStringEncodeValue",
			Types:    []reflect.Type{nullStringType},
			Received: val,
		}
	}

	c := val.Interface().(null.String)

	if c.Valid {
		return vw.WriteString(c.String)
	}

	return vw.WriteNull()
}

func (n nullStringCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{
			Name:     "NullStringDecodeValue",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var result null.String
	switch vr.Type() {
	case bsontype.String:
		s, err := vr.ReadString()
		if err != nil {
			return err
		}
		result = null.StringFrom(s)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}
		result = null.String{}
	default:
		return fmt.Errorf("received invalid BSON type to decode into null.String: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(result))
	return nil
}
