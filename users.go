// Copyright 2025- The sacloud/apigw-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apigw

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
)

type UserAPI interface {
	List(ctx context.Context) ([]v1.User, error)
	Create(ctx context.Context, request *v1.UserDetail) (*v1.UserDetail, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.UserDetail, error)
	Update(ctx context.Context, request *v1.UserDetail, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ UserAPI = (*userOp)(nil)

type userOp struct {
	client *v1.Client
}

func NewUserOp(client *v1.Client) UserAPI {
	return &userOp{client: client}
}

func (op *userOp) List(ctx context.Context) ([]v1.User, error) {
	res, err := op.client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetUsersOK:
		return p.Apigw.Users, nil
	case *v1.GetUsersBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUsersInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *userOp) Create(ctx context.Context, request *v1.UserDetail) (*v1.UserDetail, error) {
	res, err := op.client.AddUser(ctx, request)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.AddUserCreated:
		return &p.Apigw.User.Value, nil
	case *v1.AddUserBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.AddUserConflict:
		return nil, errors.New(p.Message.Value)
	case *v1.AddUserInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *userOp) Read(ctx context.Context, id uuid.UUID) (*v1.UserDetail, error) {
	res, err := op.client.GetUser(ctx, v1.GetUserParams{UserId: id})
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetUserOK:
		return &p.Apigw.User.Value, nil
	case *v1.GetUserBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUserNotFound:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUserInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *userOp) Update(ctx context.Context, request *v1.UserDetail, id uuid.UUID) error {
	res, err := op.client.UpdateUser(ctx, request, v1.UpdateUserParams{UserId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateUserNoContent:
		return nil
	case *v1.UpdateUserBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateUserNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateUserInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}

func (op *userOp) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := op.client.DeleteUser(ctx, v1.DeleteUserParams{UserId: id})
	if err != nil {
		return err
	}

	return nil
}

type UserExtraAPI interface {
	ListGroup(ctx context.Context) ([]v1.UserGroupDetail, error)
	UpdateGroup(ctx context.Context, groupIdOrName string, isAssigned bool) error
	ReadAuth(ctx context.Context) (*v1.UserAuthentication, error)
	UpdateAuth(ctx context.Context, request v1.UserAuthentication) error
}

var _ UserExtraAPI = (*userExtraOp)(nil)

type userExtraOp struct {
	client *v1.Client
	userId uuid.UUID
}

func NewUserExtraOp(client *v1.Client, userId uuid.UUID) UserExtraAPI {
	return &userExtraOp{client: client, userId: userId}
}

func (op *userExtraOp) ListGroup(ctx context.Context) ([]v1.UserGroupDetail, error) {
	res, err := op.client.GetUserGroup(ctx, v1.GetUserGroupParams{UserId: op.userId})
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetUserGroupOK:
		return p.Apigw.Groups, nil
	case *v1.GetUserGroupBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUserGroupInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *userExtraOp) UpdateGroup(ctx context.Context, groupIdOrName string, isAssigned bool) error {
	var req []byte
	idOrName, err := uuid.Parse(groupIdOrName)
	if err != nil {
		temp := []struct {
			IsAssigned bool   `json:"isAssigned"`
			Name       string `json:"name"`
		}{{
			isAssigned,
			groupIdOrName,
		}}
		req, _ = json.Marshal(temp)
	} else {
		temp := []struct {
			IsAssigned bool      `json:"isAssigned"`
			Id         uuid.UUID `json:"id"`
		}{{
			isAssigned,
			idOrName,
		}}
		req, _ = json.Marshal(temp)
	}

	res, err := op.client.UpdateUserGroup(ctx, req, v1.UpdateUserGroupParams{UserId: op.userId})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateUserGroupNoContent:
		return nil
	case *v1.UpdateUserGroupBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateUserGroupNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateUserGroupInternalServerError:
		return errors.New(p.Message.Value)
	}

	return nil
}

func (op *userExtraOp) ReadAuth(ctx context.Context) (*v1.UserAuthentication, error) {
	res, err := op.client.GetUserAuthentication(ctx, v1.GetUserAuthenticationParams{UserId: op.userId})
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetUserAuthenticationOK:
		return &p.Apigw.UserAuthentication.Value, nil
	case *v1.GetUserAuthenticationBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUserAuthenticationNotFound:
		return nil, errors.New(p.Message.Value)
	case *v1.GetUserAuthenticationInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}
func (op *userExtraOp) UpdateAuth(ctx context.Context, request v1.UserAuthentication) error {
	res, err := op.client.UpsertUserAuthentication(ctx, v1.NewOptUserAuthentication(request),
		v1.UpsertUserAuthenticationParams{UserId: op.userId})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpsertUserAuthenticationNoContent:
		return nil
	case *v1.UpsertUserAuthenticationBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpsertUserAuthenticationNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpsertUserAuthenticationInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}
