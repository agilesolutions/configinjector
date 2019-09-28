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

	param1 := flag.String("url", "http://spring.config.com:8888/cms/default/PRD", "Application yaml file.")

	param2 := flag.String("directory", "", "Directory to download configs")

	flag.Parse()

	uri := *param1
	dir := *param2

	fmt.Printf("Application config : %s processed at : %s\n", uri, dir)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// HTTP pull spring boot config from Spring Cloud Config server
	//if response, err := http.Get(uri); err == nil {
	fmt.Printf("downloaded %s.\n", uri)

	//defer response.Body.Close()

	//nbody, _ := ioutil.ReadAll(response.Body)

	var m interface{}

	error := json.Unmarshal([]byte(mock), &m)
	if error != nil {
		panic(error)
	}
	/*
	* write a new property file aggregating content of all spring boot yaml files
	* This file gets pulled in on a k8s configmap GO template on a HELM chart
	*/
	f, ferror := os.Create(dir + "/application.properties")
	if ferror != nil {
		panic(ferror)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	mu := m.(map[string]interface{})
	// parse out all property key value pairs
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
							//fmt.Println("properties found")
							values := z.(map[string]interface{})
							// write propery key value pair to new property file on k8s configmap
							for propertyName, propertyValue := range values {
								fmt.Fprintln(w, propertyName, "=", propertyValue)
								fmt.Println(propertyName, "=", propertyValue)
								w.Flush()
							}

						}

					}

				}

			}

		}

	}

	//} else {
	//	panic(err)
	//}

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
   "name":"cms",
   "profiles":[  
      "default"
   ],
   "label":"PRD",
   "version":"0322baa10d3f1d1a9452dcaab16ec6b95fb8225c",
   "state":null,
   "propertySources":[  
      {  
         "name":"https://cloud.config.com/scm/id/config.git/cms.yml",
         "source":{  
            "server.tomcat.max-http-header-size":131072,
            "server.max-http-header-size":10000000,
            "spring.http.multipart.max-file-size":"20MB",
            "spring.http.multipart.max-request-size":"20MB",
            "local.server.port":8780,
            "datasource.jdbc.url":"jdbc:oracle:thin:www.cms.com:9999:cms001",
            "datasource.jdbc.driverClassName":"oracle.jdbc.xa.client.OracleXADataSource",
            "datasource.jdbc.userName":"xxx",
            "datasource.jdbc.password":"xxx",
            "datasource.pool.initialSize":2,
            "datasource.pool.maxIdle":5,
            "datasource.pool.maxActive":20,
            "activeDirectory.url":"ldap://xxxx.compute.amazonaws.com:389",
            "activeDirectory.domain":"xxx.net",
            "activeDirectory.search-filter":"((&objectClass=user)(userPrincipalName={0})(memberof=CN=Users,DC=pub,DC=net))",
            "email.smtp.host":"smtppublic.mail.com",
            "email.smtp.port":25,
            "email.smtp.username":"xxx",
            "email.smtp.password":"test",
            "application.security.jwt.secret":"lksadfuhsdaflöweuuus",
            "application.security.jwt.validityMinutes":60,
            "application.security.allowedOrigins[0]":"http://localhost:4200"
         }
      }
   ]
}
`

