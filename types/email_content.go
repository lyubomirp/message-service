package types

import (
	"database/sql/driver"
	"encoding/json"
)

type Content struct {
	Recipients []Recipient     `json:"recipients"`
	Subject    json.RawMessage `json:"subject"`
	Content    json.RawMessage `json:"content"`
	Type       string          `json:"type"`
	Format     string          `json:"format"`
}

type Recipient struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

func (sla *Recipient) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &sla)
}

func (sla Recipient) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}
