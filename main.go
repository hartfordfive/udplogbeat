package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/hartfordfive/udplogbeat/beater"
)

func main() {
	err := beat.Run("udplogbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
