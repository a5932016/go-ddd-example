package fs

import (
	"context"

	"github.com/a5932016/go-ddd-example/config"
	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func NewFSRepository(fs afero.Fs) FSRepository {
	return FSRepository{
		fs: fs,
	}
}

type FSRepository struct {
	fs afero.Fs
}

func isProduction() bool {
	return config.Env.Core.Mode == gin.ReleaseMode
}

func (fsr FSRepository) MkdirAlbum(path string) error {
	if err := fsr.fs.MkdirAll(path, 0755); err != nil { // '-rwxr-xr-x'
		return errors.Wrap(err, "fs.MkdirAll")
	}

	if isProduction() {
		if err := fsr.fs.Chown(path, 101, 101); err != nil { // 101: Typical Nginx user and group ID
			return errors.Wrap(err, "fs.Chown")
		}
	}

	return nil
}

func (fsr FSRepository) SetImagePermission(path string) error {
	if isProduction() {
		if err := fsr.fs.Chown(path, 101, 101); err != nil { // 101: Typical Nginx user and group ID
			return errors.Wrap(err, "fs.Chown")
		}
	}

	if err := fsr.fs.Chmod(path, 0700); err != nil { // '-rwx------' Only for Owner(Nginx). Frontend will only access images through Nginx
		return errors.Wrap(err, "fs.Chmod")
	}

	return nil
}

func (fsr FSRepository) Remove(c context.Context, path string) {
	if err := fsr.fs.Remove(path); err != nil {
		log.FromContext(c).WithFields(logrus.Fields{"path": path}).WithError(err).Error("fs.Remove")
	}
}

func (fsr FSRepository) RemoveAll(c context.Context, path string) {
	if err := fsr.fs.RemoveAll(path); err != nil {
		log.FromContext(c).WithFields(logrus.Fields{"path": path}).WithError(err).Error("fs.RemoveAll")
	}
}
