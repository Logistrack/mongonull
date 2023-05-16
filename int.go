package mongonull

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/guregu/null.v4"
)

var nullIntType = reflect.TypeOf(null.Int{})

type nullIntCodec struct{}

func (n nullIntCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != nullIntType {
		return bsoncodec.ValueEncoderError{
			Name:     "NullIntEncodeValue",
			Types:    []reflect.Type{nullIntType},
			Received: val,
		}
	}

	c := val.Interface().(null.Int)

	if c.Valid {
		return vw.WriteInt64(c.Int64)
	}

	return vw.WriteNull()
}

func (n nullIntCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{
			Name:     "NullIntDecodeValue",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var result null.Int
	switch vr.Type() {
	case bsontype.Int64:
		s, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		result = null.IntFrom(s)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}
		result = null.Int{}
	default:
		return fmt.Errorf("received invalid BSON type to decode into null.Int: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(result))
	return nil
}
