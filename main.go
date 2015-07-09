package main

import (
	"fmt"
	"log"
)

func perror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	article, err := ncurses()
	perror(err)
	content, err := dlContent(article.url + "/ansi")
	perror(err)
	fmt.Println(content)
}
