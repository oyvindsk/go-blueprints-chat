package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL() should return ErrNoAvatarURL for empty client")
	}

	// set a value
	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL() should not return an error when given a avatar_url")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should retutn correct URL")
		}
	}
}
