package della

import (
	"io"
	"io/ioutil"
	"regexp"
)

var listRegexp = regexp.MustCompile(`request_code="([\d]*)`)
var dateupsRegexp = regexp.MustCompile(`dateups="([\d]*)`)

func parseData(reader io.ReadCloser) (*CargosData, error) {
	var cd = &CargosData{}

	body, err := readBody(reader)
	if err != nil {
		return cd, err
	}
	bodyString := string(body)

	cd.Ids, err = parseCargoList(bodyString)
	if err != nil {
		return cd, err
	}

	cd.Dateups = parseDateups(bodyString)

	return cd, nil
}

func readBody(reader io.ReadCloser) ([]byte, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return bytes, nil
}

func parseDateups(body string) string {
	var dateups string
	for _, match := range dateupsRegexp.FindAllStringSubmatch(body, 10) {
		dateups = match[1]
	}
	return dateups
}

func parseCargoList(body string) ([]string, error) {
	ids := []string{}
	for _, match := range listRegexp.FindAllStringSubmatch(body, 30) {
		ids = append(ids, match[1])
	}
	return ids, nil
}
