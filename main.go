package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// "github.com/!puerkito!bio/goquery"

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	// fmt.Println(totalPages)

	// c := make(chan []extractedJob)
	
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	// fmt.Println(len(jobs))
	writeJobs(jobs)
	fmt.Println("done extracting", len(jobs), "jobs")
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "TITLE", "LOCATION", "SALARY", "SUMMARY"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range(jobs) {
		jobSlice := []string{"kr.indeed.com/viewjob?jk="+job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ") 
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("requesting", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	// fmt.Println(id, exists, title, location)
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())

	return extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
		summary: summary,
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	// if err != nil {
	// 	log.Fatalln(err)
	// } else if res.StatusCode != 200 {
	// 	log.Fatalln("Request failed with Status:", res.StatusCode)
	// }
	checkErr(err)
	checkStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// fmt.Println(doc)
	doc.Find(".pagination").Each(func(i int, sel *goquery.Selection) {
		pages = sel.Find("a").Length()
	})
	
	return pages
}