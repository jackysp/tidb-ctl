// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// root command flags
var (
	host   net.IP
	port   uint16
	genDoc bool
)

const (
	rootUse   = "tidb-ctl"
	rootShort = "TiDB Controller"
	rootLong  = "TiDB Controller (tidb-ctl) is a command line tool for TiDB Server (tidb-server)."
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   rootUse,
	Short: rootShort,
	Long:  rootLong,
	RunE: func(_ *cobra.Command, args []string) error {
		if !genDoc {
			return nil
		}
		docDir := "./doc"
		docCmd := &cobra.Command{
			Use:   rootUse,
			Short: rootShort,
			Long:  rootLong,
		}
		docCmd.AddCommand(mvccRootCmd, schemaRootCmd, regionRootCmd)
		return doc.GenMarkdownTree(docCmd, docDir)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func httpPrint(path string) error {
	resp, err := http.Get("http://" + host.String() + ":" + strconv.Itoa(int(port)) + "/" + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func init() {
	hostFlagName := "host"
	portFlagName := "port"

	rootCmd.AddCommand(mvccRootCmd)
	rootCmd.AddCommand(regionRootCmd)
	rootCmd.AddCommand(schemaRootCmd)

	rootCmd.PersistentFlags().IPVarP(&host, hostFlagName, "H", net.IP("127.0.0.1"), "TiDB server host")
	rootCmd.PersistentFlags().Uint16VarP(&port, portFlagName, "P", 10080, "TiDB server port")
	rootCmd.MarkPersistentFlagRequired(hostFlagName)
	rootCmd.MarkPersistentFlagRequired(portFlagName)
	rootCmd.Flags().BoolVar(&genDoc, "doc", false, "generate doc file")
	rootCmd.Flags().MarkHidden("doc")
}
