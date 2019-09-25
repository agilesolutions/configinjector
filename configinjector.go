The Go Playground   Imports 
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
97
98
99
100
101
102
103
104
105
106
107
108
109
110
111
112
113
114
115
116
117
118
119
120
121
122
123
124
125
126
127
128
129
130
131
132
133
134
135
136
137
138
139
140
141
142
143
144
145
146
147
148
149
150
151
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

	param1 := flag.String("url", "http://configuration-service:8888/my-spring-boot-app-dev.yml", "Application yaml file.")

	param2 := flag.String("directory", "", "Directory to download configs")

	flag.Parse()

	uri := *param1
	dir := *param2

	fmt.Printf("Application config url : %s downloaded at : %s\n", uri, dir)

	fmt.Printf("DownloadToFile From: %s.\n", uri)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//if response, err := http.Get(uri); err == nil {
	//	fmt.Printf("downloaded %s.\n", uri)

	//	defer response.Body.Close()

	//nbody, err := ioutil.ReadAll(response.Body)

	var m interface{}

	error := json.Unmarshal([]byte(mock), &m)
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

Application config url : http://configuration-service:8888/my-spring-boot-app-dev.yml downloaded at : 
DownloadToFile From: http://configuration-service:8888/my-spring-boot-app-dev.yml.
properties found
AAAAAA = 9999
XXXXXX = 1.4

Program exited.