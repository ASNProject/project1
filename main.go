package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	//"html/template"
	"io/ioutil"
	"log"
	"net/http"

	//"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Dataiot struct {
	Id string `json:"id"`
	Command string `json: "command"`
	Temp string `json: "temp"`
	Hum string `json: "hum"`
}

func connect() (*sql.DB, error){
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db_webservice")
	if err != nil{
		return nil, err
	}
	return db, nil
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Welcome to the HomePage")
	fmt.Println("Enpoint Hit: homePage")
}

func esp_sensor(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var d Dataiot
	json.Unmarshal(reqBody, &d)
	json.NewEncoder(w).Encode(d)

	db, err :=connect()
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO tb_sensor (temp)VALUES (?)", d.Temp)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fmt.Println("Endpoint Hit: esp_sensor, Time: ", time.Now())
}

func web_sensor(w http.ResponseWriter, r *http.Request){
	db, err := connect()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	data, err := db.Query("SELECT * FROM tb_sensor ORDER BY id DESC LIMIT 1")
	if err != nil{
		panic(err.Error())
	}
	defer data.Close()

	var result[] Dataiot
	for data.Next(){
		var each = Dataiot{}
		var err = data.Scan(&each.Id, &each.Temp)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}

	json.NewEncoder(w).Encode(result)
	fmt.Println("Endpoint Hit: web_sensor, Time: ", time.Now())
}

func web_command(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var d Dataiot
	json.Unmarshal(reqBody, &d)
	json.NewEncoder(w).Encode(d)

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO tb_command (command)VALUES (?)", d.Command)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fmt.Println("Endpoint Hit: web_command, Time: ",time.Now())
}

func esp_command(w http.ResponseWriter, r *http.Request){
	db, err := connect()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	data, err := db.Query("SELECT * FROM tb_command ORDER BY id DESC LIMIT 1")
	if err != nil{
		panic(err.Error())
	}
	defer data.Close()

	var result[] Dataiot
	for data.Next(){
		var each = Dataiot{}
		var err = data.Scan(&each.Id, &each.Command)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}

	json.NewEncoder(w).Encode(result)
	fmt.Println("Endpoint Hit: esp_command, Time: ", time.Now())
}



func handlerRequest(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/esp_sensor", esp_sensor).Methods("POST")
	myRouter.HandleFunc("/web_command", web_command).Methods("POST")
	myRouter.HandleFunc("/web_sensor", web_sensor).Methods("GET")
	myRouter.HandleFunc("/esp_command", esp_command).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

func main(){
	fmt.Println("Program Started!!!")
	handlerRequest()
/*	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		var filepath = path.Join("views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]interface{}{
			"title" : "Beta Golang Web Project", 
			"name" : "Sembada tech",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost: 8000")
	http.ListenAndServe(":8000", nil) */
}