package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// inspired by https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f

func discoverIp() string  {

	discoverIpurl := "https://api.ipify.org?format=text"
	response, err := http.Get(discoverIpurl)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		//let's close the Body and manage the error if needed
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()


	currentPublicIp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("your current IP address is: %s\n", currentPublicIp)

	return string(currentPublicIp)
}
