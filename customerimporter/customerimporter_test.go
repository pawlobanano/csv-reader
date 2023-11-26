package customerimporter

import (
	"encoding/csv"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestProcessEmailDomainsConcurrently(t *testing.T) {
	// Given
	log := NewMockLogger()
	config, err := LoadConfig(log, "./.env")
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	file, err := os.Open(config.InputCSVFilePath10Lines)
	if err != nil {
		t.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader, err := createCSVfileReader(log, config, file)
	if err != nil {
		t.Fatalf("Error creating CSV file reader: %v", err)
	}

	// When
	emaiDomains := processEmailDomainsConcurrently(log, config, reader)

	// Then
	expectedEmailDomains := map[string]int{
		"cnet.com":        1,
		"github.com":      2,
		"github.io":       3,
		"hubpages.com":    1,
		"rediff.com":      1,
		"statcounter.com": 1,
	}

	if !reflect.DeepEqual(emaiDomains, expectedEmailDomains) {
		t.Errorf("Unexpected email domains. Expected: %v, Got: %v", expectedEmailDomains, emaiDomains)
	}
}

func TestEmptyInputFile(t *testing.T) {
	// Given
	log := NewMockLogger()
	config, err := LoadConfig(log, "./.env")
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	reader := csv.NewReader(strings.NewReader(""))

	// When
	emaiDomains := processEmailDomainsConcurrently(log, config, reader)

	// Then
	expectedEmailDomains := map[string]int{}

	if !reflect.DeepEqual(emaiDomains, expectedEmailDomains) {
		t.Errorf("Unexpected email domains. Expected: %v, Got: %v", expectedEmailDomains, emaiDomains)
	}
}

func TestSortEmailDomains(t *testing.T) {
	// Given
	emailDomainsWithOccurrences := map[string]int{
		"github.io":       5,
		"github.com":      2,
		"hubpages.com":    1,
		"statcounter.com": 7,
		"rediff.com":      4,
		"cnet.com":        2,
	}

	// When
	sortedDomains := sortEmailDomains(emailDomainsWithOccurrences)

	// Then
	expectedSortedDomains := []string{"cnet.com", "github.com", "github.io", "hubpages.com", "rediff.com", "statcounter.com"}

	if !reflect.DeepEqual(sortedDomains, expectedSortedDomains) {
		t.Errorf("Unexpected sorted domains. Expected: %v, Got: %v", expectedSortedDomains, sortedDomains)
	}
}

func TestExtractDomainInvalidEmail(t *testing.T) {
	// Given
	email := "invalid.email.com"

	// When
	domain := extractDomain(email)

	// Then
	expectedDomain := ""
	if !reflect.DeepEqual(domain, expectedDomain) {
		t.Errorf("Unexpected domain. Expected: %v, Got: %v", expectedDomain, domain)
	}
}

func TestExtractDomainWithTestCases(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		expectedValue string
	}{
		{
			name:          "OK",
			email:         "proper.email@gmail.com",
			expectedValue: "gmail.com",
		},
		{
			name:          "Invalid email",
			email:         "proper.email.gmail.com",
			expectedValue: "",
		},
		{
			name:          "Empty string as email",
			email:         "",
			expectedValue: "",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Given, When
			domain := extractDomain(tc.email)

			// Then
			if !reflect.DeepEqual(domain, tc.expectedValue) {
				t.Errorf("Test %s failed. Expected: %v, Got: %v", tc.name, tc.expectedValue, domain)
			}
		})
	}
}

// NewMockLogger is a helper function to create a mock logger for testing.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Info(msg string, keyVals ...interface{}) {
	m.InfoCalled = true
}

func (m *MockLogger) Warn(msg string, keyVals ...interface{}) {
	m.WarnCalled = true
}

func (m *MockLogger) Error(msg string, keyVals ...interface{}) {
	m.ErrorCalled = true
}
