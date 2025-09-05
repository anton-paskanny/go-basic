package main

import (
	"flag"
	"fmt"
	"log"

	"demo/bin/api"
	"demo/bin/config"
	f "demo/bin/file"
)

func main() {
	// Flags
	opCreate := flag.Bool("create", false, "Create a new bin from JSON file")
	opGet := flag.Bool("get", false, "Get a bin by id")
	opUpdate := flag.Bool("update", false, "Update a bin by id with JSON file")
	opDelete := flag.Bool("delete", false, "Delete a bin by id")
	id := flag.String("id", "", "Bin id for get/update")
	filePath := flag.String("file", "", "Path to JSON file for create/update")
	isPrivate := flag.Bool("private", false, "Create bin as private")
	flag.Parse()

	// Ensure exactly one operation
	ops := 0
	if *opCreate {
		ops++
	}
	if *opGet {
		ops++
	}
	if *opUpdate {
		ops++
	}
	if *opDelete {
		ops++
	}
	if ops != 1 {
		fmt.Println("Usage:")
		fmt.Println("  -create -file path [-private]")
		fmt.Println("  -get -id BIN_ID")
		fmt.Println("  -update -id BIN_ID -file path")
		fmt.Println("  -delete -id BIN_ID")
		flag.PrintDefaults()
		return
	}

	// Load configuration and init API client
	cfg := config.Load()
	apiClient := api.New(cfg)

	switch {
	case *opCreate:
		if *filePath == "" {
			log.Fatal("-file is required for -create")
		}
		data, err := f.ReadJSONFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		newID, err := apiClient.CreateBin(data, *isPrivate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created bin id: %s\n", newID)

	case *opGet:
		if *id == "" {
			log.Fatal("-id is required for -get")
		}
		data, err := apiClient.GetBin(*id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))

	case *opUpdate:
		if *id == "" {
			log.Fatal("-id is required for -update")
		}
		if *filePath == "" {
			log.Fatal("-file is required for -update")
		}
		data, err := f.ReadJSONFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		if err := apiClient.UpdateBin(*id, data); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Updated")

	case *opDelete:
		if *id == "" {
			log.Fatal("-id is required for -delete")
		}
		if err := apiClient.DeleteBin(*id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Deleted")
	}
}
