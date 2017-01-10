package mivo

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Channel struct {
	ID        int32
	Name      string
	Slug      string `json:"dataSlug"`
	URL       string `json:"dataUrl"`
	Image     string
	Schedules []Schedule
}

type Schedule struct {
	ID         int32
	Name       string
	JsonStart  string `json:"start"`
	JsonFinish string `json:"finish"`
	Thumbnail  string `json:"thumbnail_url"`
	Now        bool
}

var (
	channelsLastUpdate int64     = 0
	channels           []Channel = nil
	channelMap         map[int32]*Channel
)

func GetChannels() (out_channels *[]Channel, err error) {
	if channels != nil {
		return &channels, nil
	}
	out_channels = nil
	file, err := os.Open("channels.json")
	if err != nil {
		return
	}
	defer file.Close()

	var temp []Channel
	dec := json.NewDecoder(file)
	if err := dec.Decode(&temp); err != nil {
		return nil, err
	}
	channels = temp
	chMap := make(map[int32]*Channel)
	for i, c := range channels {
		chMap[c.ID] = &temp[i]
	}
	channelMap = chMap

	out_channels = &temp
	return
}

func GetChannel(id int32) (*Channel, error) {
	_, err := GetChannels()
	if err != nil {
		return nil, err
	}
	return channelMap[id], nil
}

func (c *Channel) FetchPlaylist() (io.ReadCloser, int64, error) {
	sign, err := GetSign()
	if err != nil {
		return nil, 0, err
	}
	url := c.URL + sign
	resp, err := httpGet(url)
	return resp.Body, resp.ContentLength, err
}

func formatDateTime(s string) string {
	now := time.Now().Format("20060102")
	return now + s[4:6] + s[7:9] + s[10:12] + " +0700"
}

func (s Schedule) Start() string {
	return formatDateTime(s.JsonStart)
}

func (s Schedule) Finish() string {
	return formatDateTime(s.JsonStart)
}
