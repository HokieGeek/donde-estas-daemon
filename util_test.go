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
	var expectedDummy dummyStruct
	expectedDummy.Val1 = createRandomString()
	expectedDummy.Val2 = rand.Int()

	expectedDummyJson, _ := json.Marshal(expectedDummy)
	expectedDummyStr := string(expectedDummyJson)

	var dummy dummyStruct
	if err := readCloserJsonToStruct(stringToReadCloser(expectedDummyStr), &dummy); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	}

	if expectedDummy.Val1 != dummy.Val1 || expectedDummy.Val2 != dummy.Val2 {
		t.Fatalf("Did not receive expected struct %+v but found %+v", expectedDummy, dummy)
	}

	// Incorrect with JSON object
	if err := readCloserJsonToStruct(stringToReadCloser(`{"id":"foo}`), nil); err == nil {
		t.Error("Did not receive expected error on bad JSON unmarshalling")
	}
}
