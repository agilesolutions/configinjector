package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*

This module is designed to run as a docker container at a Jenkins pipeline agent to fetch configs from Spring Cloud config server and inject the Spring boot application.yml on a k8s configmap.
1. REST fetches application yaml from config server
2. puts it on a helm chart where helm pulls it on a configmap

*/
func main() {

	exitCode := 0

	param1 := flag.String("url", "https://github.com/o7planning/spring-cloud-config-git-repo-example/blob/master/app-about-company.properties", "Application yaml file.")

	param2 := flag.String("directory", "", "Directory to download configs")

	flag.Parse()

	uri := *param1
	dir := *param2

	fmt.Printf("Application config url : %s downloaded at : %s\n", uri, dir)

	fmt.Printf("DownloadToFile From: %s.\n", uri)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if response, err := http.Get(uri); err == nil {
		fmt.Printf("downloaded %s.\n", uri)

		defer response.Body.Close()

	nbody,_ := ioutil.ReadAll(response.Body)

	var m interface{}

	error := json.Unmarshal([]byte(nbody), &m)
	if error != nil {
		panic(error)
	}

	f, ferror := os.Create("application.properties")
	if ferror != nil {
		panic(ferror)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	mu := m.(map[string]interface{})

	for k, v := range mu {

		switch vv := v.(type) {

		case []interface{}:

			if k == "propertySources" {

				//fmt.Println(k, "PROPERTY SOURCES FOUND :", v)

				for _, u := range vv {

					//fmt.Println(u)

					source := u.(map[string]interface{})

					for k, z := range source {

						if k == "source" {
							fmt.Println("properties found")
							values := z.(map[string]interface{})

							for propertyName, propertyValue := range values {
								fmt.Println(propertyName, "=", propertyValue)
								fmt.Fprintf(w, "%s\n", propertyName)
								w.Flush()
							}
						}

					}

				}

			}

		}

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

const mock = `

{

   "name":"irs",

   "profiles":[

      "default"

   ],

   "label":"PRD",

   "version":"0f8c9d17e9435b2072fab97356944e3155ad22b5",

   "state":null,

   "propertySources":[

      {

         "name":"https://mycomp.com/scm/sfs/config.git/irs.yml",

         "source":{

            "AAAAAA":9999,
            "XXXXXX":"1.4"

         }

      }

   ]

}

`