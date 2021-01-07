package main

import (
	"fmt"
	"log"

	"github.com/masahiro331/go-box-parser/vagrantcloud"
)

func main() {
	reader, err := vagrantcloud.GetBox("generic/alpine38", "3.1.20", vagrantcloud.Virtualbox)
	if err != nil {
		log.Fatal(err)
	}
	treader, err := vagrantcloud.NewBoxReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	for {
		header, err := treader.Next()
		if err != nil {
			break
		}
		fmt.Println(header.Name)
	}

}
