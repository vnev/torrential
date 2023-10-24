package main

import (
	// Uncomment this line to pass the first stage
	// "encoding/json"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if unicode.IsLetter(rune(bencodedString[0])) {
		switch bencodedString[0] {
		case 'i':
			{
				result, err := strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
				if err != nil {
					return "", fmt.Errorf("Error parsing number")
				}

				return result, nil
			}
		case 'l':
			{
				finalResult := []interface{}{}
				start := 1
				end := len(bencodedString) - 1
				for start != end {
					result, err := decodeBencode(bencodedString[start:end])
					if err != nil {
						return "", fmt.Errorf("Encountered issue while trying to decode list")
					}

					start += len(fmt.Sprintf("%v", result)) + 2 // chomp the string length and : for string decode, i and ending e for number decode
					finalResult = append(finalResult, result)
				}

				return finalResult, nil
			}
		case 'd':
			{
				finalResult := make(map[string]interface{})
				start := 1
				end := len(bencodedString) - 1

				for start != end {
					// decode the key
					result, err := decodeBencode(bencodedString[start:end])
					if err != nil {
						return "", fmt.Errorf("Encountered issue while trying to decode key in dict")
					}
					tempStr := fmt.Sprintf("%v", result)
					start += len(tempStr) + 2

					// decode the value
					result, err = decodeBencode(bencodedString[start:end])
					if err != nil {
						return "", fmt.Errorf("Encountered issue while trying to decode value in dict")
					}

					start += len(fmt.Sprintf("%v", result)) + 2
					finalResult[tempStr] = result
				}

				return finalResult, nil
			}
		default:
			return "", fmt.Errorf("Only supporting numbers")
		}
	}

	return "", fmt.Errorf("Only numbers and strings are supported at the moment")
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	//fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		//Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
