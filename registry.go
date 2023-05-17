package mongonull

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// BuildDefaultRegistry is a bsoncodec.Registry with codec for
// null.String, null.Int, null.Float, null.Bool and null.Time
func BuildDefaultRegistry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

	rb.
		RegisterCodec(nullStringType, nullStringCodec{}).
		RegisterCodec(nullIntType, nullIntCodec{}).
		RegisterCodec(nullFloatType, nullFloatCodec{}).
		RegisterCodec(nullBoolType, nullBoolCodec{}).
		RegisterCodec(nullTimeType, nullTimeCodec{})

	return rb.Build()
}
