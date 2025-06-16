package fs

import (
	"container/heap"
	"os"
	"syscall"
)

/*
	vfd.go defines virtual file descriptor. we can't do open() and close() syscall
		all the time because they are expensive.

*/

type Vfd struct {
	id        uint8
	file_path string
	os_fd     int
	is_open   bool
	use_count uint8
	index     int8
}

type VfdManager struct {
	next_vfd_id uint8
	max_opens   uint8
	open_vfds   VfdPriorityQueue // currently opened fds
	vfd_table   map[uint8]Vfd    // Vfd.id : Vfd
}

func NewVfdMgr() *VfdManager {
	return &VfdManager{
		next_vfd_id: 1,
		max_opens:   5,
		open_vfds:   make(VfdPriorityQueue, 5),
		vfd_table:   make(map[uint8]Vfd),
	}
}

func (vfdmgr *VfdManager) vfd_open(path string) uint8 {
	// check if path is already opened, dunno if there's any efficient implementation
	for vfd_id, vfd := range vfdmgr.vfd_table {

		// if vfd has path
		if vfd.file_path == path && vfd.id == vfd_id {
			// if vfd is opened then dont open just return the same fd
			if vfd.is_open {
				return vfd_id
			} else {
				os_fd, err := syscall.Open(path, os.O_RDWR, 0o660)

				if err != nil {
					panic(err)
				}

				vfd.is_open = true
				vfd.os_fd = os_fd
				vfd.use_count += 1
				return vfd_id
			}
		}
	}

	// before opening, apply lru to close un-needed vfd
	if vfdmgr.open_vfds.Len() > int(vfdmgr.max_opens) {
		vfd_id_to_close := vfdmgr.open_vfds.Pop().(uint8)
		os_fd_to_close := vfdmgr.vfd_table[vfd_id_to_close].os_fd
		syscall.Close(os_fd_to_close)
		heap.Remove(&vfdmgr.open_vfds, int(vfdmgr.vfd_table[vfd_id_to_close].index))
	}

	// dealing with opening a new file
	vfdmgr.next_vfd_id += 1
	vfd_id := vfdmgr.next_vfd_id
	os_fd, err := syscall.Open(path, os.O_RDWR, 0o660)

	if err != nil {
		panic(err)
	}

	vfd := Vfd{
		id:        vfd_id,
		file_path: path,
		os_fd:     os_fd,
		is_open:   true,
		use_count: 1,
		index:     int8(vfdmgr.open_vfds.Len()),
	}

	vfdmgr.vfd_table[vfd_id] = vfd
	vfdmgr.open_vfds.Push(OpenedVfdItem{vfd_id: vfd_id, priority: vfd.use_count, index: int8(vfdmgr.open_vfds.Len())})

	return vfd_id
}

func (vfdmgr *VfdManager) vfd_write(vfd_id uint8, data []byte) {
	// make data to be written is a stream of bytes... for now
}

func (vfdmgr *VfdManager) vfd_read(vfd_id uint8) {

}

func (vfdmgr *VfdManager) vfd_close(vfd_id uint8) {

	os_fd_to_close := vfdmgr.vfd_table[vfd_id].os_fd
	syscall.Close(os_fd_to_close)
	heap.Remove(&vfdmgr.open_vfds, int(vfdmgr.vfd_table[vfd_id].index))
}
