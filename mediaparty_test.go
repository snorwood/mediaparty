package mediaparty

import (
	"testing"
)

func TestStringToPostgresString(t *testing.T) {
	testStrings := map[string]string{
		"'":            "\\'",
		"hello\\bye":   "hello\\\\bye",
		"hello\\\\bye": "hello\\\\bye",
		"\\hello":      "\\\\hello",
		"hello\\":      "hello\\\\",
		"Steven's":     "Steven\\'s",
		"Steven\\'s":   "Steven\\'s",
	}

	for test, expected := range testStrings {
		actual := StringToPostgresString(test)

		if actual != expected {
			t.Errorf("\nTest: %s => Actual: %s != Expected: %s", test, actual, expected)
		}
	}
}

func TestExecuteTemplate(t *testing.T) {
	var err error = nil

	data := struct {
		Schema string
		Table  string
	}{
		"musicplayer",
		"music",
	}
	_, err = ExecuteTemplate(FromTableTemplate, data)
	if err != nil {
		t.Errorf("Error parsing template: %s", err)
	}

	data2 := struct {
		Columns []string
	}{
		[]string{"alpha", "beta", "gamma"},
	}
	result := "dogmakarm"
	result, err = ExecuteTemplate(SelectColumnsTemplate, data2)
	t.Log(result)
	if err != nil {
		t.Errorf("Error parsing template: %s", err)
	}
}

func TestGetSongQuery(t *testing.T) {
	query, err := GetSongQuery("musicplayer", "music", "The National", "Fake Empire", "The Boxer", "The National")
	t.Log(query)

	if err != nil {
		t.Errorf(err.Error())
	}
}
