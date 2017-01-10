package mivo

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
	"os"
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
	cacheDir string = "cache"
)

func GetChannels() (out_channels *[]Channel, err error) {
  now := time.Now().Unix()
	if channels != nil && channelsLastUpdate + 600 > now  {
		return &channels, nil
	}
	out_channels = nil
	resp, err := httpGet("https://api.mivo.com/v4/web/channels")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var temp []Channel
	dec := json.NewDecoder(resp.Body)
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
	channelsLastUpdate = now
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
  cacheName := fmt.Sprintf("%s/%d.m3u8", cacheDir, c.ID)
  if reader, size, ok := fetchPlaylistFromFile(cacheName); ok {
    return reader, size, nil
  }
  
	sign, err := GetSign()
	if err != nil {
		return nil, 0, err
	}
	url := c.URL + sign
	resp, err := httpGet(url)
	if err != nil {
	  return nil, 0, err
	}
	
	size := resp.ContentLength
	
	fmt.Printf("cacheName %s\n", cacheName)
	if file, err := os.OpenFile(cacheName, os.O_RDWR | os.O_CREATE, 0644); err == nil {
	  if written, err := io.CopyN(file, resp.Body, size); written == size && err == nil {
	    resp.Body.Close()
	    file.Seek(0, 0)
	    return file, size, nil
	  } else {
	    fmt.Printf(" copy error: %s", err.Error())
	  }
	} else {
	  fmt.Printf("open error: %s", err.Error())
	}
	
	return resp.Body, size, nil
}

func fetchPlaylistFromFile(fileName string) (reader io.ReadCloser, size int64, ok bool) {
  ok = false
  file, err := os.Open(fileName)
  if err != nil {
    return
  }
  defer func() {
    if !ok {
      file.Close()
    }
  }()
  
  stats, err := file.Stat()
  if err != nil {
    return
  }
  
  now := time.Now().Unix()
  if now > stats.ModTime().Unix() + 600 {
    return
  }
  
  reader = file
  size = stats.Size()
  ok = true
  return
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
