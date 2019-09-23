package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"log"
	"net/http"
)

/*

This module is designed to run as a docker container at a Jenkins pipeline agent to fetch configs from Spring Cloud config server and inject the Spring boot application.yml on a k8s configmap.


1. REST fetches application yaml from config server
2. puts it on a helm chart where helm pulls it on a configmap
see: https://gist.github.com/jeffjohnson9046/3f18eb0de8c9674347abbd978ba78e6d
 */
func main() {

	exitCode := 0

	param1 := flag.String("url", "http://configuration-service:8888/my-spring-boot-app-dev.yml", "Application yaml file.")

	param2 := flag.Bool("terminate", false, "Terminate jenkins pipeline on violation.")

	flag.Parse()

	uri := *param1
	terminate := *param2

	//uri = "http://configuration-service:8888/my-spring-boot-app-dev.yml"

	fmt.Printf("Application config url : %s termination is : %t\n", uri, terminate)

	fmt.Printf("DownloadToFile From: %s.\n", uri)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if d, err := HTTPDownload(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)
		if WriteFile("chart/application.yml", d) == nil {
			fmt.Printf("saved %s as application.yml\n", uri)
		}
	} else {
		panic(err)
	}



	os.Exit(exitCode)
}


func HTTPDownload(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
