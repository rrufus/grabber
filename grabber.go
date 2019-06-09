package grabber

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Item struct {
	Key string
	Url string
}

// Simultaneously call all endpoints
func Grab(client *http.Client, items ...*Item) (map[string][]byte, error) {
	result := map[string][]byte{}

	// Initialise the right number of error channels
	var errorChannels = make([]chan error, len(items))
	for i := range errorChannels {
		errorChannels[i] = make(chan error)
	}

	// Send off itemes in separate goroutines
	for i, item := range items {
		go item.get(result, errorChannels[i], client)
	}

	// Block on waiting for errors
	for i, item := range items {

		err := <-errorChannels[i]
		if err != nil {
			// log.Println(err)
			return nil, fmt.Errorf("Error from %v [%v]: %v", item.Key, item.Url, err)
		}
	}

	return result, nil

}

func (s *Item) get(r map[string][]byte, errChannel chan error, client *http.Client) {
	req, err := http.NewRequest("GET", s.Url, nil)
	if err != nil {
		errChannel <- err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		errChannel <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChannel <- errors.New(fmt.Sprintf("Error from %v|%v: %v", s.Key, s.Url, resp.Status))
	}

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errChannel <- err
		return
	}

	r[s.Key] = responseText
	close(errChannel)
}
