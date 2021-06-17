package utils

import (
	"strings"
	"errors"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
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

// If in database name not found word test panic
func ValidateTestURI(URI, TestURI string) {
	uri, err := connstring.Parse(URI)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/utils",
				"func": "ValidateTestURI",
				"err": err,
			},
		).Panic("Faield to parse mongodb uri")
	}

	testuri, err := connstring.Parse(TestURI)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/utils",
				"func": "ValidateTestURI",
				"err": err,
			},
		).Panic("Faield to parse mongodb test uri")
	}

	if uri.Database == testuri.Database {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/utils",
				"func": "ValidateTestURI",
				"err": fmt.Errorf("URI and TestURI should be different because in test - drops test database to validate all data"),
			},
		).Panic()
	}

	if !strings.Contains(
		strings.ToLower(testuri.Database),
		"test",
	) {
		logrus.WithFields(
			logrus.Fields{
				"package": "repositories/utils",
				"func": "ValidateTestURI",
				"err": fmt.Errorf(`test uri database should contains word "test" example:"mongodb://user:pass@net:port/test_database"`),
			},
		).Panic()
	}
}
