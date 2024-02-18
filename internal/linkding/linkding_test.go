package linkding

import (
	"context"
	"fmt"
	"linkding-telegram/internal/config"
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestGetBookmarks(t *testing.T) {
	// TODO delete strings
	lkdgAddr := "http://192.168.1.88:9090"
	lkdgUsrToken := "e107007a011da7773de83671397dd7d1d67f63bd"

	conf := config.LinkdingConf{
		LinkdingAddr: lkdgAddr,
		UserToken:    lkdgUsrToken,
	}

	lkdg := New(conf)

	resp, err := lkdg.GetBookmarks(context.Background(), "", "1", "")
	assert.NoError(t, err)

	fmt.Println(string(resp))
}
