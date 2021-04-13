package utils

import (
	"regexp"
	"errors"
)

const (
	uriWithoutName = 1
	dbname = 2
)

// group 1 is URI without db name
// group 2 is db name
const db_pattern = `(?m)(mongodb:\/\/[\w]+:[\w]+@[\w\d.]+:[\d]+)\/(\w+)`

// GetDbNameByReg
// return a name of a database 
// if not find name return "" and log warning
func GetDbNameByReg(URI string) (string, error) {
	URIWithoutName := getGroupFromURI(URI, dbname) 
	if URIWithoutName == "" {
		return "", errors.New("Unable get db name")
	}
	return URIWithoutName, nil
}

func GetDBURIWithoutName(URI string) (string, error) {
	URIWithoutName := getGroupFromURI(URI, uriWithoutName) 
	if URIWithoutName == "" {
		return "", errors.New("Unable get uri without name")
	}

	return URIWithoutName, nil
}

func getGroupFromURI(URI string, group int) string {
	re := regexp.MustCompile(db_pattern)
	all := re.FindStringSubmatch(URI)
	if len(all) != 3 {
		return ""
	}

	return all[group]
}