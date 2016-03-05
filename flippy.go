package flippy

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
	"strings"
)

var env envVars
var flips map[string]string

func init() {
	env.SlackTokens = os.Getenv("SLACK_TOKEN")
	env.TriggerWord = os.Getenv("SLACK_TRIGGERWORD")
	flips = getFlipMap()
	http.HandleFunc("/slack", slackPostHandler)
}

func slackPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	hookRequest := slackRequest{}
	table := "┻━┻"
	log.Debugf(ctx, "SLACK_TOKEN=%s, SLACK_TRIGGERWORD=%s", env.SlackTokens, env.TriggerWord)
	var reversed string
	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "error parsing form=%s", err)
		http.NotFound(w, r)
		return
	}
	defer r.Body.Close()
	err = validateSlackRequest(r, &hookRequest)
	if err != nil {
		log.Errorf(ctx, "validateSlackRequest error=%v", err)
	}
	// remove the triggerword, if its there
	triggerText := strings.Replace(strings.Trim(hookRequest.Text, " "), env.TriggerWord, "", 1)
	if triggerText != "" {
		reversed = flipText(triggerText)
	}
	table = reverseString(table + " " + reversed)
	// build the response
	payload := Payload{
		ResponseType: "in_channel",
		Text:         "   (╯°□°）╯ " + table,
	}
	log.Debugf(ctx, "payload=%v", payload)
	// marshal & send it
	js, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func validateSlackRequest(r *http.Request, hookRequest *slackRequest) error {
	ctx := appengine.NewContext(r)
	decoder := schema.NewDecoder()
	err := decoder.Decode(hookRequest, r.PostForm)
	if err != nil {
		return err
	}
	log.Infof(ctx, "validate error=%#v\n", hookRequest)
	if !strings.Contains(env.SlackTokens, hookRequest.Token) || hookRequest.TriggerWord != env.TriggerWord {
		log.Debugf(ctx, "triggerword=%s\n", hookRequest.TriggerWord)
		return errors.New("invalid token or trigger")
	}
	return nil
}

// utility functions
func flipText(input string) string {
	var flipped string
	for _, rune := range input {
		letter := string(rune)
		// get matches
		if flips[letter] != "" {
			flipped += flips[letter]
		} else {
			flipped += letter
		}
	}
	return flipped
}

func reverseString(input string) string {
	// Get Unicode code points.
	n := 0
	rune := make([]rune, len(input))
	for _, r := range input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	// Convert back to UTF-8.
	output := string(rune)
	return output
}

func getFlipMap() map[string]string {
	var flips = make(map[string]string)
	flips["a"] = "ɐ"
	flips["b"] = "q"
	flips["c"] = "ɔ"
	flips["d"] = "p"
	flips["e"] = "ǝ"
	flips["f"] = "ɟ"
	flips["g"] = "ƃ"
	flips["h"] = "ɥ"
	flips["i"] = "ᴉ"
	flips["j"] = "ɾ"
	flips["k"] = "ʞ"
	flips["l"] = "l"
	flips["m"] = "ɯ"
	flips["n"] = "u"
	flips["o"] = "o"
	flips["p"] = "d"
	flips["q"] = "b"
	flips["r"] = "ɹ"
	flips["s"] = "s"
	flips["t"] = "ʇ"
	flips["u"] = "n"
	flips["v"] = "ʌ"
	flips["w"] = "ʍ"
	flips["x"] = "x"
	flips["y"] = "ʎ"
	flips["z"] = "z"
	flips["A"] = "∀"
	flips["B"] = "B"
	flips["C"] = "Ɔ"
	flips["D"] = "D"
	flips["E"] = "Ǝ"
	flips["F"] = "Ⅎ"
	flips["G"] = "פ"
	flips["H"] = "H"
	flips["I"] = "I"
	flips["J"] = "ſ"
	flips["K"] = "K"
	flips["L"] = "˥"
	flips["M"] = "W"
	flips["N"] = "N"
	flips["O"] = "O"
	flips["P"] = "Ԁ"
	flips["Q"] = "Q"
	flips["R"] = "R"
	flips["S"] = "S"
	flips["T"] = "┴"
	flips["U"] = "∩"
	flips["V"] = "Λ"
	flips["W"] = "M"
	flips["X"] = "X"
	flips["Y"] = "⅄"
	flips["Z"] = "z"
	flips["0"] = "0"
	flips["1"] = "Ɩ"
	flips["2"] = "ᄅ"
	flips["3"] = "Ɛ"
	flips["4"] = "ㄣ"
	flips["5"] = "ϛ"
	flips["6"] = "9"
	flips["7"] = "ㄥ"
	flips["8"] = "8"
	flips["9"] = "6"
	flips[","] = "'"
	flips["."] = "˙"
	flips["?"] = "¿"
	flips["!"] = "¡"
	flips["\""] = ",,"
	flips["'"] = ","
	flips["`"] = ","
	flips["("] = ")"
	flips[")"] = "("
	flips["["] = "]"
	flips["]"] = "["
	flips["{"] = "}"
	flips["}"] = "{"
	flips["<"] = ">"
	flips[">"] = "<"
	flips["&"] = "⅋"
	flips["_"] = "‾"
	return flips
}
