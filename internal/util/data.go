package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
)

func RangePersistData(conf *config.VaultConfig, f func(key, data []byte)) error {

	indexFilePath, err := config.CreateFileIfNeeded(conf.PersistConf.IndexFile)
	if err != nil {
		return err
	}
	dataFilePath, err := config.CreateFileIfNeeded(conf.PersistConf.DataFile)
	if err != nil {
		return err
	}

	indexData, err := ioutil.ReadFile(indexFilePath)
	if err != nil {
		return err
	}
	bodyData, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return err
	}

	for i := 0; i < len(indexData); i += 32 {
		var keyOffset int64
		var keyLen int32
		binary.Read(bytes.NewBuffer(indexData[i+20:i+28]), binary.LittleEndian, &keyOffset)
		binary.Read(bytes.NewBuffer(indexData[i+28:i+32]), binary.LittleEndian, &keyLen)
		key := bodyData[keyOffset : keyOffset+int64(keyLen)]

		var dataOffset int64
		var dataLen int32
		binary.Read(bytes.NewBuffer(indexData[i+8:i+16]), binary.LittleEndian, &dataOffset)
		binary.Read(bytes.NewBuffer(indexData[i+16:i+20]), binary.LittleEndian, &dataLen)
		enDataAll := bodyData[dataOffset : dataOffset+int64(dataLen)]

		f(key, enDataAll)
	}
	return nil
}

func ModifyData(conf *config.VaultConfig, originKey, data []byte) error {
	//使用master password加盐生成aes-256的key

	indexFilePath, err := config.CreateFileIfNeeded(conf.PersistConf.IndexFile)
	if err != nil {
		return err
	}
	dataFilePath, err := config.CreateFileIfNeeded(conf.PersistConf.DataFile)
	if err != nil {
		return err
	}

	indexData, err := ioutil.ReadFile(indexFilePath)
	if err != nil {
		return err
	}
	bodyData, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return err
	}

	//先将data添加到body之后
	bodyOffset := len(bodyData)
	bodyLen := len(data)
	bodyData = append(bodyData, data...)

	storeKey := pbkdf2.Key(originKey, []byte{}, constant.Pbkdf2Iter, 8, sha256.New)

	for i := 0; i < len(indexData); i += 32 {
		if base64.StdEncoding.EncodeToString(storeKey) == base64.StdEncoding.EncodeToString(indexData[i:i+8]) {
			//已经存在，改密码

			var updateIndexBuf bytes.Buffer
			binary.Write(&updateIndexBuf, binary.LittleEndian, int64(bodyOffset))
			binary.Write(&updateIndexBuf, binary.LittleEndian, int32(bodyLen))
			updateByte := updateIndexBuf.Bytes()

			for j := 0; j < 12; j++ {
				indexData[i+8+j] = updateByte[j]
			}

			ioutil.WriteFile(dataFilePath, bodyData, 0644)
			ioutil.WriteFile(indexFilePath, indexData, 0644)
			return nil
		}
	}

	//不存在，添加user和密码
	keyOffset := len(bodyData)
	keyLen := len(originKey)
	bodyData = append(bodyData, originKey...)

	var addIndexBuf bytes.Buffer
	addIndexBuf.Write(storeKey)
	binary.Write(&addIndexBuf, binary.LittleEndian, int64(bodyOffset))
	binary.Write(&addIndexBuf, binary.LittleEndian, int32(bodyLen))
	binary.Write(&addIndexBuf, binary.LittleEndian, int64(keyOffset))
	binary.Write(&addIndexBuf, binary.LittleEndian, int32(keyLen))
	indexData = append(indexData, addIndexBuf.Bytes()...)

	ioutil.WriteFile(dataFilePath, bodyData, 0644)
	ioutil.WriteFile(indexFilePath, indexData, 0644)

	return nil
}
