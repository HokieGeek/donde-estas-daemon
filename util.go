package dondeestas

import (
	"encoding/json"
    "io"
    "io/ioutil"
)

func readCloserJsonToStruct(stream io.ReadCloser, data interface{}) error {
	defer stream.Close();
	
	str, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(str, &data); err != nil {
		return err
	}

	return nil
}
