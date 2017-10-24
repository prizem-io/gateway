package errorfactory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
	//"golang.org/x/text/language"
)

type (
	// Params is used to replace placeholders in an error template with the corresponding values.
	Params map[string]interface{}

	errorTemplate struct {
		Status           int    `yaml:"status"`
		Message          string `yaml:"message"`
		DeveloperMessage string `yaml:"developerMessage"`
		ErrorCode        string `yaml:"errorCode"`
	}

	APIError struct {
		XMLName          struct{}    `json:"-" xml:"error" msgpack:"-"`
		Status           int         `json:"status" xml:"status" msgpack:"status"`
		Reason           string      `json:"reason" xml:"reason" msgpack:"reason"`
		Message          string      `json:"message" xml:"message" msgpack:"message"`
		DeveloperMessage string      `json:"developerMessage,omitempty" xml:"developerMessage,omitempty" msgpack:"developerMessage,omitempty"`
		ErrorCode        string      `json:"errorCode,omitempty" xml:"errorCode,omitempty" msgpack:"errorCode,omitempty"`
		RequestID        string      `json:"requestId,omitempty" xml:"requestId,omitempty" msgpack:"requestId,omitempty"`
		Details          interface{} `json:"details,omitempty" xml:"details,omitempty" msgpack:"details,omitempty"`
	}
)

type Localized interface {
	Locale() []string //language.Tag
}

type ErrorTemplateMap map[string]errorTemplate
type LocaleMap map[string]ErrorTemplateMap

var (
	localeMap          = make(LocaleMap, 10)
	defaultTemplateMap ErrorTemplateMap
	baseFilePrefix     string
	lock               sync.RWMutex

	defaultError = APIError{
		Status:  http.StatusInternalServerError,
		Message: "Boom",
	}
)

func Initialize(messageBase string) error {
	baseFilePrefix = messageBase
	errorTemplates, err := loadMessages("")
	if err != nil {
		return errors.New(fmt.Sprint("Failed to read default locale messages: ", err))
	}

	defaultTemplateMap = errorTemplates
	localeMap[""] = errorTemplates

	return nil
}

func loadMessages(locale string) (ErrorTemplateMap, error) {
	bytes, err := ioutil.ReadFile(baseFilePrefix + locale + ".yaml")
	if err != nil {
		return nil, err
	}
	messages := ErrorTemplateMap{}
	return messages, yaml.Unmarshal(bytes, &messages)
}

func getErrorTemplateMap(localized Localized) ErrorTemplateMap {
	locale := localized.Locale()
	language := locale[0] //.Language()
	country := locale[1]  //.Country()
	variant := locale[2]  //.Variant()

	if len(country) > 0 && len(variant) > 0 {
		key := language + "_" + country + "_" + variant
		errorTemplateMap := getForKey(key)
		if errorTemplateMap != nil {
			return errorTemplateMap
		}
	} else if len(country) > 0 {
		key := language + "_" + country
		errorTemplateMap := getForKey(key)
		if errorTemplateMap != nil {
			return errorTemplateMap
		}
	} else {
		errorTemplateMap := getForKey(language)
		if errorTemplateMap != nil {
			return errorTemplateMap
		}
	}

	return defaultTemplateMap
}

func getForKey(key string) ErrorTemplateMap {
	lock.RLock()
	errorTemplateMap, ok := localeMap[key]
	lock.RUnlock()

	if !ok {
		errorTemplateMap, _ := loadMessages(key)
		lock.Lock()
		localeMap[key] = errorTemplateMap
		lock.Unlock()
	}

	return errorTemplateMap
}

func New(locale Localized, reason string, params ...Params) *APIError {
	errorTemplate := errorTemplate{}
	errorTemplates := getErrorTemplateMap(locale)
	found := commonTemplates(errorTemplates, &errorTemplate, &reason, params...)

	if found {
		return errorTemplate.newHTTPError(&reason, params...)
	}

	return &defaultError
}

func FromTemplate(locale Localized, reason string, template string, params ...Params) *APIError {
	errorTemplate := errorTemplate{}
	errorTemplates := getErrorTemplateMap(locale)
	found := commonTemplates(errorTemplates, &errorTemplate, &reason, params...)

	// Try to find a template using the reason and template together
	if template != "" {
		reasonTemplate := reason + "|" + template
		if template, ok := errorTemplates[reasonTemplate]; ok {
			found = true
			copyErrorTemplate(&template, &errorTemplate)
		}
	}

	if found {
		return errorTemplate.newHTTPError(&reason, params...)
	}

	return &defaultError
}

func commonTemplates(errorTemplates ErrorTemplateMap, errorTemplate *errorTemplate, reason *string, params ...Params) bool {
	found := false
	numParams := len(params)

	// Default to using just the reason as the key
	if template, ok := errorTemplates[*reason]; ok {
		found = true
		copyErrorTemplate(&template, errorTemplate)
	}

	// Try to find a template using the reason and template together
	if numParams > 0 {
		paramSize := len(params[0])
		if paramSize > 0 {
			reasonParams := *reason + "|" + strconv.Itoa(paramSize)
			if template, ok := errorTemplates[reasonParams]; ok {
				found = true
				copyErrorTemplate(&template, errorTemplate)
			}
		}
	}

	return found
}

func copyErrorTemplate(src *errorTemplate, dest *errorTemplate) {
	if src.Status != 0 {
		dest.Status = src.Status
	}

	if src.Message != "" {
		dest.Message = src.Message
	}

	if src.DeveloperMessage != "" {
		dest.DeveloperMessage = src.DeveloperMessage
	}

	if src.ErrorCode != "" {
		dest.ErrorCode = src.ErrorCode
	}
}

func (e errorTemplate) newHTTPError(reason *string, params ...Params) *APIError {
	message := e.Message
	developerMessage := e.DeveloperMessage
	if len(params) > 0 {
		message = replaceMessagePlaceholders(message, params[0])

		if developerMessage != "" {
			developerMessage = replaceMessagePlaceholders(developerMessage, params[0])
		}
	}

	return &APIError{
		Status:           e.Status,
		Reason:           *reason,
		Message:          message,
		DeveloperMessage: developerMessage,
		ErrorCode:        e.ErrorCode,
	}
}

func replaceMessagePlaceholders(message string, params Params) string {
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.Message
}
