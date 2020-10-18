package main

import (
	"flag"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"os"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func main() {
	var (
		// Main values
		adminerUrl      string
		adminerPassword string
		adminerUser     string
		adminerServer   string
		adminerDB       string

		// Export values
		adminerExport              string
		adminerExportOutput        string
		adminerExportFormat        string
		adminerExportDBStyle       string
		adminerExportRoutines      string
		adminerExportEvents        string
		adminerExportTableStyle    string
		adminerExportTriggers      string
		adminerExportDataStyle     string
		adminerExportAutoIncrement string
	)

	// Main args
	flag.StringVar(&adminerUrl, "url", "", "adminer url")
	flag.StringVar(&adminerPassword, "pass", "", "user password")
	flag.StringVar(&adminerUser, "user", "", "username")
	flag.StringVar(&adminerServer, "host", "", "DB host")
	flag.StringVar(&adminerDB, "db", "", "DB name")

	// Export args
	flag.StringVar(&adminerExport, "export", "", "if empty - export all")
	flag.StringVar(&adminerExportOutput, "exportOutput", "save", "save/open/gzip")
	flag.StringVar(&adminerExportFormat, "exportFormat", "sql", "sql/csv/csv;/tsv")
	flag.StringVar(&adminerExportDBStyle, "exportDBStyle", "CREATE", "CREATE/DROP+CREATE/USE")
	flag.StringVar(&adminerExportRoutines, "exportRoutines", "1", "if exmpty - off")
	flag.StringVar(&adminerExportEvents, "exportEvents", "1", "if empty - off")
	flag.StringVar(&adminerExportTableStyle, "exportTableStyle", "DROP+CREATE", "DROP+CREATE/CREATE")
	flag.StringVar(&adminerExportTriggers, "exportTriggers", "1", "if empty - off")
	flag.StringVar(&adminerExportDataStyle, "exportDataStyle", "INSERT", "INSERT/INSERT+UPDATE/TRUNCATE+INSERT")
	flag.StringVar(&adminerExportAutoIncrement, "exportAutoIncrement", "1", "if empty - off")

	flag.Parse()

	if len(strings.TrimSpace(adminerUrl)) != 0 && len(strings.TrimSpace(adminerUser)) != 0 && len(strings.TrimSpace(adminerPassword)) != 0 {
		client := createClient()
		// Login
		login(adminerUrl, adminerServer, adminerUser, adminerPassword, adminerDB, client)
		export(
			adminerUrl,
			adminerUser,
			adminerServer,
			client,
			adminerExportDataStyle,
			adminerExportDBStyle,
			adminerExportEvents,
			adminerExportFormat,
			adminerExportOutput,
			adminerExportRoutines,
			adminerExportTableStyle,
			adminerExportTriggers,
			adminerExportAutoIncrement,
		)
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
	resp, err := client.PostForm(adminerUrl, url.Values{
		"auth[driver]":   {"server"},
		"auth[server]":   {server},
		"auth[username]": {username},
		"auth[password]": {password},
		"auth[db]":       {db},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func export(
	adminerUrl string,
	adminerUser string,
	adminerServer string,
	client *http.Client,
	adminerExportDataStyle string,
	adminerExportDBStyle string,
	adminerExportEvents string,
	adminerExportFormat string,
	adminerExportOutput string,
	adminerExportRoutines string,
	adminerExportTableStyle string,
	adminerExportTriggers string,
	adminerExportAutoIncrement string,
) {
	// Parse token and databases
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
	token := html.Find("input", "name", "token").Attrs()["value"]
	dbs := html.FindAll("input", "name", "databases[]")

	// Create query
	values := url.Values{
		"output":                     {adminerExportOutput},
		"format":                     {adminerExportFormat},
		"db_style":                   {adminerExportDBStyle},
		"routines":                   {adminerExportRoutines},
		"events":                     {adminerExportEvents},
		"table_style":                {adminerExportTableStyle},
		"adminerExportAutoIncrement": {adminerExportAutoIncrement},
		"triggers":                   {adminerExportTriggers},
		"data_style":                 {adminerExportDataStyle},
		"token":                      {token},
	}

	for _, db := range dbs {
		values.Add("databases[]", db.Attrs()["value"])
	}

	resp, err = client.PostForm(adminerUrl+"?server="+adminerServer+"&username="+adminerUser+"&dump=", values)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(adminerUrl)
    if err != nil {
        log.Fatal(err)
    }

	out, err := os.Create(u.Host +  ".txt")
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
