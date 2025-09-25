package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

// bad-word-filter

func BadWordReplacement (payload string) string{	
  original := strings.Split(payload, " ")  
	out := make([]string, 0, len(original))
	wordsToFilter := []string{"kerfuffle", "sharbert", "fornax"}

	for _, o := range original {
		if slices.Contains(wordsToFilter, strings.ToLower(o)){
			out = append(out, "****")	
		}	else {
			out = append(out, o)	
		}			
	}

	return strings.Join(out, " ")
}


// respond-with-JSON

func RespondWithJSON(w http.ResponseWriter, code int, payload any) error {
  response, err := json.Marshal(payload)
  if err != nil {
      return err
  }

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

// respond-with-error

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
    return RespondWithJSON(w, code, map[string]string{"error": msg})
}

// respond-no-content
func RespondNoContent(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)
}

// parse-request

func ParseRequest[T any](r *http.Request) (T, error) {
	var params T
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		return params, fmt.Errorf("Something went wrong")
	}
	if err := json.Unmarshal(dat, &params); err != nil {
		return params, fmt.Errorf("Couldn't unmarshal parameters")
	}
	return params, nil
}

