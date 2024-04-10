package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/siddhant-vij/HTML-Link-Parser/parser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type UrlEntry struct {
	Loc      string `xml:"loc"`
	Priority string `xml:"priority"`
}

type UrlSet struct {
	Xmlns string     `xml:"xmlns,attr"`
	Urls  []UrlEntry `xml:"url"`
}

func main() {
	baseUrl := flag.String("url", "https://www.calhoun.io/", "The base URL of the sitemap")

	flag.Parse()

	internalLinks := make(map[string]int)
	traverseRecursively(*baseUrl, *baseUrl, &internalLinks)

	priorityMap := assignPriorityToLinks(internalLinks)
	sortedKeys := sortMapKeys(priorityMap)

	result := buildXmlSiteMap(&sortedKeys, &priorityMap)
	fmt.Println("Result: \n" + xml.Header + result)
}

func buildXmlSiteMap(sortedKeys *[]string, priorityMap *map[string]float64) string {
	var urlSet UrlSet
	urlSet.Xmlns = xmlns
	for _, url := range *sortedKeys {
		urlSet.Urls = append(urlSet.Urls, UrlEntry{
			Loc:      url,
			Priority: fmt.Sprintf("%.5f", (*priorityMap)[url]),
		})
	}
	xmlBytes, _ := xml.MarshalIndent(urlSet, "", "  ")
	return string(xmlBytes)
}

func sortMapKeys(m map[string]float64) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return keys
}

func assignPriorityToLinks(internalLinks map[string]int) map[string]float64 {
	if len(internalLinks) == 0 {
		return nil
	}

	var minCount, maxCount int = -1, -1

	for _, count := range internalLinks {
		if minCount == -1 || count < minCount {
			minCount = count
		}
		if maxCount == -1 || count > maxCount {
			maxCount = count
		}
	}

	if minCount == maxCount {
		priorityMap := make(map[string]float64)
		for url := range internalLinks {
			priorityMap[url] = 0.5
		}
		return priorityMap
	}

	priorityMap := make(map[string]float64)
	for url, count := range internalLinks {
		normalized := (float64(count) - float64(minCount)) / (float64(maxCount) - float64(minCount))
		priorityMap[url] = normalized
	}

	return priorityMap
}

func traverseRecursively(url, baseUrl string, internalLinks *map[string]int) error {
	links, err := filterURLsForDomain(url, baseUrl)
	if err != nil {
		return err
	}

	for _, link := range links {
		if _, ok := (*internalLinks)[link]; !ok {
			(*internalLinks)[link] = 1
			traverseRecursively(link, baseUrl, internalLinks)
		} else {
			(*internalLinks)[link]++
		}
	}

	return nil
}

func filterURLsForDomain(urlAtHand, baseUrl string) ([]string, error) {
	htmlStr, err := getHTMLFromURL(urlAtHand)
	if err != nil {
		return nil, err
	}

	links, err := parser.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var filteredLinks []string

	for _, link := range links {
		linkUrl, err := url.Parse(link.Href)
		if err != nil {
			return nil, err
		}

		if linkUrl.Scheme != "http" && linkUrl.Scheme != "https" {
			if link.Href[0] == '/' {
				if baseUrl[len(baseUrl)-1] == '/' {
					baseUrl = baseUrl[:len(baseUrl)-1]
				}
				filteredLinks = append(filteredLinks, baseUrl+link.Href)
			} else {
				filteredLinks = append(filteredLinks, link.Href)
			}
		} else {
			baseUrlParsed, err := url.Parse(baseUrl)
			if err != nil {
				return nil, err
			}
			if baseUrlParsed.Host == linkUrl.Host {
				filteredLinks = append(filteredLinks, link.Href)
			}
			// Think about subdomains - if you want them to be included in the filtered links.
			// Currently subdomains (not matching the baseURL subdomain) are ignored as hostnames (from net/url includes subdomains) are checked for equailty.
		}
	}

	return removeDuplicates(filteredLinks), nil
}

func getHTMLFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func removeDuplicates(links []string) []string {
	uniqueLinks := make(map[string]bool)
	for _, link := range links {
		uniqueLinks[link] = true
	}

	var uniqueLinksSlice []string

	for link := range uniqueLinks {
		uniqueLinksSlice = append(uniqueLinksSlice, link)
	}

	return uniqueLinksSlice
}
