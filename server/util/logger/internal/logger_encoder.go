package internal

import (
	"encoding/base64"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type SpaceOnlineEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc zapcore.ReflectedEncoder
}

func (enc *SpaceOnlineEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	enc.EncodeTime(ent.Time, final)
	final.buf.AppendString(" ")
	enc.EncodeLevel(ent.Level, final)
	final.buf.AppendString(fmt.Sprintf("\t[GO-%d]", getGOID()))
	enc.EncodeCaller(ent.Caller, final)
	final.buf.AppendString("\t- ")
	final.buf.AppendString(ent.Message)
	final.buf.AppendString(enc.LineEnding)

	ret := final.buf
	putSpaceOnlineEncoder(final)
	return ret, nil
}

func encodeFmtTime(t time.Time, p zapcore.PrimitiveArrayEncoder) {
	encoder := zapcore.TimeEncoderOfLayout("2006/01/02-15:04:05.000")
	encoder(t, p)
}

func encodeFmtLevel(l zapcore.Level, p zapcore.PrimitiveArrayEncoder) {
	text := fmt.Sprintf("[%s]", l.CapitalString())
	p.AppendString(text)
}

func encodeFmtCaller(c zapcore.EntryCaller, p zapcore.PrimitiveArrayEncoder) {
	text := fmt.Sprintf("[%s]", c.TrimmedPath())
	p.AppendString(text)
}

func getGOID() int {
	var id int
	var b [64]byte
	n := runtime.Stack(b[:], false)
	stackStr := string(b[:n])
	stackStr = strings.TrimPrefix(stackStr, "goroutine ")
	fds := strings.Fields(stackStr)
	if len(fds) > 0 {
		nb, _ := strconv.Atoi(fds[0])
		id = nb
	}
	return id
}

// For JSON-escaping; see SpaceOnlineEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var (
	_jsonPool = sync.Pool{
		New: func() interface{} {
			return &SpaceOnlineEncoder{}
		},
	}
	bufferpool = buffer.NewPool()
)

func getSpaceOnlineLogEncoder() *SpaceOnlineEncoder {
	return _jsonPool.Get().(*SpaceOnlineEncoder)
}

func putSpaceOnlineEncoder(enc *SpaceOnlineEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_jsonPool.Put(enc)
}

// NewSpaceOnlineEncoder creates a fast, low-allocation JSON encoder. The encoder
// appropriately escapes all field keys and values.
//
// Note that the encoder doesn't deduplicate keys, so it's possible to produce
// a message like
//
//	{"foo":"bar","foo":"baz"}
//
// This is permitted by the JSON specification, but not encouraged. Many
// libraries will ignore duplicate key-value pairs (typically keeping the last
// pair) when unmarshaling, but users should attempt to avoid adding duplicate
// keys.
func NewSpaceOnlineEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return newSpaceOnlineEncoder(cfg, false)
}

func newSpaceOnlineEncoder(cfg zapcore.EncoderConfig, spaced bool) *SpaceOnlineEncoder {
	if cfg.SkipLineEnding {
		cfg.LineEnding = ""
	} else if cfg.LineEnding == "" {
		cfg.LineEnding = zapcore.DefaultLineEnding
	}

	return &SpaceOnlineEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        spaced,
	}
}

func (enc *SpaceOnlineEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *SpaceOnlineEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *SpaceOnlineEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *SpaceOnlineEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *SpaceOnlineEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *SpaceOnlineEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *SpaceOnlineEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
}

func (enc *SpaceOnlineEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *SpaceOnlineEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *SpaceOnlineEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *SpaceOnlineEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *SpaceOnlineEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = enc.NewReflectedEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *SpaceOnlineEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *SpaceOnlineEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *SpaceOnlineEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *SpaceOnlineEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *SpaceOnlineEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *SpaceOnlineEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *SpaceOnlineEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *SpaceOnlineEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	// Close ONLY new openNamespaces that are created during
	// AppendObject().
	old := enc.openNamespaces
	enc.openNamespaces = 0
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	enc.closeOpenNamespaces()
	enc.openNamespaces = old
	return err
}

func (enc *SpaceOnlineEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *SpaceOnlineEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *SpaceOnlineEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *SpaceOnlineEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *SpaceOnlineEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *SpaceOnlineEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *SpaceOnlineEncoder) AppendString(val string) {
	enc.safeAddString(val)
}

func (enc *SpaceOnlineEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.addElementSeparator()
	enc.buf.AppendTime(time, layout)
}

func (enc *SpaceOnlineEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *SpaceOnlineEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *SpaceOnlineEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *SpaceOnlineEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *SpaceOnlineEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *SpaceOnlineEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *SpaceOnlineEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *SpaceOnlineEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *SpaceOnlineEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *SpaceOnlineEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *SpaceOnlineEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *SpaceOnlineEncoder) AppendComplex64(v complex64)    { enc.appendComplex(complex128(v), 32) }
func (enc *SpaceOnlineEncoder) AppendComplex128(v complex128)  { enc.appendComplex(complex128(v), 64) }
func (enc *SpaceOnlineEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *SpaceOnlineEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *SpaceOnlineEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *SpaceOnlineEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *SpaceOnlineEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *SpaceOnlineEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *SpaceOnlineEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *SpaceOnlineEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *SpaceOnlineEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *SpaceOnlineEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *SpaceOnlineEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }

func (enc *SpaceOnlineEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *SpaceOnlineEncoder) clone() *SpaceOnlineEncoder {
	clone := getSpaceOnlineLogEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()
	return clone
}

func (enc *SpaceOnlineEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *SpaceOnlineEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
	enc.openNamespaces = 0
}

func (enc *SpaceOnlineEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

func (enc *SpaceOnlineEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *SpaceOnlineEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *SpaceOnlineEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s)
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *SpaceOnlineEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
	}
}

func (enc *SpaceOnlineEncoder) tryAddRuneSelf(b byte) bool {
	if b > utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *SpaceOnlineEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
