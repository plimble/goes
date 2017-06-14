package mocks

import goes "github.com/plimble/goes"
import mock "github.com/stretchr/testify/mock"

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// GetFromRevision provides a mock function with given fields: id, from
func (_m *Storage) GetFromRevision(id string, from int) ([]goes.Event, error) {
	ret := _m.Called(id, from)

	var r0 []goes.Event
	if rf, ok := ret.Get(0).(func(string, int) []goes.Event); ok {
		r0 = rf(id, from)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]goes.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(id, from)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastEvent provides a mock function with given fields: id
func (_m *Storage) GetLastEvent(id string) ([]goes.Event, error) {
	ret := _m.Called(id)

	var r0 []goes.Event
	if rf, ok := ret.Get(0).(func(string) []goes.Event); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]goes.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSnapshot provides a mock function with given fields: id, version
func (_m *Storage) GetSnapshot(id string, version int) (*goes.Snapshot, error) {
	ret := _m.Called(id, version)

	var r0 *goes.Snapshot
	if rf, ok := ret.Get(0).(func(string, int) *goes.Snapshot); ok {
		r0 = rf(id, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*goes.Snapshot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(id, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUndispatchedEvent provides a mock function with given fields:
func (_m *Storage) GetUndispatchedEvent() ([]goes.Event, error) {
	ret := _m.Called()

	var r0 []goes.Event
	if rf, ok := ret.Get(0).(func() []goes.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]goes.Event)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkDispatchedEvent provides a mock function with given fields: es
func (_m *Storage) MarkDispatchedEvent(es []goes.Event) error {
	ret := _m.Called(es)

	var r0 error
	if rf, ok := ret.Get(0).(func([]goes.Event) error); ok {
		r0 = rf(es)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: es
func (_m *Storage) Save(es []goes.Event) error {
	ret := _m.Called(es)

	var r0 error
	if rf, ok := ret.Get(0).(func([]goes.Event) error); ok {
		r0 = rf(es)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveSnapshot provides a mock function with given fields: snap
func (_m *Storage) SaveSnapshot(snap *goes.Snapshot) error {
	ret := _m.Called(snap)

	var r0 error
	if rf, ok := ret.Get(0).(func(*goes.Snapshot) error); ok {
		r0 = rf(snap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}