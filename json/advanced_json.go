package json

import (
	"encoding/json"
	"fmt"
	// "reflect"
	"strconv"
	"strings"
	"time"
)

// Custom JSON tags dan omitempty
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Age       int       `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"-"` // field ini tidak akan di-serialize
}

// Custom JSON marshaling/unmarshaling
type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format("2006-01-02 15:04:05"))
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// Struct dengan custom time
type Event struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	EventTime   CustomTime `json:"event_time"`
}

// Interface untuk dynamic JSON handling
type Product interface {
	GetType() string
}

type Book struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

func (b Book) GetType() string {
	return "book"
}

type Electronics struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Brand    string `json:"brand"`
	Warranty int    `json:"warranty_months"`
}

func (e Electronics) GetType() string {
	return "electronics"
}

// Custom JSON Number handling
type FlexibleNumber struct {
	Value interface{}
}

func (fn *FlexibleNumber) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	
	// Try parsing as int first
	if intVal, err := strconv.Atoi(str); err == nil {
		fn.Value = intVal
		return nil
	}
	
	// Try parsing as float
	if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
		fn.Value = floatVal
		return nil
	}
	
	// Keep as string if not a number
	fn.Value = str
	return nil
}

func (fn FlexibleNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(fn.Value)
}

// Raw JSON handling
type DynamicData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// JSON dengan embedded struct
type BaseEntity struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Article struct {
	BaseEntity
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Utility functions
func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(b)
}

func ToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func MergeJSON(json1, json2 string) (string, error) {
	var map1, map2 map[string]interface{}
	
	if err := json.Unmarshal([]byte(json1), &map1); err != nil {
		return "", err
	}
	
	if err := json.Unmarshal([]byte(json2), &map2); err != nil {
		return "", err
	}
	
	// Merge map2 into map1
	for k, v := range map2 {
		map1[k] = v
	}
	
	result, err := json.Marshal(map1)
	return string(result), err
}

func DeepCopy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}