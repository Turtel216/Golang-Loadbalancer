package util

import (
	"os"
)

// parse the config file
func Config_parser(path string) (*[]string, error) {
	//read from config file
	_bytes, err := read_config(path)
	if err != nil {
		return nil, err
	}

	// convert byte slice to string
	conf_str := string(_bytes)

	//keeps track of the url of each iteratino
	var url string
	//stores all urls
	var config []string

	//iterate over string
	for _, _char := range conf_str {
		// add url to struct
		if _char == '\n' || _char == ' ' || _char == '\t' {
			config = append(config, url)
			url = ""
			continue
		}

		// append char to url
		url += string(_char)
	}

	return &config, nil
}

// Read from the config file
func read_config(path string) ([]byte, error) {
	//Read from the file
	buff, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Return read buffer
	return buff, nil
}
