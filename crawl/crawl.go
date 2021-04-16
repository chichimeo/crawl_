package crawl

import (
	"Crawl/malware"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var (
	regexLineBreak = regexp.MustCompile(`.*\n`)
	regexString    = regexp.MustCompile(`\S+`)
	regexDay       = regexp.MustCompile(`\d{4}-\d{2}-\d{2}/`)
)

const (
	urlGetData = "https://malshare.com/daily/%s"
	urlDay     = "%smalshare_fileList.%s.all.txt"
)

type MalwareHandler struct {
	Repository malware.Repository
}

func GetAllDays(s string) []string {
	item := regexDay.FindAllString(s, -1)
	keys := make(map[string]bool)
	var list []string
	for _, entry := range item {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetData(s string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(urlGetData, s))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}
func (re *MalwareHandler) GetMalware(wg *sync.WaitGroup, allDays <-chan string) {
	defer wg.Done()
	for day := range allDays {
		url := fmt.Sprintf(urlDay, day, day[:len(day)-1])
		body, err := GetData(url)
		if err != nil {
			fmt.Println("error", err)
			continue
		}

		elements := regexLineBreak.FindAllString(body, -1)

		for _, el := range elements {
			item := regexString.FindAllString(el, -1)

			it := malware.Malware{}
			if len(item) < 3 {
				continue
			}
			if item[0] != "NULL" {
				it.Md5 = item[0]
			}
			if item[1] != "NULL" {
				it.Sha1 = item[1]
			}
			if item[2] != "NULL" {
				it.Sha256 = item[2]
			}
			it.Date = day[:len(day)-1]
			err = re.Repository.Insert(it)
			if err != nil {
				fmt.Println("error", err)
			}
		}
	}
}

func (re *MalwareHandler) Crawl() {
	body, err := GetData("")
	if err != nil {
		return
	}

	allDays := GetAllDays(body)
	jobs := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go re.GetMalware(&wg, jobs)
	}
	for _, element := range allDays {
		jobs <- element
	}
	close(jobs)
	wg.Wait()
}
