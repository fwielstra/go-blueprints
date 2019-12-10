package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

//ErrNoAvatarURL is a predefined error returned when any issue occurs determining a user's avatar.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

//Avatar is an interface for fetching an avatar or avatar URL for a given client.
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

//AuthAvatar is an Avatar provider that uses the `avatar_url` value from an oauth response.
type AuthAvatar struct {
}

//UseAuthAvatar is a utility to fetch a client's avatar using the oauth url strategy.
var UseAuthAvatar AuthAvatar //nolint

//GetAvatarURL will try to pull the `avatar_url` field from the client's user data, or return
//`ErrNoAvatarURL` if the user does not have this data.
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}

	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil
}

// GravatarAvatar is an implementation of `Avatar` that will attempt to retrieve an avatar from Gravatar.com,
// based on a hashed versino of the user's email.
type GravatarAvatar struct{}

//UseGravatar is an utility to fetch a client's avatar using the gravatar strategy.
var UseGravatar GravatarAvatar //nolint:unused

//GetAvatarURL will get the user's user ID and create a Gravatar URL.
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	userIDStr, ok := userid.(string)

	if !ok {
		return "", ErrNoAvatarURL
	}

	return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIDStr), nil
}

//FileSystemAvatar is an impelentation of `Avatar` that returns avatar urls (named by the user ID) from the local
// webserver.
type FileSystemAvatar struct{}

//UseFileSystemAvatar is an utility to fetch a client's avatar url using the filesystem strategy.
var UseFileSystemAvatar FileSystemAvatar

//GetAvatarURL attempts to find a matching file in the avatar folder for the user's ID.
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	userIDStr, ok := userid.(string)

	if !ok {
		return "", ErrNoAvatarURL
	}

	return "/avatars/" + userIDStr + ".jpg", nil
}
