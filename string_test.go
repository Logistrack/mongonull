package mongonull

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsonrw/bsonrwtest"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"gopkg.in/guregu/null.v4"
)

func Test_nullStringCodec_EncodeValue(t *testing.T) {
	tests := []struct {
		name  string
		value any
		b     []byte
		err   error
	}{
		{
			name:  "string",
			value: map[string]null.String{"foo": null.StringFrom("a string")},
			b:     buildDocument(bsoncore.AppendStringElement(nil, "foo", "a string")),
			err:   nil,
		},
		{
			name:  "null",
			value: map[string]null.String{"foo": null.NewString("null", false)},
			b:     buildDocument(bsoncore.AppendNullElement(nil, "foo")),
			err:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := make(bsonrw.SliceWriter, 0, 512)
			vw, err := bsonrw.NewBSONValueWriter(&b)
			noerr(t, err)
			reg := BuildDefaultRegistry()
			enc, err := reg.LookupEncoder(reflect.TypeOf(tt.value))
			noerr(t, err)
			err = enc.EncodeValue(bsoncodec.EncodeContext{Registry: reg}, vw, reflect.ValueOf(tt.value))
			if err != tt.err {
				t.Errorf("EncodeValue() error = %v", err)
			}
			if diff := cmp.Diff([]byte(b), tt.b); diff != "" {
				t.Errorf("Bytes written differ: (-got +want)\n%s", diff)
				t.Errorf("Bytes\ngot: %v\nwant:%v\n", b, tt.b)
				t.Errorf("Readers\ngot: %v\nwant:%v\n", bsoncore.Document(b), bsoncore.Document(tt.b))
			}
		})
	}
}

func Test_nullStringCodec_DecodeValue(t *testing.T) {
	tests := []struct {
		name   string
		reader *bsonrwtest.ValueReaderWriter
		want   null.String
		err    error
	}{
		{
			name:   "string",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.String, Return: "a string"},
			want:   null.StringFrom("a string"),
		},
		{
			name:   "null",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.Null, Return: nil},
			want:   null.NewString("null", false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := nullStringCodec{}
			val := reflect.New(reflect.TypeOf(null.String{})).Elem()
			if err := n.DecodeValue(bsoncodec.DecodeContext{}, tt.reader, val); err != tt.err {
				t.Errorf("DecodeValue() error = %v", err)
			}
			got := val.Interface().(null.String)
			if !got.Equal(tt.want) {
				t.Errorf("got: %+v, want: %+v", got, tt.want)
			}
		})
	}
}

func noerr(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Errorf("Unexpected error: (%T)%v", err, err)
		t.FailNow()
	}
}

func buildDocument(elems []byte) []byte {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = append(doc, elems...)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
}
