package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func length(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/landingPage.html"))

	tmpl.Execute(w, "length")
}

func weight(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/landingPage.html"))

	tmpl.Execute(w, "weight")
}

func temperature(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/landingPage.html"))

	tmpl.Execute(w, "temperature")
}

func result(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/result.html"))

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
	}

	sUnit := r.FormValue("unit")
	fUnit, err := strconv.ParseFloat(sUnit, 64)
	if err != nil {
		http.Error(w, "Enter a valid number", http.StatusBadRequest)
	}
	convFrom := strings.ToLower(r.FormValue("cfrom"))
	convTo := strings.ToLower(r.FormValue("cto"))

	var conversion float64
	switch convFrom {
	case "mm", "cm", "m", "km", "in", "ft", "yd", "mi":
		conversion = convLength(fUnit, convFrom, convTo)
	case "mg", "g", "kg", "oz", "lb":
		conversion = convWeight(fUnit, convFrom, convTo)
	case "C", "F", "K":
		conversion = convTemperature(fUnit, convFrom, convTo)
	default:
		http.Error(w, "Invalid unit conversion request", http.StatusBadRequest)
		return
	}

	final := fmt.Sprintf("%.2f %s = %.2f %s", fUnit, convFrom, conversion, convTo)
	tmpl.Execute(w, final)
}

func convLength(num float64, convFrom, convTo string) float64 {
	lengthConverter := map[string]float64{
		"mm": .001,
		"cm": .01,
		"m":  1,
		"km": 1000,
		"in": .0254,
		"ft": .3048,
		"yd": .9144,
		"mi": 1608.34,
	}
	return (num * lengthConverter[convFrom]) / lengthConverter[convTo]
}

func convWeight(num float64, convFrom, convTo string) float64 {
	weightConverter := map[string]float64{
		"mg": .001,
		"g":  1,
		"kg": 1000,
		"oz": 28.3495,
		"lb": 453.592,
	}
	return (num * weightConverter[convFrom]) / weightConverter[convTo]
}

func convTemperature(num float64, convFrom, convTo string) float64 {
	convFrom = strings.ToUpper(convFrom)
	convTo = strings.ToUpper(convTo)
	if convFrom == convTo {
		return num
	}
	if convFrom == "C" {
		if convTo == "F" {
			return (num * 9 / 5) + 32
		}
		return num + 273.15
	} else if convFrom == "F" {
		if convTo == "C" {
			return (num - 32) * 5 / 9
		}
		return (num-32)*5/9 + 273.15
	} else {
		if convFrom == "C" {
			return num - 273.15
		}
		return (num-273.15)*9/5 + 32
	}
}

func main() {

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	http.HandleFunc("/length", length)
	http.HandleFunc("/weight", weight)
	http.HandleFunc("/temperature", temperature)
	http.HandleFunc("/result", result)
	http.HandleFunc("/", length)

	http.ListenAndServe(":8080", nil)
}
