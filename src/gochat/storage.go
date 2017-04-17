package gochat

import (
	"net/http"
	"encoding/json"
	"os"
	"io/ioutil"
)

type Storage struct {
	storageDirPath string
}

type storageData struct {
	Uuid 		string
	BaseRequest BaseRequest
	PassTicket 	string
	Cookies 	[]*http.Cookie
}

func (this *Storage) setData(Uuid string, baseRequest BaseRequest, passTicket string, cookies []*http.Cookie) error {
	storageStr, _ := json.Marshal(storageData {
		Uuid:			Uuid,
		BaseRequest: 	baseRequest,
		PassTicket: 	passTicket,
		Cookies: 		cookies,
	})

	fileName := this.storageDirPath + "storage_data.json"
	file, err := os.OpenFile(fileName, os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(storageStr)

	return err
}

func (this *Storage) getData() (string, BaseRequest, string, []*http.Cookie, error) {
	fileName := this.storageDirPath + "storage_data.json"
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", BaseRequest{}, "", nil, err
	}
	var storageData storageData
	err = json.Unmarshal(bs, &storageData)
	if err != nil {
		return "", BaseRequest{}, "", nil, err
	}

	return storageData.Uuid, storageData.BaseRequest, storageData.PassTicket, storageData.Cookies, nil
}

func (this *Storage) delData() error {
	fileName := this.storageDirPath + "storage_data.json"
	err := os.Remove(fileName)
	return err
}