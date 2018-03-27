package della

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/spearfisher/cargoUpdater/utils"
)

const dellaURI = "https://della.ua"

// Client is a della http middleware
type Client struct {
	login      string
	password   string
	dasc       string
	phpsessid  string
	httpClient *http.Client
}

// NewDellaClient returns new authorized della client instance
func NewDellaClient(login string, password string) *Client {
	dellaClient := &Client{login: login, password: password}
	dellaClient.httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	dellaClient.auth()
	return dellaClient
}

func (c *Client) auth() {
	authData := fmt.Sprintf("login_mode=enter&location_url=della.ua/&login=%s&password=%s", c.login, c.password)

	req, _ := http.NewRequest("POST", dellaURI, strings.NewReader(authData))
	setStandartHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		utils.Logger.Println(err)
		return
	}
	c.setCookies(resp.Cookies())

	if c.dasc == "" {
		utils.Logger.Println(fmt.Sprintf("Can't authrize with credentials. Login: %s, Password: %s", c.login, c.password))
		utils.Logger.Println(resp)
	}
}

// GetList returns list of freights
func (c *Client) GetList() (*CargosData, error) {
	req, _ := http.NewRequest("GET", dellaURI+"/my/", nil)
	setStandartHeaders(req)

	req.AddCookie(&http.Cookie{
		Name:  "dasc",
		Value: c.dasc,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		utils.Logger.Println(err)
	}

	c.setCookies(resp.Cookies())

	return parseData(resp.Body)
}

func (c *Client) setCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Name == "dasc" && len(cookie.Value) > 20 {
			c.dasc = cookie.Value
		}
		if cookie.Name == "PHPSESSID" {
			c.phpsessid = cookie.Value
		}
	}
}

// RefreshCargos refreshes list of cargos for della user
func (c *Client) RefreshCargos(data *CargosData) {
	crazyFormatedDellaBody := "mode=repeat_requests_v2&codes="
	for _, id := range data.Ids {
		crazyFormatedDellaBody += id + "_today%20"
	}
	crazyFormatedDellaBody += "&menu_id=22&flag=act&dateups=" + data.Dateups

	path := fmt.Sprintf("%s/new_ajax.php?PHPSESSID=%s&JsHttpRequest=14386040738230-xml", dellaURI, c.phpsessid)
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer([]byte(crazyFormatedDellaBody)))
	setStandartHeaders(req)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Referer", "//della.ua/my/")

	req.AddCookie(&http.Cookie{
		Name:  "PHPSESSID",
		Value: c.phpsessid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "dasc",
		Value: c.dasc,
	})

	_, err := c.httpClient.Do(req)
	if err != nil {
		utils.Logger.Println(err)
	}
}

func setStandartHeaders(req *http.Request) {
	req.Header.Set("Host", "della.ua")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:51.0) Gecko/20100101 Firefox/51.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", dellaURI)
	req.Header.Set("Cookie", "della_request_v2=YToxOntzOjM6IkNJRCI7czoxOToiMTcxNjExNzM5MDkwNTUwNjUzMCI7fQ%3D%3D; _pk_id.10.9c5c=4ee7339230209e73.1497191951.1.1497191951.1497191951.; _pk_ses.10.9c5c=*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Language", "ru,en-US;q=0.8,en;q=0.6,uk;q=0.4")
}
