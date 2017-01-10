package main

import (
	"fmt"
	"github.com/abihf/mivo/mivo"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	channels, _ := mivo.GetChannels()
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	fmt.Fprint(w, "#EXTM3U\n")
	for _, c := range *channels {
		fmt.Fprintf(w, "#EXTINF:-1 tvg-id=\"mivo-%d\",%s\n", c.ID, c.Name)
		if strings.HasSuffix(c.URL, ".m3u8") {
			fmt.Fprintf(w, "http://%s/ch/%d.m3u8\n\n", r.Host, c.ID)
		} else {
			fmt.Fprintf(w, "%s\n\n", c.URL)
		}

	}
}

func handleChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id_string := vars["id"]
	id, err := strconv.ParseInt(id_string, 10, 32)
	if err != nil {
		return
	}

	channel, err := mivo.GetChannel(int32(id))
	if err != nil {
		return
	}

	body, size, err := channel.FetchPlaylist()
	if err != nil {
		return
	}
	defer body.Close()

	header := w.Header()
	header.Set("Content-Length", strconv.FormatInt(size, 10))
	header.Set("Content-Type", "application/vnd.apple.mpegurl")
	io.CopyN(w, body, size)
}

func handleEpg(w http.ResponseWriter, r *http.Request) {
	channels, _ := mivo.GetChannels()

	w.Header().Set("Content-Type", "application/xml")

	fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<tv>\n")
	for _, c := range *channels {
		fmt.Fprintf(w, " <channel id=\"mivo-%d\"><display-name>%s</display-name></channel>\n", c.ID, c.Name)
		for _, p := range c.Schedules {
			fmt.Fprintf(w,
				" <programme start=\"%s\" stop=\"%s\" channel=\"mivo-%d\"><title>%s</title></programme>\n",
				p.Start(), p.Finish(), p.ID, p.Name)
		}
	}
	fmt.Fprint(w, "</tv>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/epg.xml", handleEpg)
	r.HandleFunc("/ch/{id:[0-9]+}.m3u8", handleChannel)
	http.Handle("/", r)
	http.ListenAndServe(":12321", nil)
}
