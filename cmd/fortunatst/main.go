package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jacoelho/fortuna"
	"github.com/jacoelho/fortuna/entropy/random"
)

func main() {
	f, err := fortuna.New()
	if err != nil {
		log.Fatal(err)
	}

	s := &random.Random{}
	m := fortuna.Manager{Pools: 32}

	m.Accumulate(f, s)

	seedFile := fortuna.NewSeedFile(f, "test.dat")
	b, err := seedFile.Read()
	if err != nil {
		log.Println(err)
	}

	f.SeedBytes(b)

	randomData, err := os.OpenFile("test.dat", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer randomData.Close()

	// wait for enough entropy
	time.Sleep(45 * time.Second)

	_, err = io.CopyN(randomData, f, 1<<25)
	if err != nil {
		log.Fatal(err)
	}
}
