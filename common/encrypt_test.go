package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchPass(t *testing.T) {
	result := MatchPass("516504610", "5965ffb36f1fd2e98529b8e3c2af9a1f80e0a0be0498fca9cd8a92e592cc8af0")
	assert.Equal(t, true, result)
	result = MatchPass("516504610", "5965ffb36f1fd2e98529b8e3c2af9a1f80e0a0be0498fca9cd8a92e592cc8af1")
	assert.Equal(t, false, result)
	result = MatchPass("no", "5965ffb36f1fd2e98529b8e3c2af9a1f80e0a0be0498fca9cd8a92e592cc8af0")
	assert.Equal(t, false, result)
}