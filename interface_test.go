package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFOut(t *testing.T) {
	s := newGameSessionN(10)

	msg := fmt.Sprintf(s.fOut("Traders report that %s harvested %d bushels per acre.\n", "199", 5), "uruk")
	assert.Equal(t, "", msg)
}
