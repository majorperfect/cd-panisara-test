package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

func Scanner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dir string = "build/"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name(), file.IsDir())
		}
		fmt.Println(runtime.NumCPU())
		// runtime.GOMAXPROCS(runtime.NumCPU())
		// f, err := os.Open("file.txt")
		// if err != nil {
		// 	log.Fatalf("unable to read file: %v", err)
		// }
		// defer f.Close()
		// buf := make([]byte, 1024)
		// for {
		// 	n, err := f.Read(buf)
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		continue
		// 	}
		// 	if n > 0 {
		// 		fmt.Println(string(buf[:n]))
		// 	}
		// }
		w.WriteHeader(http.StatusOK)

	}
}
