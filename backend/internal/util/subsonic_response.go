package util

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

const SubsonicAPIVersion = "1.16.1"

type SubsonicResponse struct {
	XMLName xml.Name    `xml:"subsonic-response" json:"-"`
	Status  string      `xml:"status,attr" json:"status"`
	Version string      `xml:"version,attr" json:"version"`
	Type    string      `xml:"type,attr" json:"type"`
	Body    interface{} `xml:",any" json:"-"`
}

type SubsonicError struct {
	XMLName xml.Name `xml:"error" json:"-"`
	Code    int      `xml:"code,attr" json:"code"`
	Message string   `xml:"message,attr" json:"message"`
}

type SubsonicJSONWrapper struct {
	Response interface{} `json:"subsonic-response"`
}

func SubsonicOK(c *gin.Context, body interface{}) {
	format := c.DefaultQuery("f", "xml")

	if format == "json" {
		resp := gin.H{
			"status":  "ok",
			"version": SubsonicAPIVersion,
			"type":    "pineapple-music",
		}
		if body != nil {
			// Merge body into response
			if m, ok := body.(gin.H); ok {
				for k, v := range m {
					resp[k] = v
				}
			} else {
				b, _ := json.Marshal(body)
				var m map[string]interface{}
				json.Unmarshal(b, &m)
				for k, v := range m {
					resp[k] = v
				}
			}
		}
		c.JSON(http.StatusOK, SubsonicJSONWrapper{Response: resp})
		return
	}

	// XML response
	type xmlResponse struct {
		XMLName xml.Name `xml:"subsonic-response"`
		Status  string   `xml:"status,attr"`
		Version string   `xml:"version,attr"`
		Type    string   `xml:"type,attr"`
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, buildSubsonicXML("ok", body))
}

func SubsonicErrorResp(c *gin.Context, code int, message string) {
	format := c.DefaultQuery("f", "xml")

	if format == "json" {
		c.JSON(http.StatusOK, SubsonicJSONWrapper{
			Response: gin.H{
				"status":  "failed",
				"version": SubsonicAPIVersion,
				"type":    "pineapple-music",
				"error": gin.H{
					"code":    code,
					"message": message,
				},
			},
		})
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	xmlStr := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<subsonic-response status="failed" version="` + SubsonicAPIVersion + `" type="pineapple-music">` +
		`<error code="` + intToStr(code) + `" message="` + xmlEscape(message) + `"/>` +
		`</subsonic-response>`
	c.String(http.StatusOK, xmlStr)
}

func buildSubsonicXML(status string, body interface{}) string {
	header := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<subsonic-response status="` + status + `" version="` + SubsonicAPIVersion + `" type="pineapple-music">`
	footer := `</subsonic-response>`

	if body == nil {
		return header + footer
	}

	inner, err := xml.Marshal(body)
	if err != nil {
		return header + footer
	}

	return header + string(inner) + footer
}

func intToStr(n int) string {
	s := ""
	if n == 0 {
		return "0"
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func xmlEscape(s string) string {
	var result string
	for _, c := range s {
		switch c {
		case '&':
			result += "&amp;"
		case '<':
			result += "&lt;"
		case '>':
			result += "&gt;"
		case '"':
			result += "&quot;"
		case '\'':
			result += "&apos;"
		default:
			result += string(c)
		}
	}
	return result
}
