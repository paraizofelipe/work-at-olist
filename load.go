package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"work-at-olist/storage"
)

func main() {
	const (
		baseUrl = "localhost"
		path    = "api/records/"
		port    = 8989
	)

	src := "99988526423"
	dst := "9933468278"

	var bodyReq [16]storage.Record

	bodyReq[0] = storage.Record{Source: src, Destination: dst, CallId: 70, Type: "start", Timestamp: "2016-02-29T12:00:00Z"}
	bodyReq[1] = storage.Record{Source: src, Destination: dst, CallId: 70, Type: "end", Timestamp: "2016-02-29T14:00:00Z"}

	bodyReq[2] = storage.Record{Source: src, Destination: dst, CallId: 71, Type: "start", Timestamp: "2017-12-11T15:07:13Z"}
	bodyReq[3] = storage.Record{Source: src, Destination: dst, CallId: 71, Type: "end", Timestamp: "2017-12-11T15:14:56Z"}

	bodyReq[4] = storage.Record{Source: src, Destination: dst, CallId: 72, Type: "start", Timestamp: "2017-12-12T22:47:56Z"}
	bodyReq[5] = storage.Record{Source: src, Destination: dst, CallId: 72, Type: "end", Timestamp: "2017-12-12T22:50:56Z"}

	bodyReq[6] = storage.Record{Source: src, Destination: dst, CallId: 73, Type: "start", Timestamp: "2017-12-12T21:57:13Z"}
	bodyReq[7] = storage.Record{Source: src, Destination: dst, CallId: 73, Type: "end", Timestamp: "2017-12-12T22:10:56Z"}

	bodyReq[8] = storage.Record{Source: src, Destination: dst, CallId: 74, Type: "start", Timestamp: "2017-12-12T04:57:13Z"}
	bodyReq[9] = storage.Record{Source: src, Destination: dst, CallId: 74, Type: "end", Timestamp: "2017-12-12T06:10:56Z"}

	bodyReq[10] = storage.Record{Source: src, Destination: dst, CallId: 75, Type: "start", Timestamp: "2017-12-13T21:57:13Z"}
	bodyReq[11] = storage.Record{Source: src, Destination: dst, CallId: 75, Type: "end", Timestamp: "2017-12-14T22:10:56Z"}

	bodyReq[12] = storage.Record{Source: src, Destination: dst, CallId: 76, Type: "start", Timestamp: "2017-12-12T15:07:58Z"}
	bodyReq[13] = storage.Record{Source: src, Destination: dst, CallId: 76, Type: "end", Timestamp: "2017-12-12T15:12:56Z"}

	bodyReq[14] = storage.Record{Source: src, Destination: dst, CallId: 77, Type: "start", Timestamp: "2018-02-28T21:57:13Z"}
	bodyReq[15] = storage.Record{Source: src, Destination: dst, CallId: 77, Type: "end", Timestamp: "2018-03-01T22:10:56Z"}

	url := fmt.Sprintf("http://%s:%d/%s", baseUrl, port, path)

	for _, body := range bodyReq {
		j, err := json.Marshal(body)
		if err != nil {
			continue
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}

}
