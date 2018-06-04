package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ryanbradynd05/go-tmdb"
)

var tpl *template.Template

//Shows Struct for tv Shows
type Shows struct {
	Name       string
	ImageLinks string
}

type allShows []Shows
type myShows []Shows

//Data Struct that holds data to be rendered into template
type Data struct {
	Allshows allShows
	Myshows  myShows
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

}

func main() {

	http.HandleFunc("/", indexHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tmdb1 := tmdb.Init("5beb2bf1821813990f328d5c98b5fdc5")

	var tvInfo *tmdb.TvPagedResults
	var err error
	var tvResults []tmdb.TvShort
	option := make(map[string]string)

	totPages := 1

	for i := 1; i <= totPages; i++ {
		option["page"] = strconv.Itoa(i)
		tvInfo, err = tmdb1.GetTvAiringToday(option)
		if err != nil {
			log.Fatal(err)
		}
		totPages = tvInfo.TotalPages
		tvResults = append(tvResults, tvInfo.Results...)
	}

	lines, err := readLines("showlist.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	var allshows allShows
	var myshows myShows

	for _, result := range tvResults {
		for _, line := range lines {
			if line == result.Name {
				myshows = append(myshows, Shows{
					line,
					"https://image.tmdb.org/t/p/original" + result.PosterPath,
				})
			}
		}

		allshows = append(allshows, Shows{
			result.Name,
			"https://image.tmdb.org/t/p/original" + result.PosterPath,
		})
	}

	data := &Data{allshows, myshows}

	tplerr := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if tplerr != nil {
		log.Fatalln(tplerr)
	}

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
