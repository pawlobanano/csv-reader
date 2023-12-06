// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Run opens CSV file, prepare a CSV file reader, process email domains and count the occurences and sort email domains by name.
func Run(log Logger, config *Config) error {
	file, err := os.Open(config.InputCSVFilePathDefault)
	if err != nil {
		log.Warn("Error opening CSV file.", err)
		return err
	}
	defer file.Close()

	reader, err := createCSVfileReader(log, config, file)
	if err != nil {
		return err
	}

	emailDomains := processEmailDomainsConcurrently(log, config, reader)

	sortedDomains := sortEmailDomains(emailDomains)
	for _, domain := range sortedDomains {
		log.Info("Sorted domain.", "domain_name", domain, "occurrences", emailDomains[domain])
	}

	return nil
}

// processEmailDomainsConcurrently processes email domains concurrently using worker goroutines.
// It takes a logger, configuration, and a CSV reader as input, and returns a map of email domains with their occurrences.
// The function utilizes goroutines and channels to achieve concurrent processing.
func processEmailDomainsConcurrently(log Logger, config *Config, reader *csv.Reader) map[string]int {
	var (
		emailDomains = make(map[string]int)
		emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		domainRegex  = regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		wg           sync.WaitGroup
		tasks        = make(chan Task, config.Concurrency)
		results      = make(chan DomainCounter, config.Concurrency)
		errors       = make(chan error, config.Concurrency)
	)

	// Start worker goroutines.
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range tasks {
				customer := parseCustomer(task.record)
				domain := extractDomain(customer.Email)

				// Validate email.
				if !emailRegex.MatchString(customer.Email) {
					errors <- fmt.Errorf("invalid email format: %s", customer.Email)
					continue
				}

				// Validate domain.
				if domain == "" || !domainRegex.MatchString(domain) {
					errors <- fmt.Errorf("invalid domain: %s", domain)
					continue
				}

				results <- DomainCounter{domain, 1}
			}
		}()
	}

	// Start a goroutine to close the results and errors channels when all workers are done.
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Start a goroutine to feed tasks to the workers.
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("The reader failed while reading the file.", err)
				continue
			}

			tasks <- Task{record}
		}

		close(tasks)
	}()

	// Collect results and handle errors from workers.
	for {
		select {
		case result, ok := <-results:
			if !ok { // Results channel closed, no more results to process.
				return emailDomains
			}
			emailDomains[result.domain] += result.counter

		case err, ok := <-errors:
			if !ok { // Errors channel closed, no more errors to process.
				return emailDomains
			}
			log.Warn("Error processing email domain.", err)
		}
	}
}

// createCSVfileReader sets and use buffered reader from bufio package. It returns a csvReader ready to be used for CSV file processing.
func createCSVfileReader(log Logger, config *Config, file *os.File) (*csv.Reader, error) {
	reader := bufio.NewReaderSize(file, config.ReadBufferSizeInBytes)
	csvReader := csv.NewReader(reader)

	// Skip the header line.
	_, err := csvReader.Read()
	if err != nil && err != io.EOF {
		log.Warn("Skipping the first line in the file failed.", err)
		return nil, err
	}

	return csvReader, nil
}

// parseCustomer parses record input to Customer struct for better visibility and maintability of the code.
func parseCustomer(record []string) *Customer {
	return &Customer{
		FirstName: record[0],
		LastName:  record[1],
		Email:     record[2],
		Gender:    record[3],
		IPAddress: record[4],
	}
}

// sortEmailDomains sorts map of email domains input.
func sortEmailDomains(emailDomains map[string]int) []string {
	var sortedDomains []string
	for domain := range emailDomains {
		sortedDomains = append(sortedDomains, domain)
	}

	sort.Strings(sortedDomains)

	return sortedDomains
}

// extractDomain extracts the domain from the email input. Returns empty string if the email is invalid.
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 { // Invalid email.
		return ""
	}

	return parts[1]
}
