package main

import (
  "fmt"
  "strings"
  "strconv"
  //"encoding/json"
)

type Salary struct {
  Lower  float32
  Upper  float32
}

func tokenize_string (s string) []string {
  return strings.Split(s, " ")
}

func IsNumeric(s string) bool {
   _, err := strconv.ParseFloat(s, 64)
   return err == nil
}

func string_to_float (s string) float32 {
  value, err := strconv.ParseFloat(s, 32)
  if err != nil {
    return float32(0)
  }
  return float32(value)
}

func Abs(x float32) float32 {
  if x < 0 {
    return -x
  }
  return x
}

func calculate_elements(operand1_s string ,operator string ,operand2_s string) float64 {
  operand1,_ := strconv.ParseFloat(operand1_s, 64) 
  operand2,_ := strconv.ParseFloat(operand2_s, 64)
  if operator == "+" {
    return operand1 + operand2
  }
  
  if operator == "-" {
    return operand1 - operand2
  }
  
  if operator == "*" {
    return operand1 * operand2
  }

  if operator == "/" {
    return operand1 / operand2
  }
  
  return 0
}

func get_salary (result string, bias string) Salary {
  var salary_lower = string_to_float(result) - Abs(string_to_float(bias))
  var salary_upper = string_to_float(result) + Abs(string_to_float(bias))
  s := Salary{salary_lower, salary_upper}

  return s
}

func process_postfix(postfix_string string) Salary {
  postfix_string = strings.TrimSpace(postfix_string)
  elements := tokenize_string(postfix_string)
  var stack Stack

  for i := 0; i < len(elements); i++ {
    if IsNumeric(elements[i]) {
      stack.Push(elements[i])
    } else {
      temp_operand1 := fmt.Sprintf("%v", stack.Pop())
      temp_operand2 := fmt.Sprintf("%v", stack.Pop())
      result := calculate_elements(temp_operand2,elements[i],temp_operand1)
      stack.Push(result)
    }
  }

  result := fmt.Sprintf("%v", stack.Pop())
  var bias = get_bias()

  
  var salary = get_salary(result,bias);

  return salary
}
