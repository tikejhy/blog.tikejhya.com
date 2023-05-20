// Simple webapp renders html template from yaml config
// Usages gin-gonic framework
// yaml.v3
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Gist struct {
	Title string   `yaml:"title"`
	Link  string   `yaml:"link"`
	Tags  []string `yaml:"tags"`
}

type SumRecord struct {
	Tag string
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	// Load YAML data from file
	yamlFile, err := ioutil.ReadFile("gist.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %s", err.Error())
	}

	// Unmarshal YAML data into a map with a slice of Gist structs
	var data_two map[string][]Gist
	err = yaml.Unmarshal(yamlFile, &data_two)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML data: %s", err.Error())
	}

	r.GET("/", func(c *gin.Context) {
		//all_tags to []map[string]string (a slice of maps with string keys and string values)
		all_tags := make([]map[string]string, 0)
		for _, gist := range data_two["ghGists"] {
			for _, tags := range gist.Tags {
				//tagsData to a map[string]string
				tagsData := map[string]string{
					"tags": tags,
				}
				all_tags = append(all_tags, tagsData)
			}
		}
		srMap := make(map[SumRecord]int)
		// loop over records
		for _, value := range all_tags {
			sr := SumRecord{
				Tag: value["tags"],
			}
			//creates new counter or increments existing pair counter by 1
			srMap[sr] += 1
		}

		tag_weight := make(map[string]int)
		for sr, weight := range srMap {
			tag_weight[sr.Tag] = weight
		}

		// create a slice of maps to store the data for each iteration of the loop.
		// make([]map[string]string, 0) creates an empty slice of maps with string keys and string values.
		// the second argument 0 is specifying the initial length of the slice as zero.
		gists := make([]map[string]string, 0)
		for _, gist := range data_two["ghGists"] {
			gistData := map[string]string{
				"title": gist.Title,
				"link":  gist.Link,
			}
			gists = append(gists, gistData)
		}

		// Then pass this slice to the c.HTML method to render the template.
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"gists":      gists,
			"tag_weight": tag_weight,
		})
	})

	r.GET("/tags/:name", func(c *gin.Context) {
		name := c.Param("name")
		data := loadYAML("gist.yaml", name)
		gistdata := make([]map[string]string, 0)
		for _, gist := range data {
			link := map[string]string{
				"link": gist.Link,
			}

			for _, value := range gist.Tags {
				if value == name {
					gistdata = append(gistdata, link)

				}
			}
		}

		c.HTML(http.StatusOK, "tags.tmpl", gin.H{
			"gists": gistdata,
			"name":  name,
		})
	})

	//Start the server
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func loadYAML(filename string, tag string) []Gist {
	// Read YAML file
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	// Unmarshal YAML data into a slice of Gist structs
	var data map[string][]Gist
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil
	}

	// fmt.Println(data["ghGists"])
	gists := make([]map[string]string, 0)
	for _, gist := range data["ghGists"] {
		for _, tags := range gist.Tags {
			//tagsData to a map[string]string
			dataTags := map[string]string{
				"gists": "gists",
			}
			if strings.Contains(tags, tag) {
				gists = append(gists, dataTags)

			}
		}
	}
	return data["ghGists"]
}
