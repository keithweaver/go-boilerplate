package utils

import (
	// "os"
	// "errors"
	"fmt"
	// "encoding/json"
	// "io/ioutil"
	// "reflect"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"go-boilerplate/models"
)

type TranslationUtils struct {
	versionMapping map[string]VersionStep
}

type VersionStep struct {
	Inbound  func(body interface{}) (interface{}, error)
	Outbound func(body interface{}) (interface{}, error)
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

func NewInstanceOfTranslationUtils() TranslationUtils {
	return TranslationUtils{
		versionMapping: map[string]VersionStep{
			"LATEST": VersionStep{
				Previous: "2020-09-14",
			},
			"2020-09-24": VersionStep{
				Inbound:  convertUpdateCarV2ToFinal,
				Previous: "2020-09-24",
			},
			"2020-09-14": VersionStep{
				// "model": models.UpdateCarV1,
				Inbound: convertUpdateCarV1ToV2,
				Next:    "2020-09-24",
			},
		},
	}
}

func (u *TranslationUtils) InboundConversionByRequestBody(c *gin.Context) {
	// If not set, use LATEST
	version := ""

	// Get the version from the header
	if c.Request.Header.Get("VERSION") != "" {
		version = c.Request.Header.Get("VERSION")
		fmt.Printf("Setting version from header :: version :: %+v\n", version)
	}

	// If not set, then get it from the session
	sessionAsInterface, sessionExists := c.Get("session")
	if sessionExists && sessionAsInterface != nil && version == "" {
		session := sessionAsInterface.(models.Session)
		version = session.Version
		fmt.Printf("Setting version from session :: version :: %+v\n", version)
	}
	fmt.Printf("version :: %+v\n", version)
	// TODO - Add path
	// TODO - ctx.FullPath() - https://github.com/gin-gonic/gin/issues/1986

	// First version
	var body interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	// Translate till "next" doesnt exist using the "next"
	i := 0
	for i < 1000 {
		currentVersionObj, currentVersionExists := u.versionMapping[version]
		if !currentVersionExists {
			fmt.Println("Current version doesnt exist.")
			c.AbortWithStatusJSON(500, gin.H{"message": "error: internal server error"})
			return
		}
		fmt.Println("Version :: " + version)
		fmt.Printf("currentVersionObj :: %+v\n", currentVersionObj)
		fmt.Printf("currentVersionObj :: %T\n", currentVersionObj)
		nextVersion := currentVersionObj.Next
		if nextVersion == "" {
			break
		} else {
			// functionAsInterface :=
			fmt.Printf("body :: %+v\n", body)
			body, err := currentVersionObj.Inbound(body)
			// body, err := functionAsInterface.(func())(body)
			// body, err := u.Invoke(functionName, body)
			if err != nil {
				fmt.Printf("err on conversion :: %+v", err)
				c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
				return
			}

			fmt.Printf("body :: %+v\n", body)

			version = nextVersion
		}
		i += 1
	}

	fmt.Printf("body :: %T\n", body)
	fmt.Printf("body :: %+v\n", body)
	c.Set("body", body)
	c.Next()
}

func (u *TranslationUtils) Invoke(functionName string, body interface{}) (interface{}, error) {
	// Build a list of arguments for the function
	// args := make([]interface{}, 1)
	// args[0] = body
	// inputs := make([]reflect.Value, len(args))
	// for i, _ := range args {
	// 	fmt.Printf("args :: %+v\n", args[i])
	// 	inputs[i] = reflect.ValueOf(args[i])
	// }
	// fmt.Printf("inputs :: %+v\n", inputs)
	//
	// // Translate from current version to next version
	// // meth := reflect.ValueOf(u).MethodByName(functionName)
	// fmt.Printf("u :: %T\n", u)
	// fmt.Printf("u :: %+v\n", u)
	// class := reflect.ValueOf(u)
	// fmt.Printf("class :: %+v\n", class)
	// fmt.Printf("class :: %T\n", class)
	// fmt.Println("functionName :: " + functionName)
	//
	// TestMethod()
	// class := reflect.TypeOf(u)
	// method := (functionName)
	// method := class.MethodByName("TestMethod()")
	// fmt.Printf("method :: %+v\n", method)
	// fmt.Printf("method :: %T\n", method)
	// fmt.Printf("method.IsValid :: %+v\n", method.IsValid)
	// fmt.Printf("method.IsZero :: %+v\n", method.IsZero)

	// if !method.IsValid() {
	// 	fmt.Println("method is invalid")
	// 	return nil, errors.New("error: internal server error")
	// }
	//
	// res := method.Call(inputs)
	// res := ""

	// fmt.Printf("meth :: %T\n", meth)
	// fmt.Printf("res :: %+v\n", res)
	//
	// // Call method
	// // res := meth.Call(inputs)
	//
	// // Capture the response
	// ret := res[0].Interface()
	// // var err error
	// if v := res[1].Interface(); v != nil {
	// 	return nil, v.(error)
	// }
	// return ret, nil
	return nil, nil
}

func (u *TranslationUtils) InboundConversionByQueryParams(c *gin.Context) {
}

func (u *TranslationUtils) OutboundConversion(c *gin.Context) {
	// response, exists := c.Get("response")
	// if !exists {
	// 	c.JSON(500, gin.H{"message": "Internal server error"})
	// 	return
	// }
	//
	// statusCode, exists := c.Get("statusCode")
	// if !exists {
	// 	c.JSON(500, gin.H{"message": "Internal server error"})
	// 	return
	// }

	// Get the version from the header
	// If not set, then get it from the session
	// If not set, use LATEST
	// Translate till "previous" doesnt exist using the "previous"

	// c.JSON(500, gin.H{})
}

func (u *TranslationUtils) Nothing() {

}

func TestMethod() {
	fmt.Println("test method")
}

// Take UpdateCarV1 --> UpdateCarV2
func convertUpdateCarV1ToV2(body interface{}) (interface{}, error) {
	fmt.Printf("convertUpdateCarV1ToV2 :: %T\n", body)
	fmt.Printf("convertUpdateCarV1ToV2 :: %+v\n", body)
	bodyAsMap := body.(map[string]interface{})
	var updateCarV1 models.UpdateCarV2
	mapstructure.Decode(bodyAsMap, &updateCarV1)
	// updateCarV1 := body.(models.UpdateCarV1)
	return models.UpdateCarV2{
		Make:   updateCarV1.Make,
		Model:  "",
		Year:   updateCarV1.Year,
		Status: updateCarV1.Status,
	}, nil
}

// Take UpdateCarV2 --> UpdateCar (Change this for new fields)
func convertUpdateCarV2ToFinal(body interface{}) (interface{}, error) {
	bodyAsMap := body.(map[string]interface{})
	var updateCarV2 models.UpdateCarV2
	mapstructure.Decode(bodyAsMap, &updateCarV2)
	// updateCarV2 := bodyAsInterface.()
	return models.UpdateCar{
		Make:   updateCarV2.Make,
		Model:  updateCarV2.Model,
		Year:   updateCarV2.Year,
		Status: updateCarV2.Status,
	}, nil
}
