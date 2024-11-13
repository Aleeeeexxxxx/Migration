package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"migration/src"
	"migration/src/cmd/loader/producer"
)

const (
	batchSize      = 1000
	maxWorker      = 10
	defaultBufSize = 2 * 1024 * 1024 // 2MB
	sqlInsertMany  = "INSERT INTO models ( ID, MSG, UPDATED_AT ) VALUES "
)

func consumer(ch chan producer.Rows, wg *sync.WaitGroup, cfg src.DBConfig, pool *sync.Pool) {
	wg.Add(1)
	defer wg.Done()

	db, err := cfg.Dial()
	if err != nil {
		log.Fatal("fail to connect to db", err)
	}

	for {
		records, ok := <-ch
		if !ok {
			return
		}

		if len(records) > 0 {
			sql := sqlInsertMany + records.ToValues() + ";"
			if result := db.Exec(sql); result.Error != nil {
				log.Fatal("fail to insert data to mysql", result.Error)
			}
		}

		pool.Put(records[:0])
	}
}

func main() {
	username := flag.String("username", "origin", "Username to login mysql")
	password := flag.String("password", "origin", "Password to login mysql")
	ip := flag.String("ip", "localhost", "IP address to login mysql")
	port := flag.String("port", "3306", "Port to connect to mysql")
	database := flag.String("database", "migration_origin", "Database to connect to mysql")
	file := flag.String("f", "test.csv", "data file to load")

	flag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		fmt.Printf("fail to open file: %v\n", err)
		return
	}
	defer func() { _ = f.Close() }()

	cur := 0 // how many workers are running now
	cfg := src.DBConfig{
		Username: *username,
		Password: *password,
		IP:       *ip,
		Port:     *port,
		Schema:   *database,
	}
	pool := &sync.Pool{
		New: func() interface{} {
			return make(producer.Rows, 0, batchSize)
		},
	}
	begin := time.Now()
	var wg sync.WaitGroup
	var loadErr error
	ch := make(chan producer.Rows)
	p := producer.NewProducer(f, batchSize, pool, defaultBufSize)
	total := 0

	for end := false; !end; {
		records, err := p.Produce()
		if err != nil {
			end = true
			if err != io.EOF {
				loadErr = err
			}
		}

		total += len(records)
		log.Printf("has inserted %d records\r\n", total)

		select {
		case ch <- records:
		default:
			if cur < maxWorker {
				cur++
				go consumer(ch, &wg, cfg, pool)
			}
			ch <- records
		}
	}

	close(ch)
	wg.Wait()

	log.Printf("elapsed time: %s\n", time.Since(begin))
	if err != nil {
		log.Printf(
			"shutdown due to error loading data from csv file, err: %s",
			loadErr.Error(),
		)
	} else {
		log.Println("all data has been loaded successfully")
	}
}
