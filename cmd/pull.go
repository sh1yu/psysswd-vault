package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sh1yu/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull account info from remote server",
	Long:  `pull account info from remote server`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {

		vaultConf, username, _ := runPreCheck(cmd)

		remoteServerAddr, err := cmd.Flags().GetString("remote")
		checkError(err)

		if remoteServerAddr == "" {
			remoteServerAddr = vaultConf.RemoteConf.ServerAddr
		}

		if remoteServerAddr == "" {
			fmt.Println("please provide remote server addr with -r or 'remote.server_addr' in config file.")
			os.Exit(1)
		}

		remoteToken := getRemoteCredential(vaultConf.Credentials, username)
		if remoteToken == "" {
			fmt.Println("please config 'credentials' for user " + username + " in config file before pull.")
			os.Exit(1)
		}

		err = runPull(vaultConf.PersistConf.DataFile, username, remoteToken, remoteServerAddr)
		checkError(err)
	},
}

func init() {
	pullCmd.Flags().StringP("remote", "r", "", "given remote server addr.")
	rootCmd.AddCommand(pullCmd)
}

func runPull(dataFile, username, remoteToken, remoteServerAddr string) error {

	fmt.Printf("Pulling for username %v from remote %v ...\n", username, remoteServerAddr)

	data := map[string]string{"username": username, "token": remoteToken}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(remoteServerAddr+"/down", "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return err
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(respContent))
	}

	var records []*persist.AccountRecord
	err = json.Unmarshal(respContent, &records)
	if err != nil {
		return err
	}

	return persist.ImportRecord(dataFile, records)
}
