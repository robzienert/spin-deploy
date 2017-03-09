package spinnaker

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const defaultTimeoutDuration = "5s"

type Config struct {
	Host         string
	Timeout      time.Duration
	X509CertPath string
	X509KeyPath  string
	Verbose      bool
}

func NewConfig() (*Config, error) {
	if viper.GetString("client.host") == "" {
		return nil, errors.New("missing spinnaker.client.host config value")
	}

	timeoutVal := viper.GetString("client.timeout")
	if timeoutVal == "" {
		timeoutVal = defaultTimeoutDuration
	}
	timeout, err := time.ParseDuration(timeoutVal)
	if err != nil {
		return nil, errors.Wrap(err, "parsing timeout duration")
	}

	return &Config{
		Host:         viper.GetString("client.host"),
		Timeout:      timeout,
		X509CertPath: viper.GetString("client.x509CertPath"),
		X509KeyPath:  viper.GetString("client.x509KeyPath"),
	}, nil
}

func (c Config) TLSEnabled() bool {
	return c.X509CertPath != "" && c.X509KeyPath != ""
}

type Client interface {
	StartPipeline(req StartPipelineRequest) (string, error)
}

func New(conf Config) (Client, error) {
	c := &client{config: conf}

	httpClient, err := initHTTPClient(conf)
	if err != nil {
		return nil, errors.Wrap(err, "creating HTTP client")
	}
	c.httpClient = httpClient

	return c, nil
}

type client struct {
	config     Config
	httpClient *http.Client
}

func (c *client) StartPipeline(request StartPipelineRequest) (string, error) {
	p := StartPipeline{
		Type: "templatedPipeline",
		Config: Pipeline{
			Pipeline: PipelineConfig{
				Application:   request.Application,
				Name:          request.Name,
				Template:      request.Template,
				Configuration: request.Configuration,
			},
		},
	}
	j, err := json.Marshal(p)
	if err != nil {
		return "", errors.Wrap(err, "marshaling pipeline json")
	}

	endpoint := c.config.Host + "/pipelines/start"
	fmt.Println(endpoint)
	fmt.Println(string(j))

	req, err := http.NewRequest("POST", c.config.Host+"/pipelines/start", bytes.NewBuffer(j))
	if err != nil {
		return "", errors.Wrap(err, "creating http request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "performing http request")
	}
	defer resp.Body.Close()

	if strings.HasPrefix(resp.Status, "5") {
		if b, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Println(string(b))
		}
		return "", errors.Wrap(err, "received 5xx error")
	}
	fmt.Println(resp.Status)

	ref := orcaRefResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&ref); err != nil {
		if b, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Printf(string(b))
		}
		return "", errors.Wrap(err, "unmarshaling http response")
	}

	return c.config.Host + ref.Ref, nil
}

func initHTTPClient(conf Config) (*http.Client, error) {
	if !conf.TLSEnabled() {
		log.Println("TLS not enabled!")
		return &http.Client{Timeout: conf.Timeout}, nil
	}

	cert, err := tls.LoadX509KeyPair(conf.X509CertPath, conf.X509KeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "loading x509 keypair")
	}

	clientCACert, err := ioutil.ReadFile(conf.X509CertPath)
	if err != nil {
		return nil, errors.Wrap(err, "loading client ca cert")
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()

	return &http.Client{
		Timeout: conf.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}
