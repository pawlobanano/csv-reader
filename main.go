package main

import (
	"encoding/csv"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	start := time.Now()
	filePath := "customers.csv"
	emailDomains, err := getEmailDomains(filePath)
	if err != nil {
		log.Error("Error reading the domains.", err)
		return
	}

	sortedDomains := sortEmailDomains(emailDomains)

	for _, domain := range sortedDomains {
		log.Info("Sorted domain", "domain", domain, "occurrences", emailDomains[domain])
	}

	log.Info("Program finished.", slog.String("time_taken_ms", time.Since(start).String()))
}

func getEmailDomains(filePath string) (map[string]int, error) {
	emailDomains := make(map[string]int)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip the header line.
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// Read each line and count email domains.
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		email := record[2]
		domain := extractDomain(email)
		emailDomains[domain]++
	}

	return emailDomains, nil
}

func sortEmailDomains(emailDomains map[string]int) []string {
	var sortedDomains []string
	for domain := range emailDomains {
		sortedDomains = append(sortedDomains, domain)
	}

	sort.Strings(sortedDomains)

	return sortedDomains
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "" // Invalid email, return an empty string.
	}

	return parts[1]
}
