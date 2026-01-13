package util

import (
	"context"
	"os/user"

	"github.com/pkg/errors"
)

type User struct {
	Id            string   `json:"id,omitempty"`
	Name          string   `json:"name"`
	Username      string   `json:"username"`
	HomeDirectory string   `json:"home_directory"`
	GroupIds      []string `json:"group_ids,omitempty"`
}

func (o User) GetArtifactType() ArtifactType {
	return ArtifactTypeUser
}

type UserGroup struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (o UserGroup) GetArtifactType() ArtifactType {
	return ArtifactTypeUserGroup
}

func GetCurrentUser() (*User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	gids, _ := u.GroupIds()
	return &User{
		Username:      u.Username,
		Name:          u.Name,
		HomeDirectory: u.HomeDir,
		GroupIds:      gids,
	}, nil
}

func GetUserByUsername(username string) (*User, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	usr := parseUser(u)
	return &usr, nil
}

func ListUsers(ctx context.Context) ([]User, error) {
	usernames, err := listUsernames(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list usernames")
	}
	users := make([]User, 0, len(usernames))
	for _, username := range usernames {
		user, err := GetUserByUsername(username)
		if err != nil {
			continue
		}
		users = append(users, *user)
	}
	return users, nil
}

func parseUser(u *user.User) User {
	gids, _ := u.GroupIds()
	return User{
		Username:      u.Username,
		Name:          u.Name,
		HomeDirectory: u.HomeDir,
		GroupIds:      gids,
	}
}

func ListUserGroups(ctx context.Context) ([]UserGroup, error) {
	groupNames, err := listUserGroupNames(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list user group names")
	}
	groups := make([]UserGroup, 0, len(groupNames))
	for _, g := range groupNames {
		group, err := user.LookupGroup(g)
		if err != nil {
			continue
		}
		groups = append(groups, parseUserGroup(group))
	}
	return groups, nil
}

func parseUserGroup(g *user.Group) UserGroup {
	return UserGroup{
		Id:   g.Gid,
		Name: g.Name,
	}
}
