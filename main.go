package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Println(ps)
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func enableCors(w *http.ResponseWriter) {
  (*w).Header().Set("Access-Control-Allow-Origin", "http://karir.loc:3010")
}

func PostfixReader(writer http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    
    enableCors(&writer)
    spi := SalaryPredictionInput{}    

    body,err := ioutil.ReadAll(r.Body)

    

    if err != nil {
      writer.WriteHeader(400)
      fmt.Fprintf(writer, "%s", err)
      return
    }

    fmt.Println(string(body))

    err = json.Unmarshal([]byte(body), &spi)

    if err != nil {
      writer.WriteHeader(400)
      fmt.Fprintf(writer, "%s", err)
      return
    }
    
    var postfix_string = helper_for_convert_input_to_postfix_string(spi)

    var salary = process_postfix(postfix_string)

    json_salary, err := json.Marshal(salary)
 
    fmt.Println(string(json_salary))
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(http.StatusOK)
    writer.Write(json_salary)
    
}

func main() {
    MySQLDB = connect()
    defer MySQLDB.Close()

    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)
    router.POST("/postfix_reader", PostfixReader)

    
    fmt.Println("Successfully Running Machine Learning API")
    log.Fatal(http.ListenAndServe(":8080", router))
}
