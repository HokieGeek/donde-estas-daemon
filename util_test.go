package dondeestas

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
)

type dummyStruct struct {
	Val1 string `json:"val1"`
	Val2 int    `json:"val2"`
}

func stringToReadCloser(str string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(bytes.NewBufferString(str).Bytes()))
}

func TestReadCloserJsonToStruct(t *testing.T) {
	expectedDummy := &dummyStruct{Val1: createRandomString(), Val2: rand.Int()}
	expectedDummyJSON, _ := json.Marshal(expectedDummy)
	expectedDummyStr := string(expectedDummyJSON)

	var dummy dummyStruct
	if err := readCloserJSONToStruct(stringToReadCloser(expectedDummyStr), &dummy); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	}
	if expectedDummy.Val1 != dummy.Val1 || expectedDummy.Val2 != dummy.Val2 {
		t.Fatalf("Did not receive expected struct %+v but found %+v", expectedDummy, dummy)
	}

	// Incorrect with JSON object
	if err := readCloserJSONToStruct(stringToReadCloser(`{"id":"foo}`), nil); err == nil {
		t.Error("Did not receive expected error on bad JSON unmarshalling")
	}

	// I am finding it impossible to instigate an error from the ReadAll in the function
	strm := ioutil.NopCloser(bytes.NewReader(make([]byte, 0)))
	strm.Close()
	if err := readCloserJSONToStruct(strm, nil); err == nil {
		t.Error("Did not receive expected error on attempting to read from closed stream")
	} else {
		t.Logf("ERROR: %s", err)
	}

	// Set the stream to nil
	if err := readCloserJSONToStruct(nil, nil); err == nil {
		t.Error("Did not receive expected error on attempting to read from nil stream")
	}
}
