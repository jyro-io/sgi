package sgi

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	req, err := http.NewRequest(
		"POST",
		s.Protocol+"://"+s.Host+"/auth",
		nil)
	if err != nil {
		contextLogger.Fatal(err)
	}
	req.Header.Add("Content-Type", s.Headers.ContentType)
	resp, err := s.Client.Do(req)
	if err != nil {
		contextLogger.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		contextLogger.Fatal(err)
	}
	return s
}

type Response struct {
	status bool
	response string
}

type Datasource struct {
	Type                 string   `json:"type"`
	Host                 string   `json:"host"`
	Port                 string   `json:"port"`
	TimeoutMs            int      `json:"timeout_ms"`
	TimestampField       string   `json:"timestamp_field"`
	TimestampFormat      string   `json:"timestamp_format"`
	TimestampFormatJulia string   `json:"timestamp_format_julia"`
	TopicPrefix          string   `json:"topic_prefix"`
	Topics               []string `json:"topics"`
	Metadata             struct {
		Window struct {
			Limit struct {
				Upper int `json:"upper"`
				Lower int `json:"lower"`
			} `json:"limit"`
			ScaleFactor float64 `json:"scale_factor"`
		} `json:"window"`
		Join interface{} `json:"join"`
		Etl []struct {
			Operation  string `json:"operation"`
			Name       string `json:"name"`
			Parameters interface{} `json:"parameters,omitempty"`
			PullFields bool `json:"pull_fields,omitempty"`
		} `json:"etl"`
	} `json:"metadata"`
	Replication struct {
		Type     string `json:"type"`
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"replication"`
}

func GetDefinition(s Socrates, api string, module string, name string) (Response, error) {
	contextLogger := log.WithFields(log.Fields{"function": "GetDefinition"})
	var jsonData = []byte(`{
		"operation": "get",
		"name": "`+name+`"
	}`)
	req, err := http.NewRequest(
		"POST",
		s.Protocol+"://"+s.Host+"/"+api+"/"+module,
		bytes.NewBuffer(jsonData))
	if err != nil {
		contextLogger.Fatal(err)
	}
	req.Header.Add("Content-Type", s.Headers.ContentType)
	resp, err := s.Client.Do(req)
	if err != nil {
		contextLogger.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		contextLogger.Fatal(err)
	}
	definition := &Datasource{}
	switch api {
	case "archimedes":
		switch module {
		case "datasource":
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(definition)
		}
	}
	if err != nil {
		contextLogger.Fatal(err)
	}
	return Response{}, nil
}