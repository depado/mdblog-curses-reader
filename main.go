package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func perror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	var reader *bufio.Reader
	var articles map[string][]articleType
	var article articleType

	reader = bufio.NewReader(os.Stdin)

	articles, err = fetchAllArticles()
	if err != nil {
		return
	}
	for {
		article, err = ncurses(articles)
		perror(err)
		if article == (articleType{}) {
			return
		}
		content, err := dlContent(article.url + "/ansi")
		perror(err)
		fmt.Println(content)
		fmt.Println("Press Enter to go back to the menu.")
		reader.ReadString('\n')
	}
}
