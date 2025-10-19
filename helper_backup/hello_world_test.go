package helper

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestTableHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "Wildan",
			request:  "Wildan",
			expected: "Hello, Wildan!",
		},
		{
			name:     "Akhmad",
			request:  "Akhmad",
			expected: "Hello, Akhmad!",
		},
		{			
			name:     "Budi",
			request:  "Budi",
			expected: "Hello, Budi!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HelloWorld(test.request)
			require.Equal(t, test.expected, result, "Hasil harusnya sama")
		})
	}
}

func TestSubTest(t *testing.T) {
	t.Run("Wildan", func(t *testing.T) {

		result := HelloWorld("Wildan")
		expected := "Hello, Wildan!"

		require.Equal(t, expected, result, "Hasil harusnya sama")
	})

	t.Run("Akhmad", func(t *testing.T) {

		result := HelloWorld("Akhmad")
		expected := "Hello, Akhmad!"

		require.Equal(t, expected, result, "Hasil harusnya sama")
	})
}


func TestMain(m *testing.M) {
	// before
	fmt.Println("BEFORE UNIT TEST")

	m.Run()

	// after
	fmt.Println("AFTER UNIT TEST")
}


func TestSkip(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Unit test tidak bisa dijalankan di Mac")
	}

	result := HelloWorld("Wildan")
	expected := "Hello, Wildan!"

	require.Equal(t, expected, result, "Hasil harusnya sama")
	fmt.Println("Dieksekusi")
}	


func TestHelloWorldWithAssert(t *testing.T) {
	result := HelloWorld("Wildan")
	expected := "Hello, Wildan!"

	assert.Equal(t, expected, result, "Hasil harusnya sama")
	fmt.Println("Dieksekusi")
}

func TestHelloWorldWithRequire(t *testing.T) {
	result := HelloWorld("Wildan")
	expected := "Hello, Wildan!"

	require.Equal(t, expected, result, "Hasil harusnya sama")
	fmt.Println("Dieksekusi")	
}// TestHelloWorldWithRequire FailNow()

func TestHelloWorld(t *testing.T) {
	result := HelloWorld("Wildan")
	expected := "Hello, Wildan!"

	if result != expected {
		t.FailNow()
	}

	fmt.Println("Dieksekusi")
}

func TestHelloWorldAkhmad(t *testing.T) {
	result := HelloWorld("Akhmad")
	expected := "Hello, Akhmad!"

	if result != expected {
		t.Errorf("Harusnya %s tetapi hasilnya %s", expected, result)
	}

	fmt.Println("Dieksekusi")
}