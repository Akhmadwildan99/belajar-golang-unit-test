package json

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestBasicJSONOperations(t *testing.T) {
	// Test struct dengan JSON tags
	user := User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		Password:  "secret123", // ini tidak akan muncul di JSON
	}

	// Marshal ke JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("User JSON: %s\n", jsonData)

	// Unmarshal kembali
	var newUser User
	err = json.Unmarshal(jsonData, &newUser)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Unmarshaled User: %+v\n", newUser)
}

func TestCustomTimeMarshaling(t *testing.T) {
	event := Event{
		ID:          1,
		Title:       "Meeting",
		Description: "Team meeting",
		EventTime:   CustomTime{time.Now()},
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Event JSON: %s\n", jsonData)

	var newEvent Event
	err = json.Unmarshal(jsonData, &newEvent)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Unmarshaled Event: %+v\n", newEvent)
}

func TestDynamicJSONHandling(t *testing.T) {
	// Test dengan interface
	products := []Product{
		Book{
			Type:   "book",
			Title:  "Go Programming",
			Author: "John Smith",
			ISBN:   "978-1234567890",
		},
		Electronics{
			Type:     "electronics",
			Name:     "Laptop",
			Brand:    "TechBrand",
			Warranty: 24,
		},
	}

	for _, product := range products {
		jsonData, err := json.Marshal(product)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s JSON: %s\n", product.GetType(), jsonData)
	}
}

func TestFlexibleNumberHandling(t *testing.T) {
	// Test parsing berbagai tipe number
	testCases := []string{
		`{"value": 42}`,
		`{"value": 3.14}`,
		`{"value": "123"}`,
		`{"value": "hello"}`,
	}

	for _, testCase := range testCases {
		var data struct {
			Value FlexibleNumber `json:"value"`
		}

		err := json.Unmarshal([]byte(testCase), &data)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Input: %s, Parsed Value: %v (Type: %T)\n", 
			testCase, data.Value.Value, data.Value.Value)
	}
}

func TestRawJSONHandling(t *testing.T) {
	// Test dengan json.RawMessage
	jsonStr := `{
		"type": "user_data",
		"data": {
			"name": "Alice",
			"age": 25,
			"preferences": {
				"theme": "dark",
				"language": "en"
			}
		}
	}`

	var dynamicData DynamicData
	err := json.Unmarshal([]byte(jsonStr), &dynamicData)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Type: %s\n", dynamicData.Type)
	fmt.Printf("Raw Data: %s\n", dynamicData.Data)

	// Parse raw data berdasarkan type
	if dynamicData.Type == "user_data" {
		var userData struct {
			Name        string `json:"name"`
			Age         int    `json:"age"`
			Preferences struct {
				Theme    string `json:"theme"`
				Language string `json:"language"`
			} `json:"preferences"`
		}

		err = json.Unmarshal(dynamicData.Data, &userData)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Parsed User Data: %+v\n", userData)
	}
}

func TestEmbeddedStruct(t *testing.T) {
	article := Article{
		BaseEntity: BaseEntity{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:   "Advanced Go JSON",
		Content: "This is an article about JSON in Go",
		Author:  "Go Developer",
	}

	fmt.Printf("Pretty Print:\n%s\n", PrettyPrint(article))
}

func TestUtilityFunctions(t *testing.T) {
	user := User{
		ID:    1,
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Age:   28,
	}

	// Convert to map
	userMap, err := ToMap(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("User as Map: %+v\n", userMap)

	// Test JSON merging
	json1 := `{"name": "John", "age": 30}`
	json2 := `{"email": "john@example.com", "city": "Jakarta"}`

	merged, err := MergeJSON(json1, json2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Merged JSON: %s\n", merged)

	// Test deep copy
	var userCopy User
	err = DeepCopy(user, &userCopy)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Original: %+v\n", user)
	fmt.Printf("Copy: %+v\n", userCopy)
}

func TestComplexJSONScenario(t *testing.T) {
	// Scenario kompleks: API response dengan nested data
	apiResponse := `{
		"status": "success",
		"data": {
			"users": [
				{
					"id": 1,
					"name": "Alice",
					"profile": {
						"age": 25,
						"location": "Jakarta"
					}
				},
				{
					"id": 2,
					"name": "Bob",
					"profile": {
						"age": 30,
						"location": "Bandung"
					}
				}
			],
			"meta": {
				"total": 2,
				"page": 1,
				"per_page": 10
			}
		}
	}`

	var response map[string]interface{}
	err := json.Unmarshal([]byte(apiResponse), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Extract data menggunakan type assertion
	data := response["data"].(map[string]interface{})
	users := data["users"].([]interface{})
	meta := data["meta"].(map[string]interface{})

	fmt.Printf("Status: %s\n", response["status"])
	fmt.Printf("Total Users: %.0f\n", meta["total"])

	for i, user := range users {
		userMap := user.(map[string]interface{})
		profile := userMap["profile"].(map[string]interface{})
		
		fmt.Printf("User %d: %s (Age: %.0f, Location: %s)\n",
			i+1,
			userMap["name"],
			profile["age"],
			profile["location"])
	}
}