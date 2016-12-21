package dondeestas

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

func readCloserJSONToStruct(stream io.ReadCloser, data interface{}) error {
	if stream == nil {
		return errors.New("Cannot read from nil stream")
	}

	str, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}

	defer stream.Close()

	if err := json.Unmarshal(str, &data); err != nil {
		return err
	}

	return nil
}
