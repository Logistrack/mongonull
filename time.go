package mongonull

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/guregu/null.v4"
)

var nullTimeType = reflect.TypeOf(null.Time{})

type nullTimeCodec struct{}

func (n nullTimeCodec) EncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != nullTimeType {
		return bsoncodec.ValueEncoderError{
			Name:     "NullTimeEncodeValue",
			Types:    []reflect.Type{nullTimeType},
			Received: val,
		}
	}

	c := val.Interface().(null.Time)

	if c.Valid {
		return vw.WriteDateTime(c.Time.UnixMilli())
	}

	return vw.WriteNull()
}

func (n nullTimeCodec) DecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{
			Name:     "NullTimeDecodeValue",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var result null.Time
	switch vr.Type() {
	case bsontype.DateTime:
		s, err := vr.ReadDateTime()
		if err != nil {
			return err
		}
		result = null.TimeFrom(time.UnixMilli(s))
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}
		result = null.Time{}
	default:
		return fmt.Errorf("received invalid BSON type to decode into null.Time: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(result))

	return nil
}
