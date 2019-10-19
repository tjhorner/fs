package shorty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
	"github.com/tjhorner/fs/shorty"
)

func TestShorten(t *testing.T) {
	conf := shorty.Config{"https://example.com"}

	t.Run("Success", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"POST", "https://example.com/api/shorten",
			httpmock.NewStringResponder(200, `{"result":{"suffix":"test","url":"https://example.com"}}`),
		)

		res, err := shorty.Shorten("https://example.com", &conf)
		assert.Nil(t, err)
		assert.Equal(t, res, "https://example.com/test")
	})

	t.Run("Failure", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"POST", "https://example.com/api/shorten",
			httpmock.NewStringResponder(400, `{"error":"invalid url"}`),
		)

		res, err := shorty.Shorten("https://example.com", &conf)
		assert.NotNil(t, err)
		assert.Equal(t, res, "")
	})
}
