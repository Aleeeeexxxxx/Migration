package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	batchSize    = 1000
	generatedMsg = "generated"

	defaultOutput = "test.csv"
)

func main() {
	// flags
	output := flag.String("o", defaultOutput, "output file name")
	n := flag.Int("n", batchSize, "number of models to generate")

	flag.Parse()

	file, err := os.Create(*output)
	if err != nil {
		log.Fatal("failed to create csv file, %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// generate and write in batch
	begin := time.Now()

	records := make([][]string, 0, batchSize)
	for remaining := *n; remaining > 0; remaining -= batchSize {
		now := fmt.Sprintf("%d", time.Now().Unix())
		total := batchSize
		if remaining < batchSize {
			total = remaining
		}

		for i := 0; i < total; i++ {
			records = append(records, []string{
				uuid.New().String(), generatedMsg, now, // ID, MSG, UPDATED_AT
			})
		}

		if err := writer.WriteAll(records); err != nil {
			log.Fatal("failed to write to csv file, %w", err)
		}
		records = records[:0]
	}

	log.Printf("generated %d records successfully, elapsed: %s", *n, time.Since(begin).String())
}
