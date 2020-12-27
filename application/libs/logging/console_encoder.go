package logging

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	_bufferpool = buffer.NewPool()
	bufferGet   = _bufferpool.Get
)

func MilliSecondTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewTimeEncoder(timeFormat string) func(time.Time, zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}
}

const _hex = "0123456789abcdef"

var _jsonPool = sync.Pool{New: func() interface{} {
	return &consoleEncoder{}
}}
var _linePool = sync.Pool{New: func() interface{} {
	return &lineEncoder{}
}}

func getLineEncoder(buf *buffer.Buffer) *lineEncoder {
	enc := _linePool.Get().(*lineEncoder)
	enc.Buffer = buf
	return enc
}
func putLineEncoder(e *lineEncoder) {
	e.Buffer = nil
	_linePool.Put(e)
}

func getConsoleEncoder() *consoleEncoder {
	return _jsonPool.Get().(*consoleEncoder)
}

func putConsoleEncoder(c *consoleEncoder) {
	c.EncoderConfig = nil
	c.buf = nil
	c.spaced = false
	c.openNamespaces = 0
	_jsonPool.Put(c)
}

type consoleEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int
}

func NewConsoleEncoder(cfg *zapcore.EncoderConfig) zapcore.Encoder {
	return newConsoleEncoder(cfg, false)
}

func newConsoleEncoder(cfg *zapcore.EncoderConfig, spaced bool) *consoleEncoder {
	return &consoleEncoder{
		EncoderConfig: cfg,
		buf:           bufferGet(),
		spaced:        spaced,
	}
}

func (c *consoleEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	c.addKey(key)
	return c.AppendArray(arr)
}

func (c *consoleEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	c.addKey(key)
	return c.AppendObject(obj)
}

func (c *consoleEncoder) AddBinary(key string, val []byte) {
	c.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (c *consoleEncoder) AddByteString(key string, val []byte) {
	c.addKey(key)
	c.AppendByteString(val)
}

func (c *consoleEncoder) AddBool(key string, val bool) {
	c.addKey(key)
	c.AppendBool(val)
}

func (c *consoleEncoder) AddComplex128(key string, val complex128) {
	c.addKey(key)
	c.AppendComplex128(val)
}

func (c *consoleEncoder) AddDuration(key string, val time.Duration) {
	c.addKey(key)
	c.AppendDuration(val)
}

func (c *consoleEncoder) AddFloat64(key string, val float64) {
	c.addKey(key)
	c.AppendFloat64(val)
}

func (c *consoleEncoder) AddInt64(key string, val int64) {
	c.addKey(key)
	c.AppendInt64(val)
}

func (c *consoleEncoder) AddReflected(key string, obj interface{}) error {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	c.addKey(key)
	_, err = c.buf.Write(marshaled)
	return err
}

func (c *consoleEncoder) OpenNamespace(key string) {
	c.addKey(key)
	c.buf.AppendByte('{')
	c.openNamespaces++
}

func (c *consoleEncoder) AddString(key, val string) {
	c.addKey(key)
	c.AppendString(val)
}

func (c *consoleEncoder) AddTime(key string, val time.Time) {
	c.addKey(key)
	c.AppendTime(val)
}

func (c *consoleEncoder) AddUint64(key string, val uint64) {
	c.addKey(key)
	c.AppendUint64(val)
}

func (c *consoleEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	c.addElementSeparator()
	c.buf.AppendByte('[')
	err := arr.MarshalLogArray(c)
	c.buf.AppendByte(']')
	return err
}

func (c *consoleEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	c.addElementSeparator()
	c.buf.AppendByte('{')
	err := obj.MarshalLogObject(c)
	c.buf.AppendByte('}')
	return err
}

func (c *consoleEncoder) AppendBool(val bool) {
	c.addElementSeparator()
	c.buf.AppendBool(val)
}

func (c *consoleEncoder) AppendByteString(val []byte) {
	c.addElementSeparator()
	c.buf.AppendByte('"')
	c.safeAddByteString(val)
	c.buf.AppendByte('"')
}

func (c *consoleEncoder) AppendComplex128(val complex128) {
	c.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	c.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	c.buf.AppendFloat(r, 64)
	c.buf.AppendByte('+')
	c.buf.AppendFloat(i, 64)
	c.buf.AppendByte('i')
	c.buf.AppendByte('"')
}

func (c *consoleEncoder) AppendDuration(val time.Duration) {
	cur := c.buf.Len()
	c.EncodeDuration(val, c)
	if cur == c.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		c.AppendInt64(int64(val))
	}
}

func (c *consoleEncoder) AppendInt64(val int64) {
	c.addElementSeparator()
	c.buf.AppendInt(val)
}

func (c *consoleEncoder) AppendReflected(val interface{}) error {
	marshaled, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.addElementSeparator()
	_, err = c.buf.Write(marshaled)
	return err
}

func (c *consoleEncoder) AppendString(val string) {
	c.addElementSeparator()
	c.buf.AppendByte('"')
	c.safeAddString(val)
	c.buf.AppendByte('"')
}

func (c *consoleEncoder) AppendTime(val time.Time) {
	cur := c.buf.Len()
	c.EncodeTime(val, c)
	if cur == c.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		c.AppendInt64(val.UnixNano())
	}
}

func (c *consoleEncoder) AppendUint64(val uint64) {
	c.addElementSeparator()
	c.buf.AppendUint(val)
}

func (c *consoleEncoder) AddComplex64(k string, v complex64) { c.AddComplex128(k, complex128(v)) }
func (c *consoleEncoder) AddFloat32(k string, v float32)     { c.AddFloat64(k, float64(v)) }
func (c *consoleEncoder) AddInt(k string, v int)             { c.AddInt64(k, int64(v)) }
func (c *consoleEncoder) AddInt32(k string, v int32)         { c.AddInt64(k, int64(v)) }
func (c *consoleEncoder) AddInt16(k string, v int16)         { c.AddInt64(k, int64(v)) }
func (c *consoleEncoder) AddInt8(k string, v int8)           { c.AddInt64(k, int64(v)) }
func (c *consoleEncoder) AddUint(k string, v uint)           { c.AddUint64(k, uint64(v)) }
func (c *consoleEncoder) AddUint32(k string, v uint32)       { c.AddUint64(k, uint64(v)) }
func (c *consoleEncoder) AddUint16(k string, v uint16)       { c.AddUint64(k, uint64(v)) }
func (c *consoleEncoder) AddUint8(k string, v uint8)         { c.AddUint64(k, uint64(v)) }
func (c *consoleEncoder) AddUintptr(k string, v uintptr)     { c.AddUint64(k, uint64(v)) }
func (c *consoleEncoder) AppendComplex64(v complex64)        { c.AppendComplex128(complex128(v)) }
func (c *consoleEncoder) AppendFloat64(v float64)            { c.appendFloat(v, 64) }
func (c *consoleEncoder) AppendFloat32(v float32)            { c.appendFloat(float64(v), 32) }
func (c *consoleEncoder) AppendInt(v int)                    { c.AppendInt64(int64(v)) }
func (c *consoleEncoder) AppendInt32(v int32)                { c.AppendInt64(int64(v)) }
func (c *consoleEncoder) AppendInt16(v int16)                { c.AppendInt64(int64(v)) }
func (c *consoleEncoder) AppendInt8(v int8)                  { c.AppendInt64(int64(v)) }
func (c *consoleEncoder) AppendUint(v uint)                  { c.AppendUint64(uint64(v)) }
func (c *consoleEncoder) AppendUint32(v uint32)              { c.AppendUint64(uint64(v)) }
func (c *consoleEncoder) AppendUint16(v uint16)              { c.AppendUint64(uint64(v)) }
func (c *consoleEncoder) AppendUint8(v uint8)                { c.AppendUint64(uint64(v)) }
func (c *consoleEncoder) AppendUintptr(v uintptr)            { c.AppendUint64(uint64(v)) }

func (c *consoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := bufferGet()
	lineEnc := getLineEncoder(line)
	if c.TimeKey != "" && c.EncodeTime != nil {
		c.EncodeTime(ent.Time, lineEnc)
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		nameEncoder := c.EncodeName
		if nameEncoder == nil {
			// Fall back to FullNameEncoder for backward compatibility.
			nameEncoder = zapcore.FullNameEncoder
		}
		c.addTabIfNecessary(line)
		nameEncoder(ent.LoggerName, lineEnc)
	}
	if ent.Caller.Defined && c.CallerKey != "" && c.EncodeCaller != nil {
		c.addTabIfNecessary(line)
		c.EncodeCaller(ent.Caller, lineEnc)
	}
	if c.LevelKey != "" && c.EncodeLevel != nil {
		c.addTabIfNecessary(line)
		c.EncodeLevel(ent.Level, lineEnc)
	}

	putLineEncoder(lineEnc)
	// Add the message itself.
	if c.MessageKey != "" {
		c.addTabIfNecessary(line)
		line.AppendString(ent.Message)
	}
	// Add any structured context.
	c.writeContext(line, fields)
	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	if c.LineEnding != "" {
		line.AppendString(c.LineEnding)
	} else {
		line.AppendString(zapcore.DefaultLineEnding)
	}
	return line, nil

}

func (c *consoleEncoder) Clone() zapcore.Encoder {
	clone := c.clone()
	clone.buf.Write(c.buf.Bytes())
	return clone
}

func (c *consoleEncoder) clone() *consoleEncoder {
	clone := getConsoleEncoder()
	clone.EncoderConfig = c.EncoderConfig
	clone.spaced = c.spaced
	clone.openNamespaces = c.openNamespaces
	clone.buf = bufferGet()
	return clone
}

func (c *consoleEncoder) writeContext(line *buffer.Buffer, extra []zapcore.Field) {
	context := c.Clone().(*consoleEncoder)
	defer context.buf.Free()
	addFields(context, extra)
	context.closeOpenNamespaces()
	if context.buf.Len() == 0 {
		return
	}

	context.addTabIfNecessary(line)
	line.AppendByte('{')
	line.Write(context.buf.Bytes())
	line.AppendByte('}')
	putConsoleEncoder(context)
}

func (c *consoleEncoder) addTabIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendByte(' ')
	}
}

func (c *consoleEncoder) truncate() {
	c.buf.Reset()
}

func (c *consoleEncoder) closeOpenNamespaces() {
	for i := 0; i < c.openNamespaces; i++ {
		c.buf.AppendByte('}')
	}
}

func (c *consoleEncoder) addKey(key string) {
	c.addElementSeparator()
	c.buf.AppendByte('"')
	c.safeAddString(key)
	c.buf.AppendByte('"')
	c.buf.AppendByte(':')
	if c.spaced {
		c.buf.AppendByte(' ')
	}
}

func (c *consoleEncoder) addElementSeparator() {
	last := c.buf.Len() - 1
	if last < 0 {
		return
	}
	switch c.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		c.buf.AppendByte(',')
		if c.spaced {
			c.buf.AppendByte(' ')
		}
	}
}

func (c *consoleEncoder) appendFloat(val float64, bitSize int) {
	c.addElementSeparator()
	switch {
	case math.IsNaN(val):
		c.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		c.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		c.buf.AppendString(`"-Inf"`)
	default:
		c.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (c *consoleEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if c.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if c.tryAddRuneError(r, size) {
			i++
			continue
		}
		c.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (c *consoleEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if c.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if c.tryAddRuneError(r, size) {
			i++
			continue
		}
		c.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (c *consoleEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		c.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		c.buf.AppendByte('\\')
		c.buf.AppendByte(b)
	case '\n':
		c.buf.AppendByte('\\')
		c.buf.AppendByte('n')
	case '\r':
		c.buf.AppendByte('\\')
		c.buf.AppendByte('r')
	case '\t':
		c.buf.AppendByte('\\')
		c.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		c.buf.AppendString(`\u00`)
		c.buf.AppendByte(_hex[b>>4])
		c.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (c *consoleEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		c.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

type lineEncoder struct {
	*buffer.Buffer
}

func (c *lineEncoder) AppendByteString(v []byte) { c.Buffer.AppendString(string(v)) }
func (c *lineEncoder) AppendComplex128(val complex128) {
	r, i := float64(real(val)), float64(imag(val))
	c.Buffer.AppendFloat(r, 64)
	c.Buffer.AppendByte('+')
	c.Buffer.AppendFloat(i, 64)
	c.Buffer.AppendByte('i')
}
func (c *lineEncoder) AppendComplex64(v complex64) { c.AppendComplex128(complex128(v)) }
func (c *lineEncoder) AppendFloat64(v float64)     { c.Buffer.AppendFloat(v, 64) }
func (c *lineEncoder) AppendFloat32(v float32)     { c.Buffer.AppendFloat(float64(v), 32) }
func (c *lineEncoder) AppendInt(v int)             { c.Buffer.AppendInt(int64(v)) }
func (c *lineEncoder) AppendInt64(v int64)         { c.Buffer.AppendInt(v) }
func (c *lineEncoder) AppendInt32(v int32)         { c.Buffer.AppendInt(int64(v)) }
func (c *lineEncoder) AppendInt16(v int16)         { c.Buffer.AppendInt(int64(v)) }
func (c *lineEncoder) AppendInt8(v int8)           { c.Buffer.AppendInt(int64(v)) }
func (c *lineEncoder) AppendString(v string)       { c.Buffer.AppendString(v) }
func (c *lineEncoder) AppendUint(v uint)           { c.Buffer.AppendUint(uint64(v)) }
func (c *lineEncoder) AppendUint64(v uint64)       { c.Buffer.AppendUint(v) }
func (c *lineEncoder) AppendUint32(v uint32)       { c.Buffer.AppendUint(uint64(v)) }
func (c *lineEncoder) AppendUint16(v uint16)       { c.Buffer.AppendUint(uint64(v)) }
func (c *lineEncoder) AppendUint8(v uint8)         { c.Buffer.AppendUint(uint64(v)) }
func (c *lineEncoder) AppendUintptr(v uintptr)     { c.Buffer.AppendUint(uint64(v)) }
