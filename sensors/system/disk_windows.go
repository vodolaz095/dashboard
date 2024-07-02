//go:build windows
// +build windows

package system

import (
	"context"
	"time"

	"golang.org/x/sys/windows"
)

func (ds *diskSpaceSensor) Update(ctx context.Context) (err error) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	ds.UpdatedAt = time.Now()
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64

	err = windows.GetDiskFreeSpaceEx(
		windows.StringToUTF16Ptr(ds.Path),
		&freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		ds.Error = err
		return err
	}
	ds.FreeSpace = float64(freeBytesAvailable / 1024 / 1024)
	ds.UsedSpase = float64((totalNumberOfBytes - totalNumberOfFreeBytes) / 1024 / 1024)
	ds.Ratio = ds.UsedSpase / ds.FreeSpace
}
