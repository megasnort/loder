package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	base_url, offset, limit, err := parseArguments()

	if err != nil {
		fmt.Println(err)
		return
	}

	var url string

	for counter := offset; counter < (offset+limit) || limit == -1; counter++ {
		url = fmt.Sprintf(base_url, counter)

		fmt.Println("Downloading", url)
		_, err := downloadFromUrl(url)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func downloadFromUrl(url string) (fileName string, err error) {

	tokens := strings.Split(url, "/")
	fileName = tokens[len(tokens)-1]

	if _, err := os.Stat(fileName); err == nil {
		return fileName, nil
	}

	response, err := http.Get(url)
	defer response.Body.Close()

	// network problem
	if err != nil {
		return "", err
	}

	// file not found
	if response.StatusCode != 200 {
		return "", errors.New("Not found")
	}

	output, err := os.Create(fileName)
	defer output.Close()

	// problem creating on file system
	if err != nil {
		return "", err
	}

	_, err = io.Copy(output, response.Body)

	// problem copying the file
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func parseArguments() (base_url string, offset int, limit int, err error) {
	args := os.Args[1:]
	var my_err error

	if len(args) == 0 {
		fmt.Println("Provide a valid url. f.e. http://this/is/my/image%d.jpg")
		return
	}

	base_url = args[0]
	offset = 1
	limit = -1

	if len(args) > 1 {
		i, err := strconv.Atoi(args[1])

		if err != nil {
			my_err = errors.New("Invalid argument for [offset]")
		}

		offset = i

		if len(args) > 2 {
			i, err = strconv.Atoi(args[2])

			if err != nil {
				my_err = errors.New("Invalid argument for [limit]")
			}

			limit = i
		}
	}

	return base_url, offset, limit, my_err
}
