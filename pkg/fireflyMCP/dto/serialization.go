package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

// SerializationFormat represents the supported serialization formats
type SerializationFormat string

const (
	// FormatJSON represents JSON serialization format
	FormatJSON SerializationFormat = "json"
	// FormatYAML represents YAML serialization format
	FormatYAML SerializationFormat = "yaml"
	// FormatMessagePack represents MessagePack binary serialization format
	FormatMessagePack SerializationFormat = "msgpack"
)

// Serializable is an interface for DTOs that support serialization
type Serializable interface {
	// Marshal serializes the DTO to the specified format
	Marshal(format SerializationFormat) ([]byte, error)
	// MarshalTo writes the serialized DTO to the writer in the specified format
	MarshalTo(w io.Writer, format SerializationFormat) error
}

// Deserializable is an interface for DTOs that support deserialization
type Deserializable interface {
	// Unmarshal deserializes data in the specified format into the DTO
	Unmarshal(data []byte, format SerializationFormat) error
	// UnmarshalFrom reads and deserializes data from the reader in the specified format
	UnmarshalFrom(r io.Reader, format SerializationFormat) error
}

// SerializableDTO combines both serialization interfaces
type SerializableDTO interface {
	Serializable
	Deserializable
}

// Serializer provides serialization utilities for DTOs
type Serializer struct {
	// Options for JSON encoding
	JSONIndent      string
	JSONPrefix      string
	JSONEscapeHTML  bool
	
	// Options for YAML encoding
	YAMLIndent      int
	
	// Options for MessagePack encoding
	MsgPackUseCompact bool
}

// DefaultSerializer is the default serializer with standard options
var DefaultSerializer = &Serializer{
	JSONIndent:     "  ",
	JSONPrefix:     "",
	JSONEscapeHTML: false,
	YAMLIndent:     2,
	MsgPackUseCompact: true,
}

// Marshal serializes any DTO to the specified format
func (s *Serializer) Marshal(v interface{}, format SerializationFormat) ([]byte, error) {
	switch format {
	case FormatJSON:
		return s.marshalJSON(v)
	case FormatYAML:
		return s.marshalYAML(v)
	case FormatMessagePack:
		return s.marshalMessagePack(v)
	default:
		return nil, fmt.Errorf("unsupported serialization format: %s", format)
	}
}

// Unmarshal deserializes data in the specified format into a DTO
func (s *Serializer) Unmarshal(data []byte, v interface{}, format SerializationFormat) error {
	switch format {
	case FormatJSON:
		return json.Unmarshal(data, v)
	case FormatYAML:
		return yaml.Unmarshal(data, v)
	case FormatMessagePack:
		return msgpack.Unmarshal(data, v)
	default:
		return fmt.Errorf("unsupported serialization format: %s", format)
	}
}

// MarshalTo writes serialized data to a writer
func (s *Serializer) MarshalTo(w io.Writer, v interface{}, format SerializationFormat) error {
	data, err := s.Marshal(v, format)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// UnmarshalFrom reads and deserializes data from a reader
func (s *Serializer) UnmarshalFrom(r io.Reader, v interface{}, format SerializationFormat) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.Unmarshal(data, v, format)
}

func (s *Serializer) marshalJSON(v interface{}) ([]byte, error) {
	if s.JSONIndent == "" && s.JSONPrefix == "" {
		// Use standard marshaling for compact output
		return json.Marshal(v)
	}
	
	// Use indented marshaling
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(s.JSONEscapeHTML)
	encoder.SetIndent(s.JSONPrefix, s.JSONIndent)
	
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	
	// Remove the trailing newline added by Encode
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	
	return data, nil
}

func (s *Serializer) marshalYAML(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(s.YAMLIndent)
	
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

func (s *Serializer) marshalMessagePack(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	
	if s.MsgPackUseCompact {
		encoder.UseCompactInts(true)
		encoder.UseCompactFloats(true)
	}
	
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// Extension methods for DTOs to implement SerializableDTO interface

// Account serialization methods

// Marshal serializes the Account to the specified format
func (a *Account) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(a, format)
}

// MarshalTo writes the serialized Account to the writer
func (a *Account) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, a, format)
}

// Unmarshal deserializes data into the Account
func (a *Account) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, a, format)
}

// UnmarshalFrom reads and deserializes data into the Account
func (a *Account) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, a, format)
}

// Budget serialization methods

// Marshal serializes the Budget to the specified format
func (b *Budget) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(b, format)
}

// MarshalTo writes the serialized Budget to the writer
func (b *Budget) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, b, format)
}

// Unmarshal deserializes data into the Budget
func (b *Budget) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, b, format)
}

// UnmarshalFrom reads and deserializes data into the Budget
func (b *Budget) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, b, format)
}

// Category serialization methods

// Marshal serializes the Category to the specified format
func (c *Category) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(c, format)
}

// MarshalTo writes the serialized Category to the writer
func (c *Category) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, c, format)
}

// Unmarshal deserializes data into the Category
func (c *Category) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, c, format)
}

// UnmarshalFrom reads and deserializes data into the Category
func (c *Category) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, c, format)
}

// Transaction serialization methods

// Marshal serializes the Transaction to the specified format
func (t *Transaction) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(t, format)
}

// MarshalTo writes the serialized Transaction to the writer
func (t *Transaction) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, t, format)
}

// Unmarshal deserializes data into the Transaction
func (t *Transaction) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, t, format)
}

// UnmarshalFrom reads and deserializes data into the Transaction
func (t *Transaction) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, t, format)
}

// TransactionGroup serialization methods

// Marshal serializes the TransactionGroup to the specified format
func (tg *TransactionGroup) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(tg, format)
}

// MarshalTo writes the serialized TransactionGroup to the writer
func (tg *TransactionGroup) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, tg, format)
}

// Unmarshal deserializes data into the TransactionGroup
func (tg *TransactionGroup) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, tg, format)
}

// UnmarshalFrom reads and deserializes data into the TransactionGroup
func (tg *TransactionGroup) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, tg, format)
}

// AccountList serialization methods

// Marshal serializes the AccountList to the specified format
func (al *AccountList) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(al, format)
}

// MarshalTo writes the serialized AccountList to the writer
func (al *AccountList) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, al, format)
}

// Unmarshal deserializes data into the AccountList
func (al *AccountList) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, al, format)
}

// UnmarshalFrom reads and deserializes data into the AccountList
func (al *AccountList) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, al, format)
}

// BudgetList serialization methods

// Marshal serializes the BudgetList to the specified format
func (bl *BudgetList) Marshal(format SerializationFormat) ([]byte, error) {
	return DefaultSerializer.Marshal(bl, format)
}

// MarshalTo writes the serialized BudgetList to the writer
func (bl *BudgetList) MarshalTo(w io.Writer, format SerializationFormat) error {
	return DefaultSerializer.MarshalTo(w, bl, format)
}

// Unmarshal deserializes data into the BudgetList
func (bl *BudgetList) Unmarshal(data []byte, format SerializationFormat) error {
	return DefaultSerializer.Unmarshal(data, bl, format)
}

// UnmarshalFrom reads and deserializes data into the BudgetList
func (bl *BudgetList) UnmarshalFrom(r io.Reader, format SerializationFormat) error {
	return DefaultSerializer.UnmarshalFrom(r, bl, format)
}

// Helper functions for quick serialization

// MarshalJSON is a convenience method for JSON marshaling
func MarshalJSON(v interface{}) ([]byte, error) {
	return DefaultSerializer.Marshal(v, FormatJSON)
}

// UnmarshalJSON is a convenience method for JSON unmarshaling
func UnmarshalJSON(data []byte, v interface{}) error {
	return DefaultSerializer.Unmarshal(data, v, FormatJSON)
}

// MarshalYAML is a convenience method for YAML marshaling
func MarshalYAML(v interface{}) ([]byte, error) {
	return DefaultSerializer.Marshal(v, FormatYAML)
}

// UnmarshalYAML is a convenience method for YAML unmarshaling
func UnmarshalYAML(data []byte, v interface{}) error {
	return DefaultSerializer.Unmarshal(data, v, FormatYAML)
}

// MarshalMessagePack is a convenience method for MessagePack marshaling
func MarshalMessagePack(v interface{}) ([]byte, error) {
	return DefaultSerializer.Marshal(v, FormatMessagePack)
}

// UnmarshalMessagePack is a convenience method for MessagePack unmarshaling
func UnmarshalMessagePack(data []byte, v interface{}) error {
	return DefaultSerializer.Unmarshal(data, v, FormatMessagePack)
}