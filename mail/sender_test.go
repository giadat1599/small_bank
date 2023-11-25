package mail

import (
	"testing"

	"github.com/giadat1599/small_bank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmai(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddr, config.EmailSenderPassword)
	subject := "Test email"
	content := `
		<h1>Hello World</h1>
		<p>This is a test email from <a href="https://google.com">NextTech</a></p>
	`

	to := []string{"truonggiadat15@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
