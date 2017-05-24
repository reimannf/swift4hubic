package main

import (
	"fmt"
	"os"

	"github.com/reimannf/swift4hubic/swift4hubic"
)

func main() {
	if len(os.Args) < 2 {
		s := fmt.Sprintf("usage: %s <config-file>", os.Args[0])
		fmt.Fprintln(os.Stderr, s)
		os.Exit(1)
	}

	config, err := swift4hubic.GetConfiguration(os.Args[1])
	if err != nil {
		swift4hubic.Log(swift4hubic.LogFatal, err.Error())
	}

	swift4hubic.NewServer(config)
}

// TODO Context Set User/Pwd
// TODO Context Random String
// TODO auth.go - Store Token https://gist.github.com/jfcote87/89eca3032cd5f9705ba3 https://github.com/golang/oauth2/issues/84
// TODO auth.go - Read Token http://stackoverflow.com/questions/28685033/how-to-handle-refresh-tokens-in-golang-oauth2-client-lib
// TODO auth.go - Multi Account
// TODO SSL
// TODO command line client
// TODO Set Timeout httpclient
// TODO Check hc := &http.Client{Timeout: 2 * time.Second} ctx := context.WithValue(context.Background(), oauth2.HTTPClient, hc)
