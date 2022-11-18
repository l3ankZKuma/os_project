package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

type page struct {
	Header    string
	Temp      float64
	Long      float64
	Lat       float64
	Describe  string
	Icon      string
	Humidity  int
	DtTxt1    string
	DtTxt2    string
	DtTxt3    string
	DtTxt4    string
	Icon1     string
	Icon2     string
	Icon3     string
	Icon4     string
	Describe1 string
	Describe2 string
	Describe3 string
	Describe4 string
}
type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}
type forecastData struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int     `json:"visibility"`
		Pop        float64 `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Rain  struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Celcius  float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`

	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}
func queryForcast(city string) (forecastData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return forecastData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city + "&units=metric")
	if err != nil {
		return forecastData{}, err
	}

	defer resp.Body.Close()

	var d forecastData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return forecastData{}, err
	}
	return d, nil
}
func queryWeather(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city + "&units=metric")
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

var tmp1 *template.Template

func init() {
	tmp1 = template.Must(template.ParseGlob("templates/*.html"))
}
func main() {
	// fs := http.FileServer(http.Dir("assets"))

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			tmp1.ExecuteTemplate(w, "main.html", nil)
		})
	// http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := queryWeather(city)
			dataForecast, _ := queryForcast(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//w.Header().Set("Content-Type", "application/json; charset=utf-8")
			//json.NewEncoder(w).Encode(data)
			web_icon := " http://openweathermap.org/img/wn/" + data.Weather[0].Icon + "@2x.png"
			web_icon1 := " http://openweathermap.org/img/wn/" + dataForecast.List[8].Weather[0].Icon + "@2x.png"
			web_icon2 := " http://openweathermap.org/img/wn/" + dataForecast.List[16].Weather[0].Icon + "@2x.png"
			web_icon3 := " http://openweathermap.org/img/wn/" + dataForecast.List[24].Weather[0].Icon + "@2x.png"
			web_icon4 := " http://openweathermap.org/img/wn/" + dataForecast.List[32].Weather[0].Icon + "@2x.png"

			p := page{Header: data.Name, Temp: data.Main.Celcius,
				Long: data.Coord.Lon, Lat: data.Coord.Lat,
				Describe:  data.Weather[0].Description,
				Icon:      web_icon,
				Humidity:  data.Main.Humidity,
				DtTxt1:    dataForecast.List[8].DtTxt[:11],
				DtTxt2:    dataForecast.List[16].DtTxt[:11],
				DtTxt3:    dataForecast.List[24].DtTxt[:11],
				DtTxt4:    dataForecast.List[32].DtTxt[:11],
				Icon1:     web_icon1,
				Icon2:     web_icon2,
				Icon3:     web_icon3,
				Icon4:     web_icon4,
				Describe1: dataForecast.List[8].Weather[0].Description,
				Describe2: dataForecast.List[16].Weather[0].Description,
				Describe3: dataForecast.List[24].Weather[0].Description,
				Describe4: dataForecast.List[32].Weather[0].Description,
			}
			t, err := template.ParseFiles("templates/weather.html")
			fmt.Print(err)
			t.Execute(w, p)
		})

	http.ListenAndServe(":8080", nil)
}
