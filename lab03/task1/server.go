package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("У команды должен быть один аргумент - порт")
		return
	}

	port := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path[1:]
		content, err := ioutil.ReadFile("../files/" + fileName)
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
}
