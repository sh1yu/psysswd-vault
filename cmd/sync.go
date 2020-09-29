package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync <remote-addr>",
	Short: "sync account info from remote addr",
	Long:  `sync account info from remote addr`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, password := runPreCheck(cmd)
		err := runSync(vaultConf.PersistConf.DataFile, username, password, args[0])
		checkError(err)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(dataFile, username, password, remoteAddr string) error {

	data := map[string]string{"username": username, "password": password}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(remoteAddr+"/sync", "application/json", bytes.NewReader(dataJson))
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
