package main

import "syscall"

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func (d *DiskStatus) diskUsage(path string) error {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return err
	}
	d.All = fs.Blocks * uint64(fs.Bsize)
	d.Free = fs.Bfree * uint64(fs.Bsize)
	d.Used = d.All - d.Free

	return nil
}
