package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:   "serve <port>",
	Short: "serve start a remote server",
	Long:  `serve start a remote server`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
		checkError(err)

		http.HandleFunc("/down", downHandlerWrapper(vaultConf))
		http.HandleFunc("/up", upHandlerWrapper(vaultConf))

		fmt.Println("server start at ", args[0], "...")
		err = http.ListenAndServe(":"+args[0], nil)
		checkError(err)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// 接收客户端的下载数据请求（拉取）
func downHandlerWrapper(vaultConf *config.VaultConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var data map[string]string
		err = json.Unmarshal(content, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = checkRemoteCredential(vaultConf.Credentials, data["username"], data["token"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		records, err := persist.DumpRecord(vaultConf.PersistConf.DataFile, data["username"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		dataJson, err := json.Marshal(records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(dataJson)
	}
}

// 接收客户端的上传数据请求（推送）
func upHandlerWrapper(vaultConf *config.VaultConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var data map[string]string
		err = json.Unmarshal(content, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = checkRemoteCredential(vaultConf.Credentials, data["username"], data["token"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		recordsJsonData, err := base64.StdEncoding.DecodeString(data["records"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var records []*persist.AccountRecord
		err = json.Unmarshal(recordsJsonData, &records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = persist.ImportRecord(vaultConf.PersistConf.DataFile, records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
}
