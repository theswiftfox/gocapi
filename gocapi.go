package gocapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseURI = "https://addons-ecs.forgesvc.net/api/v2/"

// GameInfo wraps the important values of the gameslist call to curseapi
type GameInfo struct {
	GameID int    `json:"id"`
	Name   string `json:"name"`
}

// AddonInfo wraps the important values for addons fetched from curseapi
type AddonInfo struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Website string      `json:"websiteUrl"`
	Files   []AddonFile `json:"latestFiles"`
}

// AddonFile wraps important values for files of an addon
type AddonFile struct {
	Name         string   `json:"fileName"`
	URL          string   `json:"downloadUrl"`
	GameVersions []string `json:"gameVersion"`
}

// FeaturedResponse maps the response from featured addons curseapi
type FeaturedResponse struct {
	Featured []AddonInfo `json:"Featured"`
	Popular  []AddonInfo `json:"Popular"`
	Updated  []AddonInfo `json:"RecentlyUpdated"`
}

func getGameList() (list []GameInfo, err error) {
	client := &http.Client{}
	var requestBuilder strings.Builder
	requestBuilder.WriteString(baseURI)
	requestBuilder.WriteString("game?true")
	req, _ := http.NewRequest("GET", requestBuilder.String(), nil)

	resp, err := client.Do(req)

	if err != nil {
		// do error handling
		fmt.Println("Errored when sending request to the server")
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var array []GameInfo
	_ = json.Unmarshal([]byte(respBody), &array)
	return array, nil
}

func getWoWGameID() (int, error) {
	var games, err = getGameList()
	if err != nil {
		return -1, err
	}
	for _, game := range games {
		if game.Name == "World of Warcraft" {
			return game.GameID, nil
		}
	}
	return -1, errors.New("No GameID found")

}

// GetFeaturedWoWAddons returns list of featured addons from curseapi
func GetFeaturedWoWAddons() (*FeaturedResponse, error) {
	wowID, err := getWoWGameID()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var uriB strings.Builder
	uriB.WriteString(baseURI)
	uriB.WriteString("addon/featured")
	reqBody := []byte(`{
			"GameId": ` + fmt.Sprint(wowID) + `,
			"addonIds": [],
			"featuredCount": 6,
			"popularCount": 14,
			"updatedCount": 14
		}`)

	req, _ := http.NewRequest("POST", uriB.String(), bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	addonList := &FeaturedResponse{}
	jsonErr := json.Unmarshal([]byte(respBody), addonList)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return nil, jsonErr
	}
	return addonList, nil

}

// SearchForWoWAddon polls curseapi for given searchterm and returns list of AddonInfos
func SearchForWoWAddon(searchterm string, gameVersion string) ([]AddonInfo, error) {
	wowID, err := getWoWGameID()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURI+"addon/search", nil)
	q := req.URL.Query()
	q.Add("gameId", fmt.Sprint(wowID))
	q.Add("searchFilter", searchterm)
	if len(gameVersion) > 0 {
		q.Add("gameVersion", gameVersion)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var list []AddonInfo
	_ = json.Unmarshal([]byte(respBody), &list)

	return list, nil
}

// DownloadAddon tries to download the addon files for the given AddonInfo
func DownloadAddon(addonInfo AddonInfo) error {

	return nil
}
