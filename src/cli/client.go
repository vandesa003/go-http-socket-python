package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	url = "http://0.0.0.0:8080"
)

// http request struct
type HttpRequest struct {
	image string `json:"image" example:"abcdefg"`
}

func main() {
	fmt.Println("welcome cli.")

	// create http client, ref: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// call server by get method.

	// resp, err := netClient.Get(url)
	// if err != nil {
	// 	log.Error().Err(err).Msg("get method faild connect server.")
	// 	return
	// }
	// // read content body.
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)

	// // print results.
	// fmt.Println(resp.StatusCode)
	// fmt.Println(string(body))

	// call server by post method.
	b64 := imgToBase64("data/test.png")
	fmt.Println(b64)
	payload := map[string]string{"image": b64}
	// payload := HttpRequest{image: b64}
	json_data, err := json.Marshal(payload)
	resp, err := netClient.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Error().Err(err).Msg("post method failed connect server.")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))

}

// convert bytes to base64 strings
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// convert local image file to base64 strings
func imgToBase64(p string) string {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error().Err(err).Msg("error in convert local image to base64.")
	}
	return toBase64(b)
}
