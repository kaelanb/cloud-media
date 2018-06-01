package main

import (
	"bufio"
	"html/template"
	"log"
	"os"

	"github.com/ryanbradynd05/go-tmdb"
)

var tpl *template.Template

type Shows struct {
	Name       string
	ImageLinks string
}

type AllShows []Shows
type MyShows []Shows

type Data struct {
	Allshows AllShows
	Myshows  MyShows
}

func init() {
	tpl = template.Must(template.ParseFiles("templates/tpl.gohtml"))

}

func main() {
	tmdb1 := tmdb.Init("5beb2bf1821813990f328d5c98b5fdc5")

	//	options := make(map[string]string)
	tvInfo, err := tmdb1.GetTvAiringToday(nil)
	if err != nil {
		log.Fatal(err)
	}

	lines, err := readLines("showlist.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	var allshows AllShows
	var myshows MyShows

	for _, result := range tvInfo.Results {
		for _, line := range lines {
			if line == result.OriginalName {
				myshows = append(myshows, Shows{
					line,
					"https://image.tmdb.org/t/p/original" + result.PosterPath,
				})
			}
		}
		allshows = append(allshows, Shows{
			result.OriginalName,
			"https://image.tmdb.org/t/p/original" + result.PosterPath,
		})
	}

	data := &Data{allshows, myshows}

	tplerr := tpl.Execute(os.Stdout, data)
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
