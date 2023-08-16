/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
