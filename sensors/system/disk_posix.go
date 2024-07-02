//go:build linux || darwin
// +build linux darwin

package system

import (
	"context"
	"time"

	"golang.org/x/sys/unix"
)

func (ds *diskSpaceSensor) Update(ctx context.Context) (err error) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	ds.UpdatedAt = time.Now()
	var stat unix.Statfs_t
	err = unix.Statfs(ds.Path, &stat)
	if err != nil {
		ds.Error = err
		return err
	}
	ds.FreeSpace = float64(stat.Bavail*uint64(stat.Bsize)) / 1024 / 1024
	ds.UsedSpase = float64((stat.Blocks-stat.Bavail)*uint64(stat.Bsize)) / 1024 / 1024
	ds.Ratio = 100 * ds.FreeSpace / ds.UsedSpase
	return nil
}
