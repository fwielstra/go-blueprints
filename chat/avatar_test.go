package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	err := ioutil.WriteFile(filename, []byte{}, 0777)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(filename)

	var fsAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "abc"}
	url, err := fsAvatar.GetAvatarURL(client)

	if err != nil {
		t.Error(err.Error())
	}

	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
