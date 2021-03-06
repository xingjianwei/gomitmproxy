package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

func httpDump(req *http.Request, resp *http.Response) {
	defer resp.Body.Close()
	var respStatusStr string
	respStatus := resp.StatusCode
	respStatusHeader := int(math.Floor(float64(respStatus / 100)))
	switch respStatusHeader {
	case 2:
		respStatusStr = Green("<--" + strconv.Itoa(respStatus))
	case 3:
		respStatusStr = Yellow("<--" + strconv.Itoa(respStatus))
	case 4:
		respStatusStr = Magenta("<--" + strconv.Itoa(respStatus))
	case 5:
		respStatusStr = Red("<--" + strconv.Itoa(respStatus))
	}
	fmt.Println(Green("Request:"))
	fmt.Printf("%s %s %s\n", Blue(req.Method), req.RequestURI, respStatusStr)
	for headerName, headerContext := range req.Header {
		fmt.Printf("%s: %s\n", Blue(headerName), headerContext)
	}
	fmt.Println(Green("Response:"))
	for headerName, headerContext := range resp.Header {
		fmt.Printf("%s: %s\n", Blue(headerName), headerContext)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("func httpDump read resp body err:", err)
	} else {
		acceptEncode := resp.Header["Content-Encoding"]
		var respBodyBin bytes.Buffer
		w := bufio.NewWriter(&respBodyBin)
		w.Write(respBody)
		w.Flush()
		for _, compress := range acceptEncode {
			switch compress {
			case "gzip":
				r, err := gzip.NewReader(&respBodyBin)
				if err != nil {
					logger.Println("gzip reader err:", err)
				} else {
					defer r.Close()
					respBody, _ = ioutil.ReadAll(r)
				}
				break
			case "deflate":
				r := flate.NewReader(&respBodyBin)
				defer r.Close()
				respBody, _ = ioutil.ReadAll(r)
				break
			}
		}
		fmt.Printf("%s\n", string(respBody))
	}

	fmt.Printf("%s%s%s\n", Black("####################"), Cyan("END"), Black("####################"))
}
