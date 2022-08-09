package starmap

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type POI struct {
	Name     string  `json:"Name"`
	X        float64 `json:"X"`
	Y        float64 `json:"Y"`
	Z        float64 `json:"Z"`
	QW       float64 `json:"qw"`
	QX       float64 `json:"qx"`
	QY       float64 `json:"qy"`
	QZ       float64 `json:"qz"`
	QTMarker bool    `json:"QTMarker"`
}

type Container struct {
	Name           string  `json:"Name"`
	Container      string  `json:"Container"`
	X              float64 `json:"X"`
	Y              float64 `json:"Y"`
	Z              float64 `json:"Z"`
	QW             float64 `json:"qw"`
	QX             float64 `json:"qx"`
	QY             float64 `json:"qy"`
	QZ             float64 `json:"qz"`
	QTMarker       bool    `json:"QTMarker"`
	OMRadius       float64 `json:"OM Radius"`
	BodyRadius     float64 `json:"Body Radius"`
	ArrivalRadius  float64 `json:"Arrival Radius"`
	TimeLines      float64 `json:"Time Lines"`
	RotationSpeed  float64 `json:"Rotation Speed"`
	RotationAdjust float64 `json:"Rotation Adjust"`
	OrbitalRadius  float64 `json:"Orbital Radius"`
	OrbitalSpeed   float64 `json:"Orbital Speed"`
	OrbitalAngle   float64 `json:"Orbital Angle"`
	GridRadius     float64 `json:"Grid Radius"`
	POI            POI     `json:"POI"`
}

var Containers map[string]*Container = map[string]*Container{}
var Keys []string = []string{}

func Load() error {
	jf, err := os.Open("./database.json")
	if err != nil {
		return err
	}
	defer jf.Close()

	b, _ := ioutil.ReadAll(jf)

	var res map[string]interface{}
	err = json.Unmarshal([]byte(b), &res)
	if err != nil {
		return err
	}

	for _, v := range res["Containers"].(map[string]interface{}) {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		c := &Container{}
		if err := json.Unmarshal(b, c); err != nil {
			return err
		}

		Containers[c.Name] = c
	}

	for k := range Containers {
		Keys = append(Keys, k)
	}

	return nil
}

func KeysUpper() []string {
	r := []string{}
	for _, v := range Keys {
		r = append(r, strings.ToUpper(v))
	}
	return r
}
