package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/crypto/pkcs12"
)

const errorStr = "http request error, url=%v , statusCode=%v"

const (
	contentTypeJson = "application/json;charset=utf-8"
	contentTypeXml  = "application/xml;charset=utf-8"
)

// HttpGet get request
func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorStr, url, response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}

// HttpPost post request
func HttpPost(url string, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))
	response, err := http.Post(url, "", body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorStr, url, response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}

// PostJson post json request
func PostJson(url string, obj interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)
	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(url, contentTypeJson, body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorStr, url, response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}

// PostJsonWithRespContentType post json request return:Content-Type + Body
func PostJsonWithRespContentType(url string, obj interface{}) ([]byte, string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(url, contentTypeJson, body)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf(errorStr, url, response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	contentType := response.Header.Get("Content-Type")
	return responseData, contentType, err
}

// PostFile 上传文件
func PostFile(fieldname, filename, url string) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			FieldName: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(fields, url)
}

// MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	FieldName string
	Value     []byte
	Filename  string
}

// PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, url string) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.FieldName, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.FieldName)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := http.Post(url, contentType, bodyBuf)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}

// PostXML perform a HTTP/POST request with XML body
func PostXML(url string, obj interface{}) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	response, err := http.Post(url, contentTypeXml, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorStr, url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// httpWithTLS CA证书
func httpWithTLS(rootCa, key string) (*http.Client, error) {
	var client *http.Client
	certData, err := ioutil.ReadFile(rootCa)
	if err != nil {
		return nil, fmt.Errorf("unable to find cert path=%s, error=%v", rootCa, err)
	}
	cert := pkcs12ToPem(certData, key)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tr := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr}
	return client, nil
}

// pkcs12ToPem 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) tls.Certificate {
	blocks, err := pkcs12.ToPEM(p12, password)
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()
	if err != nil {
		panic(err)
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}

// PostXMLWithTLS perform a HTTP/POST request with XML body and TLS
func PostXMLWithTLS(url string, obj interface{}, ca, key string) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	client, err := httpWithTLS(ca, key)
	if err != nil {
		return nil, err
	}
	response, err := client.Post(url, contentTypeXml, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errorStr, url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
