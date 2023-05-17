package mongonull

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsonrw/bsonrwtest"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"gopkg.in/guregu/null.v4"
)

func Test_EncodeValue(t *testing.T) {
	nullDocument := buildDocument(bsoncore.AppendNullElement(nil, "foo"))
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
			name:  "null.String",
			value: map[string]null.String{"foo": {}},
			b:     nullDocument,
			err:   nil,
		},
		{
			name:  "int",
			value: map[string]null.Int{"foo": null.IntFrom(10)},
			b:     buildDocument(bsoncore.AppendInt64Element(nil, "foo", 10)),
			err:   nil,
		},
		{
			name:  "null.Int",
			value: map[string]null.Int{"foo": {}},
			b:     nullDocument,
			err:   nil,
		},
		{
			name:  "float",
			value: map[string]null.Float{"foo": null.FloatFrom(10.5)},
			b:     buildDocument(bsoncore.AppendDoubleElement(nil, "foo", 10.5)),
			err:   nil,
		},
		{
			name:  "null.Float",
			value: map[string]null.Float{"foo": {}},
			b:     nullDocument,
			err:   nil,
		},
		{
			name:  "bool",
			value: map[string]null.Bool{"foo": null.BoolFrom(true)},
			b:     buildDocument(bsoncore.AppendBooleanElement(nil, "foo", true)),
			err:   nil,
		},
		{
			name:  "null.Bool",
			value: map[string]null.Bool{"foo": {}},
			b:     nullDocument,
			err:   nil,
		},
		{
			name:  "time",
			value: map[string]null.Time{"foo": null.TimeFrom(time.UnixMilli(1684328944555))},
			b:     buildDocument(bsoncore.AppendDateTimeElement(nil, "foo", int64(1684328944555))),
			err:   nil,
		},
		{
			name:  "null.Time",
			value: map[string]null.Time{"foo": {}},
			b:     nullDocument,
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
			if err := enc.EncodeValue(bsoncodec.EncodeContext{Registry: reg}, vw, reflect.ValueOf(tt.value)); err != tt.err {
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

func Test_DecodeValue(t *testing.T) {
	nullReader := &bsonrwtest.ValueReaderWriter{BSONType: bsontype.Null, Return: nil}
	tests := []struct {
		name   string
		reader *bsonrwtest.ValueReaderWriter
		want   any
		err    error
	}{
		{
			name:   "string",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.String, Return: "a string"},
			want:   null.StringFrom("a string"),
		},
		{
			name:   "null.String",
			reader: nullReader,
			want:   null.String{},
		},
		{
			name:   "int",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.Int64, Return: int64(10)},
			want:   null.IntFrom(10),
		},
		{
			name:   "null.Int",
			reader: nullReader,
			want:   null.Int{},
		},
		{
			name:   "float",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.Double, Return: 10.5},
			want:   null.FloatFrom(10.5),
		},
		{
			name:   "null.Float",
			reader: nullReader,
			want:   null.Float{},
		},
		{
			name:   "bool",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.Boolean, Return: true},
			want:   null.BoolFrom(true),
		},
		{
			name:   "null.Bool",
			reader: nullReader,
			want:   null.Bool{},
		},
		{
			name:   "time",
			reader: &bsonrwtest.ValueReaderWriter{BSONType: bsontype.DateTime, Return: int64(1684328944555)},
			want:   null.TimeFrom(time.UnixMilli(1684328944555)),
		},
		{
			name:   "null.Time",
			reader: nullReader,
			want:   null.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := BuildDefaultRegistry()
			enc, err := reg.LookupDecoder(reflect.TypeOf(tt.want))
			noerr(t, err)
			val := reflect.New(reflect.TypeOf(tt.want)).Elem()
			if err := enc.DecodeValue(bsoncodec.DecodeContext{}, tt.reader, val); err != tt.err {
				t.Errorf("DecodeValue() error = %v", err)
			}
			if !val.Equal(reflect.ValueOf(tt.want)) {
				t.Errorf("got: %+v, want: %+v", val, tt.want)
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
