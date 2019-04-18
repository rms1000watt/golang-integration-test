package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/rms1000watt/golang-integration-test/person"

	"github.com/stretchr/testify/assert"
)

var (
	serverURL string
	fullURL   string
)

func TestMain(m *testing.M) {
	serverURL = os.Getenv("PERSON_SVC_TEST_SERVER_URL")
	fullURL = serverURL + "/person"

	os.Exit(m.Run())
}

func TestInsertGet(t *testing.T) {
	people := []person.Person{
		{
			Name: "Ryan",
			Age:  99,
		},
		{
			Name: "Bryan",
			Age:  81,
		},
	}

	for _, personIn := range people {
		url := fullURL + fmt.Sprintf("?name=%s&age=%d", personIn.Name, personIn.Age)

		res, err := http.Post(url, "", nil)
		if err != nil {
			t.Error("Error POST person:", err)
			return
		}
		if res.StatusCode != http.StatusOK {
			t.Error("Failed POST person:", res.Status)
			return
		}
		res.Body.Close()

		res, err = http.Get(url)
		if err != nil {
			t.Error("Error GET person:", err)
			return
		}
		if res.StatusCode != http.StatusOK {
			t.Error("Failed GET person:", res.Status)
			return
		}

		var personOut person.Person
		if err := json.NewDecoder(res.Body).Decode(&personOut); err != nil {
			t.Error("Failed json decode:", err)
		}
		res.Body.Close()

		assert.Equal(t, personIn, personOut)
	}
}
