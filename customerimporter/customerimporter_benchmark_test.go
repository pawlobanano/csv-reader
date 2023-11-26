package customerimporter

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkProcessEmailDomainsConcurrently(b *testing.B) {
	log := NewMockLogger()
	config, err := LoadConfig(log, "./.env")
	if err != nil {
		b.Fatalf("Error loading config: %v", err)
	}

	filePaths := []string{
		config.InputCSVFilePath10Lines,
		config.InputCSVFilePath3kLines,
		config.InputCSVFilePath10mLines,
	}

	for _, filePath := range filePaths {
		b.Run(fmt.Sprintf("File: %s", filePath), func(b *testing.B) {
			for _, parallelism := range []int{1, 2, 4, 8, 12} {
				b.SetParallelism(parallelism)

				b.Run(fmt.Sprintf("CPU core parallelism: %d", parallelism), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						file, err := os.Open(filePath)
						if err != nil {
							b.Fatal(err)
						}
						defer file.Close()

						reader, err := createCSVfileReader(log, config, file)
						if err != nil {
							b.Fatal(err)
						}

						_ = processEmailDomainsConcurrently(log, config, reader)
					}
				})
			}
		})
	}
}
