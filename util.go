package dondeestas

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

func readCloserJsonToStruct(stream io.ReadCloser, data interface{}) error {
	if stream == nil {
		return errors.New("Cannot read from nil steam")
	}

	defer stream.Close()

	str, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(str, &data); err != nil {
		return err
	}

	return nil
}
