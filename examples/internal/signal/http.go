package signal

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"os"

)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer() chan string {
	port := flag.Int("port", 8080, "http server port")
	flag.Parse()

	sdpChan := make(chan string)
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		//fmt.Fprintf(w, "done")
		sdpChan <- string(body)

		// ---- try dual way ---
		answer := <-sdpChan
		fmt.Fprintln(os.Stderr, "-- read asnwer from channel ---")
		fmt.Fprintln(os.Stderr, answer)
		fmt.Fprintln(os.Stderr, "-- send response ---")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type")
		fmt.Fprint(w, answer)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return sdpChan
}
