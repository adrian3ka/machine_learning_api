package main
import (
  //"fmt"
  "strconv"
  "github.com/ghodss/yaml"
  "io/ioutil"
  "os"
  "encoding/json"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)


var MySQLDB *sql.DB

type Conf struct {
  Database string `json:"database"`
  Host     string `json:"host"`
  Password string `json:"password"`
  Port     int    `json:"port"`
  Username string `json:"username"`
}


func getEnv() string {
  return os.Args[1]
}

func getConf() Conf{
  c := Conf{}

  yamlFile, err := ioutil.ReadFile("config/"+getEnv()+".yaml")
  
  config, err := yaml.YAMLToJSON(yamlFile)
  if err != nil {
    panic (err)
  }

  err = json.Unmarshal(config, &c)
  if err != nil {
    panic (err)
  }

  return c
}


func connect() *sql.DB {

  config := getConf()

  port := strconv.Itoa(config.Port)
  db, err := sql.Open("mysql",config.Username+":" + config.Password + "@tcp" + 
                      "(" + config.Host + ":" + port + ")/"+ config.Database)

  if err != nil {
    panic (err)
  }

  err = db.Ping()
  if err != nil {
    panic (err)
  }

  return db
}
