package customerimporter

import (
	log "log/slog"
	"os"
	"strings"
	"testing"
)

func TestLoadConfigTestEnvFile(t *testing.T) {
	testCases := []struct {
		name    string
		envFile string
		logs    []string
	}{
		{
			name:    "OK",
			envFile: "data/test/.env",
			logs: []string{
				"ERROR: Parsing CONCURRENCY variable failed.",
				"ERROR: CONCURRENCY must be greater than 0. But was 0",
				"ERROR: Parsing READ_BUFFER_SIZE_IN_BYTES failed.",
				"ERROR: READ_BUFFER_SIZE_IN_BYTES must be greater than 0. But was 0",
			},
		},
		{
			name:    "Loading non-existent file",
			envFile: "data/test/.envNonExistentFile",
			logs: []string{
				"ERROR: Loading .env file.",
				"ERROR: open",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockLogger := NewMockLogger()
			os.Clearenv()
			f, err := os.Create("../data/test/.env")
			if err != nil {
				log.Error("Error creating mock .env for tests.")
			}
			f.Close()

			// When
			_, _ = LoadConfigTest(mockLogger, tc.envFile)

			// Check the logs.
			for _, expectedMessage := range tc.logs {
				found := false
				for _, log := range mockLogger.Logs {
					if strings.Contains(log, expectedMessage) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Log message not found: \"%s\"", expectedMessage)
				}
			}

			os.Remove(f.Name())
		})
	}
}

func TestLoadConfigTestEnvVars(t *testing.T) {
	testCases := []struct {
		name    string
		envVars string
		logs    []string
	}{
		{
			name: "CONCURRENCY variable valid",
			envVars: `CONCURRENCY=4
INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv
READ_BUFFER_SIZE_IN_BYTES=4096`,
			logs: []string{},
		},
		{
			name: "CONCURRENCY variable missing",
			envVars: `INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv
READ_BUFFER_SIZE_IN_BYTES=4096`,
			logs: []string{"ERROR: Parsing CONCURRENCY variable failed.", "ERROR: CONCURRENCY must be greater than 0. But was 0"},
		},
		{
			name: "CONCURRENCY variable negative",
			envVars: `CONCURRENCY=-4
INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv
READ_BUFFER_SIZE_IN_BYTES=4096`,
			logs: []string{"ERROR: CONCURRENCY must be greater than 0. But was -4"},
		},
		{
			name: "READ_BUFFER_SIZE_IN_BYTES variable valid",
			envVars: `CONCURRENCY=4
INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv
READ_BUFFER_SIZE_IN_BYTES=4096`,
			logs: []string{},
		},
		{
			name: "READ_BUFFER_SIZE_IN_BYTES variable missing",
			envVars: `CONCURRENCY=4
INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv`,
			logs: []string{"Parsing READ_BUFFER_SIZE_IN_BYTES failed."},
		},
		{
			name: "READ_BUFFER_SIZE_IN_BYTES variable negative",
			envVars: `CONCURRENCY=4
INPUT_CSV_FILE_PATH_DEFAULT=./data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_0_LINES=../data/test/customers_0_lines.csv
INPUT_CSV_FILE_PATH_10_LINES=../data/test/customers_10_lines.csv
INPUT_CSV_FILE_PATH_3K_LINES=../data/test/customers_3k_lines.csv
INPUT_CSV_FILE_PATH_10M_LINES=../data/test/customers_10m_lines.csv
READ_BUFFER_SIZE_IN_BYTES=-4096`,
			logs: []string{"ERROR: READ_BUFFER_SIZE_IN_BYTES must be greater than 0. But was -4096"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockLogger := NewMockLogger()
			os.Clearenv()
			f, err := os.Create("../data/test/.env")
			if err != nil {
				log.Error("Error creating mock .env for tests.")
			}
			if _, err := f.WriteString(tc.envVars); err != nil {
				log.Error(err.Error())
			}
			f.Close()

			// When
			_, _ = LoadConfigTest(mockLogger, "data/test/.env")

			// Check the logs.
			for _, expectedMessage := range tc.logs {
				found := false
				for _, log := range mockLogger.Logs {
					if strings.Contains(log, expectedMessage) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Log message not found: \"%s\"", expectedMessage)
				}
			}

			os.Remove(f.Name())
		})
	}
}
