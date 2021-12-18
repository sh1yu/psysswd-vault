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

		http.HandleFunc("/down", downHandlerWrapper(vaultConf.PersistConf.DataFile))
		http.HandleFunc("/up", upHandlerWrapper(vaultConf.PersistConf.DataFile))

		fmt.Println("server start at ", args[0], "...")
		err = http.ListenAndServe(":"+args[0], nil)
		checkError(err)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// 接收客户端的下载数据请求（拉取）
func downHandlerWrapper(dataFile string) func(w http.ResponseWriter, r *http.Request) {
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

		exist, valid, err := persist.CheckUser(dataFile, data["username"], data["password"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if !exist {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("user not registered: " + data["username"]))
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Permission Denied."))
			return
		}

		records, err := persist.DumpRecord(dataFile, data["username"])
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
func upHandlerWrapper(dataFile string) func(w http.ResponseWriter, r *http.Request) {
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

		exist, valid, err := persist.CheckUser(dataFile, data["username"], data["password"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if !exist {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("user not registered: " + data["username"]))
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Permission Denied."))
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

		err = persist.ImportRecord(dataFile, records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
}
