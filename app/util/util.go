package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReloadService(url string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Check(err)
	defer resp.Body.Close()
}

func CreateConfigurationFile(configFile string, templateFile string, jsonInt interface{}) {
	fTmpl, err := ioutil.ReadFile(templateFile)
	Check(err)
	fConf, err := os.Create(configFile)
	Check(err)

	t := template.New("t")
	t, err = t.Parse(string(fTmpl))
	Check(err)

	err = t.Execute(fConf, jsonInt)
	Check(err)

	fConf.Close()
	Check(err)
}
