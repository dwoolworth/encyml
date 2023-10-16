package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type DataStruct struct {
	Data map[string]string `yaml:"data"`
}

func main() {
	// Flags
	outputFlag := flag.Bool("o", false, "Output values in various states")
	encodeFlag := flag.Bool("e", false, "Base64 encode the values in the YAML file")
	decodeFlag := flag.Bool("d", false, "Base64 decode the values in the YAML file")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: go run program.go [flags] <filename>")
		return
	}

	filename := flag.Arg(0)

	// Read the YAML file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var data DataStruct
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	for key, value := range data.Data {
		original := value

		var decodedValue string
		if *decodeFlag || (!*decodeFlag && !*encodeFlag) {
			decoded, _ := base64.StdEncoding.DecodeString(value)
			decodedValue = string(decoded)
		} else {
			decodedValue = original
		}

		trimmed := strings.TrimSpace(decodedValue) // trimmed state

		if *outputFlag {
			fmt.Printf("     Key: %s\n", key)
			fmt.Printf("Original: %s\n", original)
			if *decodeFlag || (!*decodeFlag && !*encodeFlag) {
			  fmt.Printf(" Decoded: %s\n", decodedValue)
			}
			fmt.Printf(" Trimmed: %s\n", trimmed)
		}

		if *encodeFlag || (!*decodeFlag && !*encodeFlag) {
			encoded := base64.StdEncoding.EncodeToString([]byte(trimmed))
			data.Data[key] = encoded
			if *outputFlag {
        fmt.Printf(" Encoded: %s\n\n", encoded)
			}
		}

		if *decodeFlag && !*encodeFlag {
			data.Data[key] = trimmed
			if *outputFlag {
        fmt.Printf("  Output: %s\n\n", trimmed)
			}
		}
	}

	// Convert the modified data back to YAML format
	output, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Write the YAML back to the original file
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("File processed successfully!")
}

