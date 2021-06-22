package sgi

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

// ConnectToMongo connect to Mongo in a robust manner
func ConnectToMongo(host string) mongo.Client {
	contextLogger := log.WithFields(log.Fields{"function": "ConnectToMongo"})
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+host+":27017"))
	if err != nil {
		contextLogger.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			contextLogger.Fatal(err)
		}
	}()
	return *client
}

// Socrates object for communicating with the API
type Socrates struct {
	LogLevel int
	Host string
	Protocol string
	VerifySSL bool
	Username string
	Password string
	Headers struct {
		ContentType string
		Authorization string
		Token string
	}
	Client *http.Client
}

// NewClient construct an authenticated Socrates client object
func NewClient(
	logLevel int,
	host string,
	protocol string,
	username string,
	password string) Socrates {
	contextLogger := log.WithFields(log.Fields{"function": "NewClient"})
	s := Socrates{
		LogLevel: logLevel,
		Host: host,
		Protocol: protocol,
		Username: username,
		Password: password,
	}
	s.Headers.ContentType = "application/json"
	s.Client = &http.Client{}
	req, err := http.NewRequest("POST", s.Protocol+"://"+s.Host+"/auth", nil)
	if err != nil {
		contextLogger.Fatal(err)
	}
	req.Header.Add("Content-Type", s.Headers.ContentType)
	resp, err := s.Client.Do(req)
	if err != nil {
		contextLogger.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return s
}

type Response struct {
	status bool
	response string
}

func GetDefinition(api string, module string, name string) Response {
	return Response{}
}