package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Author ...
type Author struct {
	ID           int
	Name         string
	Affiliations []string
	PC           int
	CN           int
	HI           int
	PI           float64
	UPI          float64
	Tags         []string
}

// Topic ...
type Topic struct {
	Name  string
	Count int
}

func main() {

	authorsFile := "../AMiner-Author.txt"

	f, err := os.OpenFile(authorsFile, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	/*
		#index 1
		#n O. Willum
		#a Res. Center for Microperipherik, Technische Univ. Berlin, Germany
		#pc 1
		#cn 0
		#hi 0
		#pi 0.0000
		#upi 0.0000
		#t new product;product group;active product;long product lifetime;old product;product generation;new technology;environmental benefit;environmental choice;environmental consequence

		#index 2
		#n D. Wei
		#a Dept. of Electr. & Comput. Eng., Drexel Univ., Philadelphia, PA, USA
		#pc 1
		#cn 0
		#hi 0
		#pi 0.0000
		#upi 0.0000
		#t lowpass filter;multidimensional product filter;orthonormal filterbanks;product filter;new approach;novel approach;challenging problem;iterative quadratic programming;negligible reconstruction error;spectral factorization

	*/

	scanner := bufio.NewScanner(f)
	index := map[int]Author{}

	for scanner.Scan() {
		//
		line := scanner.Text()
		if strings.Index(line, "#index") == 0 {
			indexID := strings.Replace(line, "#index ", "", -1)

			id, _ := strconv.Atoi(indexID)

			author := Author{
				ID: id,
			}

			for scanner.Scan() {
				if isEmpty(scanner.Text()) {
					// add to index
					index[id] = author
					break
				}
				switch {
				case strings.Index(scanner.Text(), "#n") == 0:
					author.Name = stripReturnString(scanner.Text(), "#n")
				case strings.Index(scanner.Text(), "#a") == 0:
					author.Affiliations = stripReturnSlice(scanner.Text(), "#a")
				case strings.Index(scanner.Text(), "#pc") == 0:
					author.PC = stripReturnInt(scanner.Text(), "#pc")
				case strings.Index(scanner.Text(), "#cn") == 0:
					author.CN = stripReturnInt(scanner.Text(), "#cn")
				case strings.Index(scanner.Text(), "#hi") == 0:
					author.HI = stripReturnInt(scanner.Text(), "#hi")
				case strings.Index(scanner.Text(), "#pi") == 0:
					author.PI = stripReturnFloat(scanner.Text(), "#pi")
				case strings.Index(scanner.Text(), "#upi") == 0:
					author.UPI = stripReturnFloat(scanner.Text(), "#upi")
				case strings.Index(scanner.Text(), "#t") == 0:
					author.Tags = stripReturnSlice(scanner.Text(), "#t")

				}
			}

		}
	}

	topicsIndex := map[string][]Author{}

	for _, a := range index {
		for _, t := range a.Tags {
			if _, ok := topicsIndex[t]; !ok {
				topicsIndex[t] = []Author{}
			}
			topicsIndex[t] = append(topicsIndex[t], a)
		}
	}

	topics := []Topic{}
	for k, v := range topicsIndex {
		topics = append(topics, Topic{
			Name:  k,
			Count: len(v),
		})
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Count > topics[j].Count
	})

	os.MkdirAll("./data", 0777)

	// filter the dataset
	for i, t := range topics[0:50] {
		if len(strings.TrimSpace(t.Name)) == 0 {
			continue
		}
		outf, _ := os.Create(fmt.Sprintf("data/nodes_%v.csv", i))
		outf.WriteString(fmt.Sprintf("#topic: %v, #count: %v\n", t.Name, t.Count))
		for _, a := range topicsIndex[t.Name] {
			outf.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", a.ID, a.Name, a.PC, a.CN, a.HI, a.PI, a.UPI))
		}
		outf.WriteString(fmt.Sprintf("\n"))
	}

	fmt.Println(len(index))
	fmt.Println(len(topics))

}

func stripReturnString(s string, t string) string {
	return strings.TrimSpace(strings.Replace(s, t, "", -1))
}

func stripReturnSlice(s string, t string) []string {
	text := strings.TrimSpace(strings.Replace(s, t, "", -1))
	return strings.Split(text, ";")
}

func stripReturnInt(s string, t string) int {
	text := strings.TrimSpace(strings.Replace(s, t, "", -1))
	n, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	return n
}

func stripReturnFloat(s string, t string) float64 {
	text := strings.TrimSpace(strings.Replace(s, t, "", -1))
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
