package main

import (
  "fmt"
  "strings" 
  "github.com/fatih/structs" 
  //"strconv"
)


type SalaryPredictionInput struct {
  TipePerusahaan     string  `json:"tipe_perusahaan"`
  PendidikanTerakhir string  `json:"pendidikan_terakhir"`
  KategoriPekerjaan  string  `json:"kategori_pekerjaan"`
  TipeIndustri       string  `json:"tipe_industri"`
  JenjangKarir       string  `json:"jenjang_karir"`
  Location           string  `json:"location"`
  PengalamanKerja    float32 `json:"pengalaman_kerja"`
}


var TIPE_PERUSAHAAN = "tipe_perusahaan"
var PENDIDIKAN_TERAKHIR = "pendidikan_terakhir"
var LOCATION = "location"
var TIPE_INDUSTRI = "tipe_industri"
var KATEGORI_PEKERJAAN = "kategori_pekerjaan"
var JENJANG_KARIR = "jenjang_karir"
var PENGALAMAN_KERJA = "pengalaman_kerja"

var P_COLLECTION = []string{"tipe_perusahaan", "pendidikan_terakhir", "kategori_pekerjaan", "tipe_industri", "jenjang_karir", "location", "pengalaman_kerja"}

var W_TIPE_PERUSAHAAN = "w_tipe_perusahaan"
var W_PENDIDIKAN_TERAKHIR = "w_pendidikan_terakhir"
var W_LOCATION = "w_location"
var W_TIPE_INDUSTRI = "w_tipe_industri"
var W_KATEGORI_PEKERJAAN = "w_kategori_pekerjaan"
var W_JENJANG_KARIR = "w_jenjang_karir"
var W_PENGALAMAN_KERJA = "w_pengalaman_kerja"

var W_COLLECTION = []string{"w_tipe_perusahaan", "w_pendidikan_terakhir", "w_kategori_pekerjaan", "w_tipe_industri", "w_jenjang_karir", "w_location", "w_pengalaman_kerja"}

func get_predefine_value(tipe string,kunci string) string {
  var value float32
  err := MySQLDB.QueryRow("SELECT nilai FROM Salary_PredefinedData " + 
                          "WHERE tipe = '" + tipe + "'"+" AND " + 
                          "kunci = '" + kunci + "' " + 
                          "ORDER BY id DESC LIMIT 1").Scan(&value)

  if err != nil {
    fmt.Println(tipe,kunci)
    return string(0)
  }

  return fmt.Sprintf("%f", value)
}

func get_weight_of (kunci string) string {
  var value float32
  err := MySQLDB.QueryRow("SELECT nilai FROM Salary_Weighting " + 
                          "WHERE kunci = '" + kunci + "' " + 
                          "ORDER BY id DESC LIMIT 1").Scan(&value)

  if err != nil {
    panic(err.Error())
  }

  return fmt.Sprintf("%f", value)
}

func get_bias () string {
  var value float32
  err := MySQLDB.QueryRow("SELECT nilai FROM Salary_Weighting " + 
                          "WHERE kunci = 'bias' " + 
                          "ORDER BY id DESC LIMIT 1").Scan(&value)

  if err != nil {
    panic(err.Error())
  }

  return fmt.Sprintf("%f", value)
}

func map_pendidikan_terakhir (pendidikan_terakhir string) string {
  mapped_value := map[string]string{
                    "SLTA":"high school",
                    "Diploma": "associate",
                    "Sarjana/S1":"bachelors",
                    "Master/S2":"masters",
                  }
  return mapped_value[pendidikan_terakhir]  
}

func map_jenjang_karir (jenjang_karir string) string {
  mapped_value := map[string]string{
                    "Pemula / Staf":"entry level staff",
                    "Staf Senior": "senior staff",
                    "Supervisor":"Supervisor",
                    "Asisten Manajer":"assistant manager",
                    "Asisten Manajer Senior":"senior assistant manager",
                    "Manajer - Departemen":"manager department",
                    "Manajer - Cabang/Regional":"manager branch",
                    "Insinyur":"engineer",
                    "Manajer Senior": "senior manager",
                    "Asisten Wakil Presiden":"assistant vice president",
                    "General Manajer":"general manager",
                    "Kepala Unit Bisnis":"business unit head",
                    "Wakil Presiden":"vice president",
                    "Wakil Presiden Senior":"senior vice president",
                    "Wakil Presiden Eksekutif":"executive vp",
                    "Direktur": "director",
                    "Presiden Direktur - CEO":"president ceo",
                  }
  return mapped_value[jenjang_karir]  
}


func map_salary_prediction_input (spi SalaryPredictionInput) SalaryPredictionInput {
  spi.PendidikanTerakhir = map_pendidikan_terakhir(spi.PendidikanTerakhir)
  spi.JenjangKarir = map_jenjang_karir(spi.JenjangKarir)
  return spi
}

func helper_for_convert_input_to_postfix_string (spi SalaryPredictionInput) string {
  var postfix_string string
  var pengalaman_kerja_f string
  
  pengalaman_kerja_f = fmt.Sprintf("%f", spi.PengalamanKerja/12)

  spi = map_salary_prediction_input (spi)

  fmt.Println(spi)
  
  var spi_value = structs.Values(spi)

  err := MySQLDB.QueryRow("SELECT postfix FROM Salary_Formula ORDER BY id DESC LIMIT 1").Scan(&postfix_string)

  if err != nil {
    panic(err.Error())
  }

  postfix_string = strings.ToLower(postfix_string)

  for i:=0; i < len(P_COLLECTION) - 1 ; i++ {
    postfix_string = strings.Replace(postfix_string,P_COLLECTION[i],get_predefine_value(P_COLLECTION[i],spi_value[i].(string)),1)
  }

  postfix_string = strings.Replace(postfix_string,PENGALAMAN_KERJA,pengalaman_kerja_f,1)
  
  for i:= 0; i < len(W_COLLECTION) ; i++ {
    postfix_string = strings.Replace(postfix_string,W_COLLECTION[i],get_weight_of (P_COLLECTION[i]) ,1)
  }

  return postfix_string
}
