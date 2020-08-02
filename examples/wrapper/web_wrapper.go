package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	// fix url
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello go")
	})

	// post (get) text
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		// -- read body ---
		// https://qiita.com/futoase/items/ea86b750bbb36d7d859a

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		//rBody := "answer-" + body;
		rBody := execProc("wc", body)

		fmt.Fprint(w, rBody)
	})

	// post (get)
	http.HandleFunc("/reflect", func(w http.ResponseWriter, r *http.Request) {
		// -- read body ---
		// https://qiita.com/futoase/items/ea86b750bbb36d7d859a

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		//rBody := "answer-" + body;
		rBody := execProc("./bin/reflect", body)

		fmt.Fprint(w, rBody)
	})

	// post (get)
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		// -- read body ---
		// https://qiita.com/futoase/items/ea86b750bbb36d7d859a

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		//rBody := "answer-" + body;
		rBody := execProc("./save-to-disk", body)

		fmt.Fprint(w, rBody)
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		// -- read body ---
		// https://qiita.com/futoase/items/ea86b750bbb36d7d859a

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		//rBody := "answer-" + body;
		rBody := execProc("./play-from-disk", body)

		fmt.Fprint(w, rBody)
	})

	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		// -- read body ---
		// https://qiita.com/futoase/items/ea86b750bbb36d7d859a

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		//rBody := "answer-" + body;
		//rBody := execProc("./bin/broadcast", body)

		startProc("./bin/broadcast", body)

		fmt.Fprint(w, "proc started")
	})

	// 静的ファイル配信.
	// ディレクトリ名をURLパスに使う場合
	// 例：http://localhost:8080/www/index.html
	//http.Handle("/www/", http.FileServer(http.Dir("./"))) // OK. serve ./www/* as http://localhost/www/*

	http.Handle("/", http.FileServer(http.Dir("./www"))) // OK. server . /www/*  as http://localhost/*

	// ディレクトリ名とURLパスを変える場合
	// 例：http://localhost:8080/mysecret/sample1.txt
	//http.Handle("/mysecret/", http.StripPrefix("/mysecret/", http.FileServer(http.Dir("./contents"))))

	// 3000ポートで起動
	fmt.Println("start server prot 3000")
	http.ListenAndServe(":3000", nil)
}

func execProc(prog string, arg string) string {
	fmt.Println("call exeProc. prog=" + prog + ", arg=" + arg)
	//cmd := exec.Command("wc")
	cmd := exec.Command(prog)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err2 := cmd.StdoutPipe()
	if err2 != nil {
		log.Fatal(err2)
	}
	stderr, err3 := cmd.StderrPipe();
	if err3 != nil {
		log.Fatal(err3)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("exeProc. start()")
	fmt.Fprintln(stdin, arg)
	fmt.Println("exeProc. write stdin")
	stdin.Close()

	var result string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		result = scanner.Text()
		fmt.Println("exeProc. result=" + result)
		if result != "" {
			break
		}
	}

	var result2 string
	scanner2 := bufio.NewScanner(stderr)
	for scanner2.Scan() {
		result2 = scanner2.Text()
		fmt.Println("exeProc. stderr=" + result2)
		if result2 == "-- send asnwer ---" {
			break
		}
	}

	fmt.Println("--- exeProc. return ----")
	return result
}

func startProc(prog string, arg string) *exec.Cmd {
	fmt.Println("call exeProc. prog=" + prog + ", arg=" + arg)
	//cmd := exec.Command("wc")
	cmd := exec.Command(prog)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("exeProc. start()")
	return cmd
}
