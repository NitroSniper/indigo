package server

import (
	"bytes"
	_ "embed" // embed is not used directly but for its macro
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/gorilla/websocket"
)

func mdToHtml(source []byte) []byte {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // value chosen at random
	WriteBufferSize: 1024,
}

func (config *ServerConfig) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
	} else {
		var lastMod time.Time
		if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
			lastMod = time.Unix(0, n)
		}
		go config.writer(ws, lastMod)
		reader(ws)
	}

}

const (
	WAITING_PONG = 60 * time.Second
	// Ping need to be less than pong so it can ping pong...
	PING_PERIOD = WAITING_PONG * 9 / 10

	// Poll file for any changes
	FILE_TIMEOUT = 1 * time.Second
)

func readFileIfModified(filename string, lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := os.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}

func (config *ServerConfig) writer(ws *websocket.Conn, lastMod time.Time) {
	lastErr := ""
	pingTicker := time.NewTicker(PING_PERIOD)
	fileTicker := time.NewTicker(FILE_TIMEOUT)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {

		case <-fileTicker.C:
			var md []byte
			var err error
			md, lastMod, err = readFileIfModified(config.Name, lastMod)
			//
			if err != nil {
				if s := err.Error(); s != lastErr {
					lastErr = s
					md = []byte(s)
				}
			} else {
				lastErr = ""
				md = mdToHtml(md)
			}

			if md != nil {
				ws.SetWriteDeadline(time.Now().Add(WAITING_PONG))
				// don't know what writeMessage does need to look into it
				if err := ws.WriteMessage(websocket.TextMessage, md); err != nil {
					return
				}
			}

		case <-pingTicker.C:
			// socket is good so it can still live for another day
			ws.SetWriteDeadline(time.Now().Add(WAITING_PONG))
			// don't know what writeMessage does need to look into it
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512) // random number
	ws.SetWriteDeadline(time.Now().Add(WAITING_PONG))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(WAITING_PONG))
		return nil
	})
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
}

//go:embed assets/template/base.html
var base_template string

//go:embed assets/flavor/pico.css
var flavor string

func (config *ServerConfig) preview(w http.ResponseWriter, r *http.Request) {
	boxing := `
	.markdown-body {
		padding: 64px;
	}
	body {
	  margin: 0;
	}
	`
	templ := template.Must(template.New("index").Parse(base_template))

	md, lastMod, err := readFileIfModified(config.Name, time.Time{})
	if err != nil {
		md = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	} else {
		md = mdToHtml(md)
	}
	data := struct {
		// Markdown template.HTML
		Flavor   template.CSS
		Markdown template.HTML
		Host     string
		LastMod  string
	}{
		// Markdown: template.HTML(markdown.Bytes()),
		Flavor:   template.CSS(flavor + boxing),
		Markdown: template.HTML(md),
		Host:     r.Host,
		LastMod:  strconv.FormatInt(lastMod.UnixNano(), 16),
	}

	if err := templ.Execute(w, data); err != nil {
		panic(err)
	}
}

type ServerConfig struct {
	Name string
}

func (config *ServerConfig) HostServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", config.serveWs)
	mux.HandleFunc("GET /", config.preview)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
