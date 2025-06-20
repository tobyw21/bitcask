package fs

import (
	"container/list"
	"errors"
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
	vfd_table   map[int8]Vfd           // Vfd.path : Vfd
	vfd_path_id map[string]int8        // path: vfd.id
}

func NewVfdMgr() *VfdManager {
	return &VfdManager{
		next_vfd_id: 1,
		max_opens:   5,
		open_vfds:   list.New(),                   // the front is recently used, tail is least used
		vfd_lru_map: make(map[int8]*list.Element), // keep a index of vfd.id : open vfd element
		vfd_table:   make(map[int8]Vfd),           // vfd.id : Vfd
		vfd_path_id: make(map[string]int8),        // path: vfd.id
	}

}

func (vfdmgr *VfdManager) VfdOpen(path string) (int8, error) {

	// use path to open vfd
	if vfd_id, ok := vfdmgr.vfd_path_id[path]; ok {
		vfd := vfdmgr.vfd_table[vfd_id]
		if vfd.is_open {
			e := vfdmgr.vfd_lru_map[vfd.id]
			vfdmgr.open_vfds.PushFront(e)
			return vfd.id, nil
		} else {
			// if close, open it and manipulate with vfd lru map
			// remove current and add to open vfds front
			fd, err := syscall.Open(path, syscall.O_RDWR|syscall.O_CREAT, 0o644)
			if err != nil {
				return -1, err
			}

			vfd.is_open = true
			vfd.os_fd = fd
			vfd.file_path = path
			e := vfdmgr.vfd_lru_map[vfd.id]
			vfdmgr.open_vfds.Remove(e)
			vfdmgr.open_vfds.PushFront(e)
			return vfd.id, nil

		}
	}
	// if over max open
	if vfdmgr.open_vfds.Len() >= vfdmgr.max_opens {
		// close least used
		e := vfdmgr.open_vfds.Back()
		vfd := e.Value.(Vfd)
		// fmt.Printf("Closing vfd: %d, with fd: %d\n", vfd.id, vfd.os_fd)
		// clean lru list
		// can't clean vfd table, what if file needs to be opened again next time?
		vfdmgr.open_vfds.Remove(e)
		delete(vfdmgr.vfd_lru_map, vfd.id)
		err := syscall.Close(vfd.os_fd)
		if err != nil {
			return -1, err
		}
	}
	new_vfd_id := vfdmgr.next_vfd_id
	vfdmgr.next_vfd_id += 1
	// deal with unopened new file
	fd, err := syscall.Open(path, syscall.O_RDWR|syscall.O_CREAT, 0o644)
	if err != nil {
		return -1, err
	}

	vfd := Vfd{
		id:        new_vfd_id,
		file_path: path,
		os_fd:     fd,
		is_open:   true,
	}
	vfdmgr.vfd_path_id[path] = vfd.id
	vfdmgr.vfd_table[vfd.id] = vfd

	vfd_elem := vfdmgr.open_vfds.PushFront(vfd)
	vfdmgr.vfd_lru_map[new_vfd_id] = vfd_elem

	return vfd.id, nil
}

func (vfdmgr *VfdManager) VfdWrite(vfd_id int8, data []byte, offset int64) (int, error) {
	// make data to be written is a stream of bytes... for now
	if e, ok := vfdmgr.vfd_lru_map[vfd_id]; ok {
		vfd := e.Value.(Vfd)
		vfdmgr.open_vfds.Remove(e)
		vfdmgr.open_vfds.PushFront(e)
		nw, err := syscall.Pwrite(vfd.os_fd, data, offset)

		if err != nil {
			return -1, err
		}

		return nw, nil
	} else {
		// if not in lru map means not opened
		// open it from vfd table
		if vfd, ok := vfdmgr.vfd_table[vfd_id]; ok {

			vfd_id, err := vfdmgr.VfdOpen(vfd.file_path)
			if err != nil {
				return -1, err
			}

			vfd := vfdmgr.vfd_table[vfd_id]
			// manage vfd lru mapping
			vfd_elem := vfdmgr.open_vfds.PushFront(vfd)
			vfdmgr.vfd_lru_map[vfd_id] = vfd_elem
			vfdmgr.vfd_path_id[vfd.file_path] = vfd.id

			nw, err := syscall.Pwrite(vfd.os_fd, data, offset)

			if err != nil {
				return -1, err
			}

			return nw, nil

		}
	}

	return -1, errors.New("unable to find Vfd Id corresponding file")
}

func (vfdmgr *VfdManager) VfdRead(vfd_id int8, buffer *[]byte, offset int64) (int, error) {
	// can't assume it is always being opened
	// if not in lru map, then find it in vfd table
	// repoen it, change lru and other metadata
	// else the vfd_id is invalid
	// find correct fd to read
	if e, ok := vfdmgr.vfd_lru_map[vfd_id]; ok {
		vfd := e.Value.(Vfd)
		vfdmgr.open_vfds.Remove(e)
		vfdmgr.open_vfds.PushFront(e)
		nread, err := syscall.Pread(vfd.os_fd, *buffer, offset)

		if err != nil {
			return -1, err
		}

		return nread, nil
	} else {
		// if not in lru map means not opened
		// open it from vfd table
		if vfd, ok := vfdmgr.vfd_table[vfd_id]; ok {

			vfd_id, err := vfdmgr.VfdOpen(vfd.file_path)
			if err != nil {
				return -1, err
			}

			vfd := vfdmgr.vfd_table[vfd_id]
			// manage vfd lru mapping
			vfd_elem := vfdmgr.open_vfds.PushFront(vfd)
			vfdmgr.vfd_lru_map[vfd_id] = vfd_elem
			vfdmgr.vfd_path_id[vfd.file_path] = vfd.id
			nread, err := syscall.Pread(vfd.os_fd, *buffer, offset)

			if err != nil {
				return -1, err
			}

			return nread, nil

		}
	}

	return -1, errors.New("unable to find Vfd Id corresponding file")
}

func (vfdmgr *VfdManager) VfdClose(vfd_id int8) error {
	if e, ok := vfdmgr.vfd_lru_map[vfd_id]; ok {
		vfd := e.Value.(Vfd)
		// clean lru list and vfd table mapping
		vfdmgr.open_vfds.Remove(e)
		// delete(vfdmgr.vfd_table, vfd.file_path)
		delete(vfdmgr.vfd_lru_map, vfd.id)
		err := syscall.Close(vfd.os_fd)
		if err != nil {
			return err
		}

		return nil
	}
	return os.ErrClosed

}

func (vfdmgr *VfdManager) VfdClean() error {
	// this functioni is a placeholder
	// leave it here maybe need to clean vfd manager when bitstore is closing
	// need to handle metadata discard or write
	return nil

}
