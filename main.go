package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	flag "github.com/spf13/pflag"
)

type Environment struct {
	Key   string
	Value string
}

var (
	outputFlag   = flag.StringP("output", "o", "stdout", "output to 'stdout', 'http', 'fluent', 'all'")
	portFlag     = flag.StringP("port", "p", "80", "port number for http server")
	fluentHost   = flag.StringP("fluent-host", "h", "localhost", "fluent host")
	fluentPort   = flag.IntP("fluent-port", "t", 24224, "fluent port")
	fluentTag    = flag.StringP("fluent-tag", "g", "env", "fluent tag")
	Environments []Environment
)

func genHtml() []byte {
	var envs []Environment
	for _, env := range os.Environ() {
		kv := strings.Split(env, "=")
		k, v := kv[0], kv[1]
		envs = append(Environments, Environment{Key: k, Value: v})
	}
	html := []byte(`
		<html>
		<head>
		<title>Environment Variables</title>
		</head>
		<body>
		<h1>Environment Variables</h1>
		<table>
		<tr>
		<th>Key</th>
		<th>Value</th>
		</tr>
		`)
	for _, env := range envs {
		html = append(html, []byte("<tr><td>"+env.Key+"</td><td>"+env.Value+"</td></tr>")...)
	}
	html = append(html, []byte(`</table></body></html>`)...)
	return html
}

func main() {
	flag.Parse()
	for _, env := range os.Environ() {
		kv := strings.Split(env, "=")
		k, v := kv[0], kv[1]
		Environments = append(Environments, Environment{Key: k, Value: v})
	}
	wg := new(sync.WaitGroup)
	if *outputFlag == "stdout" || *outputFlag == "all" {
		log.Println("start logging environment variables in stdout")
		wg.Add(1)
		go func() {
			log.Printf("start logging environment variables: %d", len(Environments))
			for {
				for _, env := range Environments {
					log.Printf("%s=%s", env.Key, env.Value)
				}
				time.Sleep(1 * time.Minute)
			}
			//wg.Done()
		}()
	}
	if *outputFlag == "fluent" || *outputFlag == "all" {
		log.Println("start logging environment variables in fluent")
		wg.Add(1)
		go func() {
			log.Printf("start logging environment variables: %d", len(Environments))
			logger, err := fluent.New(fluent.Config{
				FluentHost: *fluentHost,
				FluentPort: *fluentPort,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer logger.Close()
			for {
				for _, env := range Environments {
					logger.PostWithTime(*fluentTag, time.Now(), env)
				}
				time.Sleep(1 * time.Minute)
			}
			//wg.Done()
		}()
	}
	if *outputFlag == "http" || *outputFlag == "all" {
		log.Println("start logging environment variables in http")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write(genHtml())
		})
		if err := http.ListenAndServe(":"+*portFlag, nil); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}
