package mongonull

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/guregu/null.v4"
)

var nullFloatType = reflect.TypeOf(null.Float{})

type nullFloatCodec struct{}

func (n nullFloatCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != nullFloatType {
		return bsoncodec.ValueEncoderError{
			Name:     "NullFloatEncodeValue",
			Types:    []reflect.Type{nullFloatType},
			Received: val,
		}
	}

	c := val.Interface().(null.Float)

	if c.Valid {
		return vw.WriteDouble(c.Float64)
	}

	return vw.WriteNull()
}

func (n nullFloatCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{
			Name:     "NullFloatDecodeValue",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var result null.Float
	switch vr.Type() {
	case bsontype.Double:
		s, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		result = null.FloatFrom(s)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}
		result = null.Float{}
	default:
		return fmt.Errorf("received invalid BSON type to decode into null.Float: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(result))
	return nil
}
