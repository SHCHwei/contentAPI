package main

import (
    "log"
    "os"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "time"
    "github.com/go-playground/validator/v10"
)


type ContentData struct {
	ID int8 `json:"id"`
	Name string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required,max=20,min=1"`
	Phone string `json:"phone" validate:"required,len=10"`
	Email string `json:"email" validate:"required,email"`
	Time int64 `json:"time"`
}

var validate *validator.Validate

func main(){

    validate = validator.New()


    http.HandleFunc("/add", add)

	logF("Server Run", nil)
    err := http.ListenAndServe(":8099", nil)
    if err != nil {
        logF("Server ERROR", err)
    }
}

func add(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    w.Header().Set("Content-Type", "application/json")

    logF("Run add", nil)
    
    var rows []ContentData

    //讀取全部
    jsonFile, err := os.Open("data.json")
    if err != nil {
        logF("讀取json檔案失敗 err: ", err)
    }
    defer jsonFile.Close()

    dataList, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(dataList, &rows)

    newRow := ContentData{
        ID: rows[len(rows) - 1].ID + 1,
        Name: r.FormValue("name"),
        Content: r.FormValue("content"),
        Phone: r.FormValue("phone"),
        Email: r.FormValue("email"),
        Time: time.Now().Unix(),
    }


    errs := validate.Struct(newRow)
	if errs != nil {
        log.Println(errs)
        w.Write(buildResponse(map[string]string{"statue": "false", "message": "格式錯誤"}))
        return 
    }

    rows = append(rows, newRow)
	newTeamBytes, err := json.Marshal(rows)
    if err != nil {
        logF("err: ", err)
    }

    err = ioutil.WriteFile("data.json", newTeamBytes, 0644)

    if err != nil {
        w.Write(buildResponse(map[string]string{"statue": "false", "message": "建立失敗"}))
    }else{
        w.Write(buildResponse(map[string]string{"statue": "true", "message": "成功"}))
    }
}

func logF(s string , err error){
    if err != nil {
		log.Fatal(s, err.Error())
	}else{
        log.Println(s)
    }
}

func buildResponse(temp map[string]string)([]byte){

    jsonResp, err := json.Marshal(temp)
	if err != nil {
		logF("Error happened in JSON marshal err:", err)
	}

    return jsonResp
}