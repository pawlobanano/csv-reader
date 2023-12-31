package customerimporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "log/slog"

	"github.com/joho/godotenv"
)

// LoadConfig loads the configuration from the .env file.
func LoadConfig(log Logger, envFilePath string) (*Config, error) {
	err := Load(envFilePath)
	if err != nil {
		log.Error("Loading .env file.", err)
		return nil, err
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		log.Error("Parsing CONCURRENCY variable failed.")
	}

	if concurrency <= 0 {
		log.Error(fmt.Sprintf("%s %d", "CONCURRENCY must be greater than 0. But was", concurrency))
	}

	readBufferSizeInBytes, err := strconv.Atoi(os.Getenv("READ_BUFFER_SIZE_IN_BYTES"))
	if err != nil {
		log.Error("Parsing READ_BUFFER_SIZE_IN_BYTES failed.")
	}

	if readBufferSizeInBytes <= 0 {
		log.Error(fmt.Sprintf("%s %d", "READ_BUFFER_SIZE_IN_BYTES must be greater than 0. But was", readBufferSizeInBytes))
	}

	config := &Config{
		Concurrency:              concurrency,
		InputCSVFilePathDefault:  os.Getenv("INPUT_CSV_FILE_PATH_DEFAULT"),
		InputCSVFilePath0Lines:   os.Getenv("INPUT_CSV_FILE_PATH_0_LINES"),
		InputCSVFilePath10Lines:  os.Getenv("INPUT_CSV_FILE_PATH_10_LINES"),
		InputCSVFilePath3kLines:  os.Getenv("INPUT_CSV_FILE_PATH_3K_LINES"),
		InputCSVFilePath10mLines: os.Getenv("INPUT_CSV_FILE_PATH_10M_LINES"),
		ReadBufferSizeInBytes:    readBufferSizeInBytes,
	}

	return config, nil
}

// LoadConfigTest loads the configuration from the .env file for tests
func LoadConfigTest(log Logger, envFilePath string) (*Config, error) {
	config, err := LoadConfig(log, envFilePath)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	config = &Config{
		Concurrency:              config.Concurrency,
		InputCSVFilePathDefault:  config.InputCSVFilePath10Lines,
		InputCSVFilePath0Lines:   config.InputCSVFilePath0Lines,
		InputCSVFilePath10Lines:  config.InputCSVFilePath10Lines,
		InputCSVFilePath3kLines:  config.InputCSVFilePath3kLines,
		InputCSVFilePath10mLines: config.InputCSVFilePath10mLines,
		ReadBufferSizeInBytes:    config.ReadBufferSizeInBytes,
	}

	return config, nil
}

// Load loads the environment variables from the .env file.
func Load(envFile string) error { // Solution to differentiate .env file path for unit or benchmark tests; source: https://github.com/joho/godotenv/issues/126#issuecomment-1474645022
	err := godotenv.Load(dir(envFile))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
