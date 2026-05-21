package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var flipMap = map[rune]rune{
	'a': 'ɐ', 'b': 'q', 'c': 'ɔ', 'd': 'p', 'e': 'ǝ',
	'f': 'ɟ', 'g': 'ɓ', 'h': 'ɥ', 'i': 'ᴉ', 'j': 'ɾ',
	'k': 'ʞ', 'l': 'l', 'm': 'ɯ', 'n': 'u', 'o': 'o',
	'p': 'd', 'q': 'b', 'r': 'ɹ', 's': 's', 't': 'ʇ',
	'u': 'n', 'v': 'ʌ', 'w': 'ʍ', 'x': 'x', 'y': 'ʎ',
	'z': 'z',
	'A': '∀', 'B': 'B', 'C': 'Ɔ', 'D': 'D', 'E': 'Ǝ',
	'F': 'Ⅎ', 'G': 'פ', 'H': 'H', 'I': 'I', 'J': 'ſ',
	'K': 'K', 'L': '˥', 'M': 'W', 'N': 'N', 'O': 'O',
	'P': 'Ԁ', 'Q': 'Q', 'R': 'R', 'S': 'S', 'T': '⊥',
	'U': '∩', 'V': 'Λ', 'W': 'M', 'X': 'X', 'Y': '⅄',
	'Z': 'Z',
	'0': '0', '1': 'Ɩ', '2': 'ᄅ', '3': 'Ɛ', '4': 'ㄣ',
	'5': 'ގ', '6': '9', '7': 'ㄥ', '8': '8', '9': '6',
	'.': '˙', ',': '\'', '?': '¿', '!': '¡', '"': '„',
	'\'': ',', '_': '‾', '(': ')', ')': '(', '[': ']',
	']': '[', '{': '}', '}': '{', '<': '>', '>': '<',
	'&': '⅋', '\n': '\n',
}

func flipText(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	var b strings.Builder
	for _, r := range runes {
		if f, ok := flipMap[r]; ok {
			b.WriteRune(f)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func main() {
	token := flag.String("token", "", "Slack slash command token")
	webhook := flag.String("webhook", "", "Incoming webhook URL")
	norage := flag.Bool("norage", false, "Disable rage emoji")
	port := flag.Int("port", 3000, "Listen port")
	flag.Parse()

	if *token == "" || *webhook == "" {
		fmt.Fprintln(os.Stderr, "token and webhook are required")
		os.Exit(1)
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"wow": "such health"})
	})

	http.HandleFunc("/tableflip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.FormValue("token") != *token {
			http.Error(w, `{"success":false,"error":"Invalid token."}`, http.StatusUnauthorized)
			return
		}
		channel := r.FormValue("channel_name")
		if channel == "directmessage" {
			channel = r.FormValue("channel_id")
		} else {
			channel = "#" + channel
		}
		text := r.FormValue("text")
		if text == "" {
			text = "┻━┻"
		} else {
			text = flipText(text)
		}
		payload := map[string]interface{}{
			"channel": channel,
			"text":    "(╯°□°）╯︵ " + text,
		}
		if !*norage {
			payload["icon_emoji"] = ":rage1:"
		}
		data, _ := json.Marshal(payload)
		resp, err := http.PostForm(*webhook, map[string][]string{"payload": {string(data)}})
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false,"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Listening on :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
