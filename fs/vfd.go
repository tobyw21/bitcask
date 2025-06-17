package fs

import (
	"container/list"
	"os"
	"syscall"
)

/*
	vfd.go defines virtual file descriptor. we can't do open() and close() syscall
		all the time because they are expensive.

*/

type Vfd struct {
	id        int8
	file_path string
	os_fd     int
	is_open   bool
}
type VfdManager struct {
	next_vfd_id int8
	max_opens   int
	open_vfds   *list.List             // currently opened fds
	vfd_lru_map map[int8]*list.Element // vfd.id : element(Vfd)
	vfd_table   map[string]Vfd         // Vfd.path : Vfd
}

func NewVfdMgr() *VfdManager {
	return &VfdManager{
		next_vfd_id: 1,
		max_opens:   5,
		open_vfds:   list.New(),                   // the front is recently used, tail is least used
		vfd_lru_map: make(map[int8]*list.Element), // keep a index of vfd.id : open vfd element
		vfd_table:   make(map[string]Vfd),         // vfd.file_path : Vfd
	}

}

func (vfdmgr *VfdManager) VfdOpen(path string) (int8, error) {

	// use path to open vfd
	if vfd, ok := vfdmgr.vfd_table[path]; ok {
		if vfd.is_open {
			e := vfdmgr.vfd_lru_map[vfd.id]
			vfdmgr.open_vfds.PushFront(e)
			return vfd.id, nil
		} else {
			// if close, open it and manipulate with vfd lru map
			// remove current and add to open vfds front
			fd, err := syscall.Open(path, os.O_RDWR|os.O_CREATE, 0o644)
			if err != nil {
				return -1, err
			}

			vfd.is_open = true
			vfd.os_fd = fd
			vfd.file_path = path
			e := vfdmgr.vfd_lru_map[vfd.id]
			vfdmgr.open_vfds.PushFront(e)
			return vfd.id, nil

		}
	}
	// if over max open
	if vfdmgr.open_vfds.Len() >= vfdmgr.max_opens {
		// close least used
		e := vfdmgr.open_vfds.Back()
		vfd := e.Value.(Vfd)
		vfdmgr.open_vfds.Remove(e)
		// fmt.Printf("Closing vfd: %d, with fd: %d\n", vfd.id, vfd.os_fd)
		// clean lru list and vfd table mapping
		delete(vfdmgr.vfd_table, vfd.file_path)
		delete(vfdmgr.vfd_lru_map, vfd.id)
		err := syscall.Close(vfd.os_fd)
		if err != nil {
			return -1, err
		}
	}
	new_vfd_id := vfdmgr.next_vfd_id
	vfdmgr.next_vfd_id += 1
	// deal with unopened new file
	fd, err := syscall.Open(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return -1, err
	}

	vfd := Vfd{
		id:        new_vfd_id,
		file_path: path,
		os_fd:     fd,
		is_open:   true,
	}
	vfdmgr.vfd_table[path] = vfd

	vfd_elem := vfdmgr.open_vfds.PushFront(vfd)
	vfdmgr.vfd_lru_map[new_vfd_id] = vfd_elem

	return vfd.id, nil
}

func (vfdmgr *VfdManager) VfdWrite(vfd_id int8, data []byte) {
	// make data to be written is a stream of bytes... for now
}

func (vfdmgr *VfdManager) VfdRead(vfd_id int8) {

}

func (vfdmgr *VfdManager) VfdClose(vfd_id int8) error {
	if e, ok := vfdmgr.vfd_lru_map[vfd_id]; ok {
		vfd := e.Value.(Vfd)
		// clean lru list and vfd table mapping
		vfdmgr.open_vfds.Remove(e)
		delete(vfdmgr.vfd_table, vfd.file_path)
		delete(vfdmgr.vfd_lru_map, vfd.id)
		err := syscall.Close(vfd.os_fd)
		if err != nil {
			return err
		}

		return nil
	}
	return os.ErrClosed

}
