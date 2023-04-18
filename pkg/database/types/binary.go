package types

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type BinaryData uuid.UUID

func StringToBinaryData(s string) (BinaryData, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return BinaryData{}, err
	}
	return BinaryData(u), nil
}

func (b BinaryData) String() string {
	return uuid.UUID(b).String()
}

func (b BinaryData) GormDataType() string {
	return "binary(16)"
}

func (b BinaryData) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uuid.UUID(b).String() + `"`), nil
}

func (b *BinaryData) UnmarshalJSON(data []byte) error {
	s, err := uuid.ParseBytes(data)
	*b = BinaryData(s)
	return err
}

func (b BinaryData) Value() (driver.Value, error) {
	return uuid.UUID(b).String(), nil
}

func (b *BinaryData) Scan(src interface{}) error {
	bytes, _ := src.([]byte)
	parseByte, err := uuid.FromBytes(bytes)
	*b = BinaryData(parseByte)
	return err
}
