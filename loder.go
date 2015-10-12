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

	c := make(chan int)

	if err != nil {
		fmt.Println(err)
		return
	}

	var url string

	for counter := offset; counter < (offset+limit) || limit == -1; counter++ {
		url = fmt.Sprintf(base_url, counter)
		go downloadFromUrl(url, c)
	}

	for i := 0; i < limit; i++ {
		<-c
	}
}

func downloadFromUrl(url string, c chan int) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	if _, err := os.Stat(fileName); err == nil {
		fmt.Println(fileName, "exists")
		c <- 0
		return
	} else {
		response, err := http.Get(url)
		defer response.Body.Close()

		// network problem
		if err != nil {
			fmt.Println(fileName, "network problem")
			c <- 0
			return
		}

		// file not found
		if response.StatusCode != 200 {
			fmt.Println(fileName, "not found")
			c <- 0
			return
		}

		output, err := os.Create(fileName)
		defer output.Close()

		// problem creating on file system
		if err != nil {
			fmt.Println(fileName, "creating problem")
			c <- 0
			return
		}

		_, err = io.Copy(output, response.Body)

		// problem copying the file
		if err != nil {
			fmt.Println(fileName, "copying problem")
			c <- 0
			return
		}

		fmt.Println("Downloaded", fileName)
	}

	c <- 1
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
		} else {
			// this is a tempfix
			// eventually we want a fixed amount of go routines, fetching all
			// images, until the increasing of numbers hits an 404 image
			limit = 10
		}
	}

	return base_url, offset, limit, my_err
}
