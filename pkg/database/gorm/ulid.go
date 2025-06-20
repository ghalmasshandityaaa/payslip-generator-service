package gorm

import (
	"database/sql/driver"
	"fmt"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ULID is a custom GORM data type for handling ULID
type ULID ulid.ULID

// Value implements the driver.Valuer interface
func (u ULID) Value() (driver.Value, error) {
	if u == (ULID{}) {
		return nil, nil
	}
	return ulid.ULID(u).String(), nil
}

// Scan implements the sql.Scanner interface
func (u *ULID) Scan(value interface{}) error {
	if value == nil {
		*u = ULID{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		parsed, err := ulid.Parse(string(v))
		if err != nil {
			return fmt.Errorf("failed to parse ULID from bytes: %v", err)
		}
		*u = ULID(parsed)
	case string:
		parsed, err := ulid.Parse(v)
		if err != nil {
			return fmt.Errorf("failed to parse ULID from string: %v", err)
		}
		*u = ULID(parsed)
	default:
		return fmt.Errorf("cannot scan %T into ULID", value)
	}
	return nil
}

// String returns the string representation of the ULID
func (u ULID) String() string {
	return ulid.ULID(u).String()
}

// MarshalJSON implements json.Marshaler
func (u ULID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, u.String())), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ULID) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed, err := ulid.Parse(str)
	if err != nil {
		return fmt.Errorf("failed to parse ULID from JSON: %v", err)
	}
	*u = ULID(parsed)
	return nil
}

// GormDataType returns the GORM data type
func (ULID) GormDataType() string {
	return "ulid"
}

// GormDBDataType returns the database data type
func (ULID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "ulid"
}
