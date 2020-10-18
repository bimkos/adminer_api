package main

import (
	"net/http"
	"log"
	"net/url"
	"net/http/cookiejar"
	"io/ioutil"
	"flag"
	"github.com/anaskhan96/soup"
	"strings"
)

func main() {
	var (
		adminerUrl string
		adminerPassword string
		adminerUser string
		adminerServer string 
		adminerDB string
		adminerExport string
	) 

	// Args
	flag.StringVar(&adminerUrl, "url", "", "adminer url")
	flag.StringVar(&adminerPassword, "pass", "", "user password")
	flag.StringVar(&adminerUser, "user", "", "username")
	flag.StringVar(&adminerServer, "host", "", "DB host")
	flag.StringVar(&adminerDB, "db", "", "DB name")
	flag.StringVar(&adminerExport, "export", "", "if empty - export all")
	
	flag.Parse()

	if len(strings.TrimSpace(adminerUrl)) != 0 && len(strings.TrimSpace(adminerUser)) != 0 && len(strings.TrimSpace(adminerPassword)) != 0 {
		client := createClient()
		// Login
		login(adminerUrl, adminerServer, adminerUser, adminerPassword, adminerDB, client)
		export(adminerUrl, adminerUser, client)
	}
}

func createClient() (client *http.Client) {
	// Create http.Client with cookies
	jar, err := cookiejar.New(nil)
	if err != nil { 
		log.Fatal(err)
	}
	client = &http.Client{
		Jar: jar,
	}
	return
}

func login(adminerUrl string, server string, username string, password string, db string, client *http.Client) {
	resp, err := client.PostForm(adminerUrl, url.Values {
		"auth[driver]"   : {"server"},
        "auth[server]"   : {server},
        "auth[username]" : {username},
        "auth[password]" : {password},
        "auth[db]"       : {db},
	})
	if err != nil { 
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func export(adminerUrl string, adminerUser string, client *http.Client) {
	resp, err := client.Get(adminerUrl + "?username=" + adminerUser + "&dump=")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	html := soup.HTMLParse(string(body))
	token := html.Find("input", "name", "token")
	log.Print(token.Attrs()["value"])
}