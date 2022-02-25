package machete

import (
	"encoding/json"
	"os"
)

type routeInfo struct {
	Parent map[string]interface{}
	Name   string
}

type GenericJson struct {
	jsonMap map[string]interface{}
	routes  map[string]routeInfo
}

func NewGenericJson(filename string) (*GenericJson, error) {
	fileData, readFileErr := os.ReadFile(filename)

	if readFileErr != nil {
		return nil, readFileErr
	}

	var jsonMap map[string]interface{}

	unmErr := json.Unmarshal(fileData, &jsonMap)

	if unmErr != nil {
		return nil, unmErr
	}

	gj := GenericJson{jsonMap, make(map[string]routeInfo)}

	gj.readJsonMap(gj.jsonMap, "root")

	return &gj, nil
}

func (gj GenericJson) GetJsonMap() map[string]interface{} {
	return gj.jsonMap
}

func (gj *GenericJson) readJsonMap(m map[string]interface{}, name string) {
	for key, val := range m {

		switch element := val.(type) {
		case []interface{}:
			sample := element[0]

			// if sample is an object (map of interfaces), we want to mark element as REST route
			_, isMap := sample.(map[string]interface{})

			if isMap {
				gj.routes[key] = routeInfo{Parent: m, Name: key}
			}

		case map[string]interface{}:
			gj.readJsonMap(element, key)
		}
	}
}

func (gj *GenericJson) AddItem(routeName string, itemData map[string]interface{}) bool {
	routeInfo := gj.routes[routeName]

	// eg: cities[...]
	childElement := routeInfo.Parent[routeInfo.Name]

	if childElement != nil {
		itemList := childElement.([]interface{})
		routeInfo.Parent[routeInfo.Name] = append(itemList, itemData)

		// sample := itemList[0].(map[string]interface{})

		// attrs := identifyObjectAttributes(sample)

		// log.Println("attrs for", routeName, attrs)

		return true
	}

	return false
}

func (gj GenericJson) GetRoutes() []string {
	names := []string{}

	for key := range gj.routes {
		names = append(names, key)
	}

	return names
}

func (gj GenericJson) GetEntity(routeName string, params map[string]interface{}) []interface{} {

	// log.Printf("GetEntity: routeName=%s, params=%v", routeName, params)

	routeInfo := gj.routes[routeName]

	// eg: cities[...]
	childElement := routeInfo.Parent[routeInfo.Name]

	if childElement != nil {

		if len(params) == 0 {
			return childElement.([]interface{})
		} else {
			filteredList := []interface{}{}

			for _, element := range childElement.([]interface{}) {
				elementMap := element.(map[string]interface{})
				// log.Println(idx, elementMap)

				queryMatches := true

				for k, v := range params {
					// log.Printf("testing k=%v, v=%v against %v", k, v, elementMap[k])
					if elementMap[k] != v {
						queryMatches = false
						break
					}
				}

				if queryMatches {
					filteredList = append(filteredList, elementMap)
				}
			}
			return filteredList
		}
	}

	return nil
}

func (gj GenericJson) GetJson() ([]byte, error) {
	return json.Marshal(gj.jsonMap)
}

type attribute struct {
	Key  string
	Type string // TODO: change to a better type. eg THE type itself
}

func IdentifyObjectAttributes(obj map[string]interface{}) []attribute {

	var attrs []attribute

	for key, val := range obj {

		var valType string

		switch val.(type) {
		case bool:
			valType = "bool"
		case string:
			valType = "string"
		case float64:
			valType = "number"
		default:
			valType = "complex type"
		}

		attrs = append(attrs, attribute{key, valType})
	}

	return attrs
}
