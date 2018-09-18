// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"
import store "github.com/mattermost/mattermost-server/store"

// GroupStore is an autogenerated mock type for the GroupStore type
type GroupStore struct {
	mock.Mock
}

// CreateMember provides a mock function with given fields: groupMember
func (_m *GroupStore) CreateMember(groupMember *model.GroupMember) store.StoreChannel {
	ret := _m.Called(groupMember)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.GroupMember) store.StoreChannel); ok {
		r0 = rf(groupMember)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: groupId
func (_m *GroupStore) Delete(groupId string) store.StoreChannel {
	ret := _m.Called(groupId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(groupId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Get provides a mock function with given fields: groupId
func (_m *GroupStore) Get(groupId string) store.StoreChannel {
	ret := _m.Called(groupId)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(groupId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// GetAllPage provides a mock function with given fields: offset, limit
func (_m *GroupStore) GetAllPage(offset int, limit int) store.StoreChannel {
	ret := _m.Called(offset, limit)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(int, int) store.StoreChannel); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Save provides a mock function with given fields: group
func (_m *GroupStore) Save(group *model.Group) store.StoreChannel {
	ret := _m.Called(group)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(*model.Group) store.StoreChannel); ok {
		r0 = rf(group)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}