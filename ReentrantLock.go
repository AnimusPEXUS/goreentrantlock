package goreentrantlock

import (
	"github.com/AnimusPEXUS/golockercheckable"
	"github.com/AnimusPEXUS/goroutineid"
)

var _ golockercheckable.LockerCheckable = &ReentrantMutexCheckable{}

// just like github.com/AnimusPEXUS/golockercheckable.MutexCheckable:
// Unlock() on unlocked Mutex - doesn't lead to error
type ReentrantMutexCheckable struct {
	main_mtx       *golockercheckable.MutexCheckable
	local_mtx      *golockercheckable.MutexCheckable
	passed_id      uint64
	passed_counter int64
}

func NewReentrantMutexCheckable(locked bool) *ReentrantMutexCheckable {
	self := new(ReentrantMutexCheckable)
	self.main_mtx = golockercheckable.NewMutexCheckable(false)
	self.local_mtx = golockercheckable.NewMutexCheckable(false)
	if locked {
		self.main_mtx.Lock()
	}
	return self
}

func (self *ReentrantMutexCheckable) IsLocked() bool {
	return self.main_mtx.IsLocked()
}

func (self *ReentrantMutexCheckable) Lock() {
	self.local_mtx.Lock()
	defer self.local_mtx.Unlock()

	id, err := goroutineid.GetCurrentGoId_byRuntimeStack()
	if err != nil {
		panic("can't get goroutineid")
	}

	if !self.main_mtx.IsLocked() {
		self.passed_id = id
		self.main_mtx.Lock()
		self.passed_counter++
	} else {
		if self.passed_id == id {
			self.passed_counter++
		} else {
			self.local_mtx.Unlock()
			self.main_mtx.Lock()
		}
	}

	return
}

func (self *ReentrantMutexCheckable) Unlock() {
	self.local_mtx.Lock()
	defer self.local_mtx.Unlock()

	id, err := goroutineid.GetCurrentGoId_byRuntimeStack()
	if err != nil {
		panic("can't get goroutineid")
	}

	if !self.main_mtx.IsLocked() {
		return
	}

	if self.passed_id == id {

		if self.passed_counter < 0 {
			panic("programming error: this should not be happening")
		}

		if self.passed_counter == 0 {
			return
		}

		if self.passed_counter == 1 {
			self.passed_counter--
			self.main_mtx.Unlock()
		} else {
			self.passed_counter--
		}

	} else {
		panic("programming error: this should not be happening: you should not try unlock what is locked not by you")
	}

}

func (self *ReentrantMutexCheckable) IsLocakedByMe() (locked bool, byme bool) {
	self.local_mtx.Lock()
	defer self.local_mtx.Unlock()

	id, err := goroutineid.GetCurrentGoId_byRuntimeStack()
	if err != nil {
		panic("can't get goroutine id")
	}

	if self.passed_counter > 0 {
		return true, id == self.passed_id
	} else {
		return false, false
	}
}

func (self *ReentrantMutexCheckable) LocekdByWho() (locked bool, goid uint64) {
	self.local_mtx.Lock()
	defer self.local_mtx.Unlock()

	if self.passed_counter > 0 {
		return true, self.passed_id
	} else {
		return false, 0
	}
}
