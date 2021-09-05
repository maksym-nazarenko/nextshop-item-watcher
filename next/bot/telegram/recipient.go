package telegram

type ChatID string

func (cid ChatID) Recipient() string {
	return string(cid)
}
