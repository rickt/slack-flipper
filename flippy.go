package flippy

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	"google.golang.org/appengine"     // required to use the appengine logging package
	"google.golang.org/appengine/log" // appengine logging package
	"net/http"
	"os"
	"strings"
)

// our lovely table
const table = "┻━┻"

// vars
var (
	env   envVars
	flips map[string]string
)

// no main() since this is an appengine app
func init() {
	// get runtime options from the app.yaml
	env.SlackTokens = os.Getenv("SLACK_TOKEN")
	env.TriggerWord = os.Getenv("SLACK_TRIGGERWORD")
	// get the flipped characters
	flips = getFlipMap()
	// setup the url handler
	http.HandleFunc("/slack", slackhandler)
}

// func that handles the POST to /slack from slack
func slackhandler(w http.ResponseWriter, r *http.Request) {
	var newtable, reversed string
	// create a google appengine context
	ctx := appengine.NewContext(r)
	hook := slackRequest{}
	// get the data from the POST from slack
	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "parsing form error! err=%s", err)
		http.NotFound(w, r)
		return
	}
	defer r.Body.Close()
	// decode slack request
	err = decodeslackrequest(r, &hook)
	if err != nil {
		log.Errorf(ctx, "validateSlackRequest error! err=%v", err)
	}
	// remove the triggerword, if its there
	triggerText := strings.Replace(strings.Trim(hook.Text, " "), env.TriggerWord, "", 1)
	if triggerText != "" {
		reversed = flipText(triggerText)
	}
	// do it!
	newtable = reverseString(table + " " + reversed)
	// build the response
	payload := Payload{
		ResponseType: "in_channel",
		Text:         "   (╯°□°）╯ " + newtable,
	}
	log.Debugf(ctx, "payload=%v", payload)
	// json it up & send it
	js, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// func that decodes the slack request
func decodeslackrequest(r *http.Request, hook *slackRequest) error {
	// create a google appengine context
	ctx := appengine.NewContext(r)
	// create a decoder
	decoder := schema.NewDecoder()
	// decode
	err := decoder.Decode(hook, r.PostForm)
	if err != nil {
		log.Errorf(ctx, "decode error! err=%s", err)
		return err
	}
	// check it
	if !strings.Contains(env.SlackTokens, hook.Token) || hook.TriggerWord != env.TriggerWord {
		log.Errorf(ctx, "invalid token or trigger! env.Slacktokens=%s, hook.Token=%s, hook.TriggerWord=%s, env.TriggerWord=%s", env.SlackTokens, hook.Token, hook.TriggerWord, env.TriggerWord)
		return errors.New("invalid token or trigger")
	}
	return nil
}

//
// utility functions
//

// flips text up side down
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

// the flipped characters
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

// EOF
