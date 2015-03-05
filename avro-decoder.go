package main

import (
	"fmt"
	"log"

	"github.com/stealthly/go-avro/decoder"
)

type Record struct {
	Id        int64
	Subdomain string
}

func (r Record) Print() {
	fmt.Printf("Account ID: %d, Subdomain: %s\n", r.Id, r.Subdomain)
}

func main() {
	datumReader := decoder.NewGenericDatumReader()
	reader, err := decoder.NewDataFileReader("data.avro", datumReader)
	if err != nil {
		log.Panic("Unexpected error reading file", err)
	}
	record := Record{}
	_ = reader.Next(&record)
	record.Print()
}
