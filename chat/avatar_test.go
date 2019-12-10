package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, noAvatarErr := authAvatar.GetAvatarURL(client)

	// TODO Can we assert what error was returned instead instead of this opposite approach?
	// OTOH, it NOT being an ErrNoAvatarURL is the unexpected / exceptional flow in this case (#LOS)
	if noAvatarErr != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	testURL := "http://example.com"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err := authAvatar.GetAvatarURL(client)

	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}

	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "60a6c20d49f49bc210ac98d7e47c74a0"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}

	if url != "//www.gravatar.com/avatar/60a6c20d49f49bc210ac98d7e47c74a0" {
		t.Errorf("GravatarAvatar.GetAvatarURL returned incorrect url or hash %s", url)
	}
}
