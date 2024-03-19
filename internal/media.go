package internal

import (
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
	"rerere.org/fedi-games/config"
)

// StoreMedia stores a file and makes it publically available via the media route
// returns the full url of the media file
func StoreMedia(data []byte, ext string) (string, error) {
	cfg := config.GetConfig()
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	name := id.String() + "." + ext
	fPath := path.Join(cfg.MediaPath, name)
	os.MkdirAll(cfg.MediaPath, 0777)
	err = os.WriteFile(fPath, data, 0777)
	if err != nil {
		return "", err
	}
	return url.JoinPath(cfg.FullUrl(), "/media/"+name)
}
