package server

/*import (
	"ITLab-Projects/config"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)
type InstallToken struct {
	Token 	string 	`json:"token"`
}

func generateTokenForGithub() string{
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "generateTokenForGithub",
			"error"	:	err,
		},
		).Fatal("Error in RSA private key parsing!")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * time.Duration(7)).Unix(),
		"iss": cfg.Auth.Github.AppID,
	})
	tokenString, err := token.SignedString(privateKey)
	return tokenString
}

func getInstallationTokenFor(login string) string {
	installations := make([]config.Installation, 0)
	var installToken InstallToken

	tokenString := generateTokenForGithub()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/app/installations", nil)
	req.Header.Set("Accept", "application/vnd.github.machine-man-preview+json")
	req.Header.Set("Authorization", "Bearer " + tokenString)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"function" : "http.Get",
			"handler" : "getInstallationTokenFor",
			"url"	: "https://api.github.com/app/installations",
			"error"	:	err,
		},
		).Warn("Can't reach API!")
	}


	var id int64
	json.NewDecoder(resp.Body).Decode(&installations)
	for _, installation := range installations {
		if installation.Account.Login == login {
			id = installation.ID
			break
		}
	}
	resp.Body.Close()
	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", id)
	req, err = http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/vnd.github.machine-man-preview+json")
	req.Header.Set("Authorization", "Bearer " + tokenString)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)

	json.NewDecoder(resp.Body).Decode(&installToken)
	resp.Body.Close()
	return installToken.Token
}*/