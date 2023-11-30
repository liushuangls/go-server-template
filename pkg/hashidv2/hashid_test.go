package hashidv2

import (
	"testing"
)

func TestHashID_Encode(t *testing.T) {
	var (
		userHash = New(&Config{
			Prefix: "U0",
			Salt:   "user__",
			Type:   0,
		})
		gptDraft = New(&Config{
			Type:   2,
			Prefix: "GD0",
			Salt:   "gpt_draft__",
		})
	)

	var (
		originUserID  = 3
		originDraftID = 7282
	)

	t.Logf("userID: %d encode: %s", originUserID, userHash.EncodeNotE(originUserID))
	t.Logf("gptDraftID: %d encode: %s", originDraftID, gptDraft.EncodeNotE(originDraftID))

	decode, err := userHash.Decode("U07DY3H3B5")
	t.Log(err, decode)
}
