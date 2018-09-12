// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

func TestGroupStore(t *testing.T, ss store.Store) {
	t.Run("Save", func(t *testing.T) { testGroupStoreSave(t, ss) })
	t.Run("Get", func(t *testing.T) { testGroupStoreGet(t, ss) })
	t.Run("GetAllPage", func(t *testing.T) { testGroupStoreGetAllPage(t, ss) })
	t.Run("Delete", func(t *testing.T) { testGroupStoreDelete(t, ss) })
	t.Run("CreateMember", func(t *testing.T) { testGroupCreateMember(t, ss) })
}

func testGroupStoreSave(t *testing.T, ss store.Store) {
	// Save a new group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
		Description: model.NewId(),
		TypeProps:   model.NewId(),
	}

	// Happy path
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)
	assert.Equal(t, g1.Name, d1.Name)
	assert.Equal(t, g1.DisplayName, d1.DisplayName)
	assert.Equal(t, g1.Description, d1.Description)
	assert.Equal(t, g1.TypeProps, d1.TypeProps)

	// Requires name and display name
	g2 := &model.Group{
		Name:        "",
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res2 := <-ss.Group().Save(g2)
	assert.Nil(t, res2.Data)
	assert.NotNil(t, res2.Err)
	assert.Equal(t, res2.Err.Id, "model.group.name.app_error")

	g2.Name = model.NewId()
	g2.DisplayName = ""
	res3 := <-ss.Group().Save(g2)
	assert.Nil(t, res3.Data)
	assert.NotNil(t, res3.Err)
	assert.Equal(t, res3.Err.Id, "model.group.display_name.app_error")

	// Can't invent an ID and save it
	g3 := &model.Group{
		Id:          model.NewId(),
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
		CreateAt:    1,
		UpdateAt:    1,
	}
	res4 := <-ss.Group().Save(g3)
	assert.Nil(t, res4.Data)
	assert.Equal(t, res4.Err.Id, "store.sql_group.save.missing.app_error")

	// Won't accept a duplicate name
	g4 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res5 := <-ss.Group().Save(g4)
	assert.Nil(t, res5.Err)
	g4b := &model.Group{
		Name:        g4.Name,
		DisplayName: model.NewId(),
		Type:        model.GroupTypeLdap,
	}
	res5b := <-ss.Group().Save(g4b)
	assert.Nil(t, res5b.Data)
	assert.Equal(t, res5b.Err.Id, "store.sql_group.save.insert.app_error")

	// Fields cannot be greater than max values
	g5 := &model.Group{
		Name:        strings.Repeat("x", model.GroupNameMaxLength),
		DisplayName: strings.Repeat("x", model.GroupDisplayNameMaxLength),
		Description: strings.Repeat("x", model.GroupDescriptionMaxLength),
		TypeProps:   strings.Repeat("x", model.GroupTypePropsMaxLength),
		Type:        model.GroupTypeLdap,
	}
	assert.Nil(t, g5.IsValidForCreate())

	g5.Name = g5.Name + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.name.app_error")
	g5.Name = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	g5.DisplayName = g5.DisplayName + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.display_name.app_error")
	g5.DisplayName = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	g5.Description = g5.Description + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.description.app_error")
	g5.Description = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	g5.TypeProps = g5.TypeProps + "x"
	assert.Equal(t, g5.IsValidForCreate().Id, "model.group.type_props.app_error")
	g5.TypeProps = model.NewId()
	assert.Nil(t, g5.IsValidForCreate())

	// Must use a valid type
	g6 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		TypeProps:   model.NewId(),
		Type:        "fake",
	}
	assert.Equal(t, g6.IsValidForCreate().Id, "model.group.type.app_error")
}

func testGroupStoreGet(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
		TypeProps:   model.NewId(),
	}
	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)

	// Get the group
	res2 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res2.Err)
	d2 := res2.Data.(*model.Group)
	assert.Equal(t, d1.Id, d2.Id)
	assert.Equal(t, d1.Name, d2.Name)
	assert.Equal(t, d1.DisplayName, d2.DisplayName)
	assert.Equal(t, d1.Description, d2.Description)
	assert.Equal(t, d1.TypeProps, d2.TypeProps)
	assert.Equal(t, d1.CreateAt, d2.CreateAt)
	assert.Equal(t, d1.UpdateAt, d2.UpdateAt)
	assert.Equal(t, d1.DeleteAt, d2.DeleteAt)

	// Get an invalid group
	res3 := <-ss.Group().Get(model.NewId())
	assert.NotNil(t, res3.Err)
	assert.Equal(t, res3.Err.Id, "store.sql_group.get.app_error")
}

func testGroupStoreGetAllPage(t *testing.T, ss store.Store) {
	numGroups := 10

	groups := []*model.Group{}

	// Create groups
	for i := 0; i < numGroups; i++ {
		g := &model.Group{
			Name:        model.NewId(),
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Type:        model.GroupTypeLdap,
			TypeProps:   model.NewId(),
		}
		groups = append(groups, g)
		res := <-ss.Group().Save(g)
		assert.Nil(t, res.Err)
	}

	// Returns all the groups
	res1 := <-ss.Group().GetAllPage(0, 999)
	d1 := res1.Data.([]*model.Group)
	assert.Condition(t, func() bool { return len(d1) >= numGroups })
	for _, expectedGroup := range groups {
		present := false
		for _, dbGroup := range d1 {
			if dbGroup.Id == expectedGroup.Id {
				present = true
				break
			}
		}
		assert.True(t, present)
	}

	// Returns the correct number based on limit
	res2 := <-ss.Group().GetAllPage(0, 2)
	d2 := res2.Data.([]*model.Group)
	assert.Len(t, d2, 2)

	// Check that result sets are different using an offset
	res3 := <-ss.Group().GetAllPage(0, 5)
	d3 := res3.Data.([]*model.Group)
	res4 := <-ss.Group().GetAllPage(5, 5)
	d4 := res4.Data.([]*model.Group)
	for _, d3i := range d3 {
		for _, d4i := range d4 {
			if d4i.Id == d3i.Id {
				t.Error("Expected results to be unique.")
			}
		}
	}
}

func testGroupStoreDelete(t *testing.T, ss store.Store) {
	// Save a group
	g1 := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Type:        model.GroupTypeLdap,
		TypeProps:   model.NewId(),
	}

	res1 := <-ss.Group().Save(g1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.Group)
	assert.Len(t, d1.Id, 26)

	// Check the group is retrievable
	res2 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res2.Err)

	// Get the before count
	res7 := <-ss.Group().GetAllPage(0, 999)
	d7 := res7.Data.([]*model.Group)
	beforeCount := len(d7)

	// Delete the group
	res3 := <-ss.Group().Delete(d1.Id)
	assert.Nil(t, res3.Err)

	// Check the group is deleted
	res4 := <-ss.Group().Get(d1.Id)
	assert.Nil(t, res4.Err)
	d2 := res4.Data.(*model.Group)
	assert.NotZero(t, d2.DeleteAt)

	// Check the after count
	res5 := <-ss.Group().GetAllPage(0, 999)
	d5 := res5.Data.([]*model.Group)
	afterCount := len(d5)
	assert.Condition(t, func() bool { return beforeCount == afterCount+1 })

	// Try and delete a nonexistent group
	res6 := <-ss.Group().Delete(model.NewId())
	assert.NotNil(t, res6.Err)
	assert.Equal(t, res6.Err.Id, "store.sql_group.get.app_error")
}

func testGroupCreateMember(t *testing.T, ss store.Store) {
	gm1 := &model.GroupMember{
		GroupId: model.NewId(),
		UserId:  model.NewId(),
	}

	// Happy path
	res1 := <-ss.Group().CreateMember(gm1)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(*model.GroupMember)
	assert.Equal(t, d1.GroupId, gm1.GroupId)
	assert.Equal(t, d1.UserId, gm1.UserId)
	assert.NotNil(t, d1.CreateAt)
	assert.Equal(t, d1.DeleteAt, int64(0))

	// Duplicate composite key (GroupId, UserId)
	res2 := <-ss.Group().CreateMember(gm1)
	assert.NotNil(t, res2.Err)
	assert.Equal(t, res2.Err.Id, "store.sql_group.save_member.exists.app_error")
}
