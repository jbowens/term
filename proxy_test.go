package term

import (
	"bytes"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestEscapeProxyRead(t *testing.T) {
	t.Run("no escape keys, keys [a]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("")
		keys, _ := ToBytes("a")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, len(keys))
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("no escape keys, keys [a,b,c]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("")
		keys, _ := ToBytes("a,b,c")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, len(keys))
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("no escape keys, no keys", func(t *testing.T) {
		escapeKeys, _ := ToBytes("")
		keys, _ := ToBytes("")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.Assert(t, is.ErrorContains(err, ""), "Should throw error when no keys are to read")
		assert.Equal(t, nr, 0)
		assert.Check(t, is.Len(keys, 0))
		assert.Check(t, is.Len(buf, 0))
	})

	t.Run("DEL escape key, keys [a,b,c,+]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("DEL")
		keys, _ := ToBytes("a,b,c,+")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, len(keys))
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("DEL escape key, no keys", func(t *testing.T) {
		escapeKeys, _ := ToBytes("DEL")
		keys, _ := ToBytes("")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.Assert(t, is.ErrorContains(err, ""), "Should throw error when no keys are to read")
		assert.Equal(t, nr, 0)
		assert.Check(t, is.Len(keys, 0))
		assert.Check(t, is.Len(buf, 0))
	})

	t.Run("ctrl-x,ctrl-@ escape key, keys [DEL]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-x,ctrl-@")
		keys, _ := ToBytes("DEL")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, 1)
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("ctrl-c escape key, keys [ctrl-c]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-c")
		keys, _ := ToBytes("ctrl-c")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, len(keys))
		nr, err := reader.Read(buf)
		assert.Error(t, err, "read escape sequence")
		assert.Equal(t, nr, 0)
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[ctrl-z]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-c,ctrl-z")
		keys, _ := ToBytes("ctrl-c,ctrl-z")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, 0)
		assert.DeepEqual(t, keys[0:1], buf)

		nr, err = reader.Read(buf)
		assert.Error(t, err, "read escape sequence")
		assert.Equal(t, nr, 0)
		assert.DeepEqual(t, keys[1:], buf)
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c,ctrl-z]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-c,ctrl-z")
		keys, _ := ToBytes("ctrl-c,ctrl-z")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 2)
		nr, err := reader.Read(buf)
		assert.Error(t, err, "read escape sequence")
		assert.Equal(t, nr, 0, "nr should be equal to 0")
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[DEL,+]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-c,ctrl-z")
		keys, _ := ToBytes("ctrl-c,DEL,+")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, 0)
		assert.DeepEqual(t, keys[0:1], buf)

		buf = make([]byte, len(keys))
		nr, err = reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, len(keys))
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("ctrl-c,ctrl-z escape key, keys [ctrl-c],[DEL]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-c,ctrl-z")
		keys, _ := ToBytes("ctrl-c,DEL")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, 0)
		assert.DeepEqual(t, keys[0:1], buf)

		buf = make([]byte, len(keys))
		nr, err = reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, nr, len(keys))
		assert.DeepEqual(t, keys, buf)
	})

	t.Run("a,b,c,d escape key, keys [a,b],[c,d]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("a,b,c,d")
		keys, _ := ToBytes("a,b,c,d")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 2)
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, 0, nr)
		assert.DeepEqual(t, keys[0:2], buf)

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		assert.Error(t, err, "read escape sequence")
		assert.Equal(t, 0, nr)
		assert.DeepEqual(t, keys[2:4], buf)
	})

	t.Run("ctrl-p,ctrl-q escape key, keys [ctrl-p],[a],[ctrl-p,ctrl-q]", func(t *testing.T) {
		escapeKeys, _ := ToBytes("ctrl-p,ctrl-q")
		keys, _ := ToBytes("ctrl-p,a,ctrl-p,ctrl-q")
		reader := NewEscapeProxy(bytes.NewReader(keys), escapeKeys)

		buf := make([]byte, 1)
		nr, err := reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, 0, nr)

		buf = make([]byte, 1)
		nr, err = reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, 1, nr)
		assert.DeepEqual(t, keys[:1], buf)

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		assert.NilError(t, err)
		assert.Equal(t, 1, nr)
		assert.DeepEqual(t, keys[1:3], buf)

		buf = make([]byte, 2)
		nr, err = reader.Read(buf)
		assert.Error(t, err, "read escape sequence")
		assert.Equal(t, 0, nr)
	})
}
