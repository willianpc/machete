package machete_test

import (
	"fmt"
	"testing"

	"github.com/willianpc/machete"
)

func TestAll(t *testing.T) {

	gj, err := machete.NewGenericJson("./template.json")

	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(gj)
	}
}
