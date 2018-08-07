package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	activity_proto "server/activity-srv/proto/activity"
	content_proto "server/content-srv/proto/content"
	static_proto "server/static-srv/proto/static"
	"strconv"
	"strings"

	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
)

// type const
const (
	TypeGeo            = "geo"
	TypeGeoVariable    = "geoVariable"
	TypeGeoPost        = "post"
	TypeGeoAddress     = "address"
	TypeObject         = "object"
	TypeString         = "string"
	TypeStringArray    = "stringArray"
	TypeVariable       = "variable"
	TypeStringAddress  = "stringAddress"
	TypeStringAddress1 = "stringAddress1"
)

// path const
const (
	PathId        = "id"
	PathActivity  = "activity"
	PathOffer     = "offer"
	PathLocation  = "location"
	PathStartDate = "start_date"
	PathEndDate   = "end_date"
)

var (
	// PostCodes store geo location by postcode key
	PostCodes = make(map[string]*content_proto.Geo)
)

// type Transform struct {
// 	Type     string `json:"type"`
// 	Path     string `json:"path"`
// 	Variable string `json:"variable"`
// }

type ExtBackend struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	AppURL       string                      `json:"appURL"`
	Next         string                      `json:"next"`
	Weight       int64                       `json:"weight"`
	TimeInterval int64                       `json:"timeInterval"`
	Enabled      bool                        `json:"enabled"`
	Transform    []*activity_proto.Transform `json:"transform"`

	Finished bool
}

// ResExt is struct of external apis response
type ResExt struct {
	Next  string `json:"next,omitempty"`
	Items []struct {
		ID    interface{}            `json:"id"`
		State string                 `json:"state"`
		Data  map[string]interface{} `json:"data,omitempty"`
	} `json:"items,omitempty"`
	Data []map[string]interface{} `json:"data,omitempty"`
}

// InitExternal initializes the client of better external api
func InitExternal(client *ExtBackend) {
	Backends[client.Name] = client
	Weights[client.Name] = client.Weight
}

// Query external apis, return data array
func (b *ExtBackend) Query(source string) ([]*content_proto.Activity, error) {
	// making url
	url := b.AppURL + b.Next
	if strings.HasPrefix(b.Next, "https://") || strings.HasPrefix(b.Next, "http://") {
		url = b.Next
	}

	log.Info("external request url: ", url)
	// send request to external service
	resp, err := grequests.Get(url, nil)
	if err != nil {
		log.Error("Request error:", err)
		return nil, err
	}
	// parsing response
	res := &ResExt{}
	err = resp.JSON(res)
	if err != nil {
		log.Error("Json parsing err:", err)
		return nil, err
	}

	// replace next field
	if b.Next == res.Next || (len(res.Items) == 0 && len(res.Data) == 0) {
		b.Finished = true
	} else {
		b.Next = res.Next
	}

	datas := make([]*content_proto.Activity, 0)
	if len(res.Items) > 0 {
		for _, item := range res.Items {
			if item.Data == nil || item.State != "updated" {
				continue
			}
			// transform
			transformFields(item.Data, b.Transform)

			if activity := activityFromJSON(item.Data); activity != nil {
				activity.Id = fmt.Sprintf("%v", item.ID)
				activity.Source = source

				datas = append(datas, activity)
			}
		}
	} else if len(res.Data) > 0 {
		for _, d := range res.Data {
			// transform
			transformFields(d, b.Transform)

			if activity := activityFromJSON(d); activity != nil {
				activity.Source = source
				datas = append(datas, activity)
			}
		}
	}

	return datas, nil
}

// itemFromJSON method
func activityFromJSON(data map[string]interface{}) *content_proto.Activity {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error("Json parsing with real object err:", err)
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(b))
	var activity content_proto.Activity
	err = decoder.Decode(&activity)
	if err == nil {
		return &activity
	}
	log.Error("Json parsing with real object err:", err)
	return nil
}

func transformFields(data map[string]interface{}, transform []*activity_proto.Transform) {
	for _, t := range transform {
		switch t.Type {
		case TypeGeo, TypeGeoVariable, TypeGeoPost, TypeGeoAddress:
			retrieveGeoObject(data, t)
		default:
			retrieveObject(data, t)
		}
	}
}

func transformVariable(data map[string]interface{}, o interface{}, t *activity_proto.Transform) {
	variables := strings.Split(t.Variable, "/")
	d := data
	for i, v := range variables {
		if i == len(variables)-1 {
			d[v] = o
		} else {
			if _, ok := d[v]; ok {
				if d[v] == nil {
					continue
				}
				d = d[v].(map[string]interface{})
			} else {
				d[v] = new(map[string]interface{})
			}
		}
	}
}

func retrieveGeoObject(data map[string]interface{}, t *activity_proto.Transform) {
	var geo *content_proto.Geo
	// parsing geo path
	fields := strings.Split(t.Path, "/")
	d := data
	r := ""
	for _, field := range fields {
		if d[field] == nil {
			continue
		}
		switch d[field].(type) {
		case string:
			r = d[field].(string)
		default:
			d = d[field].(map[string]interface{})
		}
	}

	switch t.Type {
	case TypeGeo:
		if d["latitude"] != nil && d["longitude"] != nil {
			geo = new(content_proto.Geo)
			geo.Latitude = d["latitude"].(float64)
			geo.Longitude = d["longitude"].(float64)
		}
	case TypeGeoVariable:
		geo = new(content_proto.Geo)
		geo.Latitude, _ = strconv.ParseFloat(d["latitude"].(string), 64)
		geo.Longitude, _ = strconv.ParseFloat(d["longitude"].(string), 64)
	case TypeGeoPost:
		postcode := r
		geo = retreiveGeolocationFromPostcode(postcode)
	case TypeGeoAddress:
		address := r
		postcode := retrievePostcodeFromAddress(address)
		geo = retreiveGeolocationFromPostcode(postcode)
	}

	transformVariable(data, geo, t)
}

func retrieveObject(data map[string]interface{}, t *activity_proto.Transform) {
	var result interface{}

	switch t.Type {
	case TypeObject:
		object := data[t.Path]
		if object != nil {
			result = []interface{}{object}
		} else {
			result = []interface{}{}
		}
	case TypeString:
		object := data[t.Path]
		if object == nil {
			object = ""
		}
		str := object.(string)
		{
			switch t.Path {
			case PathActivity:
				result = []*content_proto.ActivityAttr{{PrefLabel: str}}
			case PathOffer:
				result = []*content_proto.Offer{{Name: str}}
			default:
				result = str
			}
		}
	case TypeStringArray:
		strs := data[t.Path].([]interface{})
		{
			switch t.Path {
			case PathActivity:
				activities := []*content_proto.ActivityAttr{}
				for _, str := range strs {
					activities = append(activities, &content_proto.ActivityAttr{PrefLabel: str.(string)})
				}
				result = activities
			case PathOffer:
				offers := []*content_proto.Offer{}
				for _, str := range strs {
					offers = append(offers, &content_proto.Offer{Name: str.(string)})
				}
				result = offers
			}
		}
	case TypeVariable:
		{
			switch t.Path {
			case PathId:
				result = fmt.Sprintf("%v", data[t.Path])
			case PathOffer:
				objects := data[t.Path].([]interface{})
				for _, object := range objects {
					obj := object.(map[string]interface{})
					for k, v := range obj {
						switch v.(type) {
						case string:
							obj[k] = v
						default:
							obj[k] = fmt.Sprintf("%v", v)
						}
					}
				}
				result = objects
			}
		}
	case TypeStringAddress:
		{
			if data[t.Path] != nil {
				log.Warn(t.Path)
				log.Warn(data[t.Path])
				address := data[t.Path].(string)
				postcode := retrievePostcodeFromAddress(address)
				geo := retreiveGeolocationFromPostcode(postcode)
				result = content_proto.Place{
					Geo: geo,
					Address: &static_proto.Address{
						StreetAddress: address,
					},
				}
			}
		}
	case TypeStringAddress1:
		{
			fields := strings.Split(t.Path, "/")
			d := data
			r := ""
			for _, field := range fields {
				if d[field] == nil {
					continue
				}
				switch d[field].(type) {
				case string:
					r = d[field].(string)
				default:
					d = d[field].(map[string]interface{})
				}
			}
			if len(r) > 0 {
				address := r
				result = static_proto.Address{
					StreetAddress: address,
				}
			}
		}
	}

	transformVariable(data, result, t)
}

// reference https://regex101.com/
func retrievePostcodeFromAddress(address string) string {
	re := regexp.MustCompile("((GIR 0AA)|((([A-PR-UWYZ][0-9][0-9]?)|(([A-PR-UWYZ][A-HK-Y][0-9][0-9]?)|(([A-PR-UWYZ][0-9][A-HJKSTUW])|([A-PR-UWYZ][A-HK-Y][0-9][ABEHMNPRVWXY])))) [0-9][ABD-HJLNP-UW-Z]{2}))")
	match := re.FindStringSubmatch(address)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}

// reference http://postcodes.io/
func retreiveGeolocationFromPostcode(postcode string) *content_proto.Geo {
	if len(postcode) == 0 {
		return nil
	}
	// check stored geolocation
	if geo, ok := PostCodes[postcode]; ok {
		return geo
	}
	// make request to retrieve geo-location with postcode
	url := "http://api.postcodes.io/postcodes/" + postcode
	// log.Error("postcode.io retrieve url:", url)

	// response struct of
	type postio struct {
		Status int                `json:"status"`
		Result *content_proto.Geo `json:"result"`
	}
	// send request to external service
	resp, err := grequests.Get(url, nil)
	if err != nil {
		log.Error("Request error:", err)
		return nil
	}
	// parsing response
	res := &postio{}
	err = resp.JSON(res)
	if err != nil {
		log.Error("Json parsing err:", err)
		return nil
	}

	PostCodes[postcode] = res.Result
	return res.Result
}
