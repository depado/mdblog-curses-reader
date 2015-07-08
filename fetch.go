package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	articleURL = "http://markdownblog.com/api/v1/article"
)

type articleAPIType struct {
	NumResults int `json:"num_results"`
	Objects    []struct {
		Content   string `json:"content"`
		ID        int    `json:"id"`
		PubDate   string `json:"pub_date"`
		Title     string `json:"title"`
		TitleSlug string `json:"title_slug"`
		User      struct {
			BlogBg          string `json:"blog_bg"`
			BlogDescription string `json:"blog_description"`
			BlogImage       string `json:"blog_image"`
			BlogPublic      string `json:"blog_public"`
			BlogSlug        string `json:"blog_slug"`
			BlogTitle       string `json:"blog_title"`
			ID              int    `json:"id"`
			Username        string `json:"username"`
		} `json:"user"`
		UserID int `json:"user_id"`
	} `json:"objects"`
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
}

type articleType struct {
	date      string
	title     string
	titleSlug string
	url       string
}

func dlSingle(url string) (content articleAPIType, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		return
	}
	return
}

func fetchAllArticles() (content map[string][]articleType, err error) {
	content = make(map[string][]articleType)
	first, err := dlSingle(articleURL)
	if err != nil {
		return
	}
	tp := first.TotalPages
	for _, item := range first.Objects {
		content[item.User.BlogSlug] = append(content[item.User.BlogSlug], articleType{
			date:      item.PubDate,
			title:     item.Title,
			titleSlug: item.TitleSlug,
			url:       "http://" + item.User.BlogSlug + ".markdownblog.com/" + item.TitleSlug,
		})
	}
	for i := 2; i < tp; i++ {
		current, err := dlSingle(articleURL + "?page=" + strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		for _, item := range current.Objects {
			content[item.User.BlogSlug] = append(content[item.User.BlogSlug], articleType{
				date:      item.PubDate,
				title:     item.Title,
				titleSlug: item.TitleSlug,
				url:       "http://" + item.User.BlogSlug + ".markdownblog.com/" + item.TitleSlug,
			})
		}
	}
	return
}
