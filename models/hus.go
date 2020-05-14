package models

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "github.com/prometheus/common/log"
    "io"
    "io/ioutil"
    "math/rand"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type Hus struct {
}

type Condition struct {
    Field string
    Operator string
    Values []interface{}
}

type Paging struct {
    Index int
    Size int
}

type Paginator struct {
    Url string
    Count int
    CurrentPage int
    PageSize int
    Padding int
}

var (
    signingKey             = []byte("Husol!@#123ok")
    keyFunc    jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return signingKey, nil }
)

func (hus *Hus) Log(data interface{}, overwrite bool) error {
    jsonData := time.Now().Format("02-Jan-2006 15:04:05") + "  "
    json, err := json.MarshalIndent(data, "", "\t")
    jsonData += string(json)
    if overwrite {
        err = ioutil.WriteFile("/tmp/debug", []byte(jsonData), 0664)
    } else {
        f, err := os.OpenFile("/tmp/debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
        if err != nil {
            return err
        }
        defer f.Close()
        _, err = f.Write([]byte(jsonData))
        _, err = f.WriteString("\n\n")
    }

    if err != nil {
        return err
    }
    return nil
}

func (hus *Hus) ErrorLog(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func (hus *Hus) RandomString(n int) string {
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func (hus *Hus) DecodeStrToObj(encode_str string, obj interface{}) {
    json.Unmarshal([]byte(encode_str), obj)
}

func (hus *Hus) DeleteDirFile(path string) bool {
    if path == "" {
        return false
    }

    if err := os.RemoveAll(path); err != nil {
        return false
    }

    return true
}

func (hus *Hus) DeleteEmptyDir(pathDir string, isRecursive bool) bool {
    if pathDir == "" {
        return false
    }

    isEmptyDir, _ := hus.IsEmptyDir(pathDir)
    if isEmptyDir {
        if err := os.RemoveAll(pathDir); err != nil {
            return false
        }
        if isRecursive {
            hus.DeleteEmptyDir(filepath.Dir(pathDir), isRecursive)
        }
    }

    return true
}

func (hus *Hus) IsEmptyDir(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    // Read in ONLY one file
    _, err = f.Readdir(1)

    // And if the file is EOF, the directory is empty.
    if err == io.EOF {
        return true, nil
    }
    return false, err
}

func (hus *Hus) GetFilePathSuffix(file string, suffix string) string {
    filePath := filepath.Dir(file)
    fileName := filepath.Base(file)
    extension := filepath.Ext(file)
    return filepath.Join(filePath, strings.TrimSuffix(fileName, extension) + suffix + extension)
}

func (hus *Hus) ValidateFiles(files []*multipart.FileHeader, types []string) bool {
    for _, file := range files {
        check := false
        for _, fileType := range types {
            if file.Header.Get("Content-Type") == fileType {
                check = true
            }
        }
        if !check {
            return false
        }
    }

    return true
}

func (hus *Hus) CallAPI(host, uri string, params map[string]interface{}) []byte {
    method := params["method"].(string)
    jsonData := params["data"]
    jsonValue, _ := json.Marshal(jsonData)
    request, _ := http.NewRequest(method, host +"/"+ uri, bytes.NewBuffer(jsonValue))
    request.Header.Add("Content-Type", "application/json")

    if token, ok := params["token"]; ok {
        request.Header.Add("Authorization", token.(string))
    }
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    httpClient := &http.Client{Transport: tr}
    resp, _ := httpClient.Do(request)
    response, _ := ioutil.ReadAll(resp.Body)

    return response
}

func (hus *Hus) CallFormUploadAPI(host, uri string, authToken string, params map[string]string) []byte {
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    if len(params["fileField"]) > 0 {
        file, err := os.Open(params["filePath"])
        if err != nil {
            log.Fatal(err)
        }
        fileContents, err := ioutil.ReadAll(file)
        if err != nil {
            log.Fatal(err)
        }
        fi, err := file.Stat()
        if err != nil {
            log.Fatal(err)
        }
        file.Close()

        part, err := writer.CreateFormFile(params["fileField"], fi.Name())
        if err != nil {
            log.Fatal(err)
        }
        part.Write(fileContents)
    }

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }
    err := writer.Close()
    if err != nil {
        log.Fatal(err)
    }

    request, _ := http.NewRequest("POST", host +"/"+ uri, body)
    request.Header.Add("Content-Type", writer.FormDataContentType())

    if len(authToken) > 0 {
        request.Header.Add("Authorization", authToken)
    }
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    httpClient := &http.Client{Transport: tr}
    resp, _ := httpClient.Do(request)
    response, _ := ioutil.ReadAll(resp.Body)

    if (len(params["filePath"]) > 0) {
        //Delete file after uploading
        err = os.Remove(params["filePath"])
        if err != nil {
            log.Fatal(err)
        }
    }

    return response
}
