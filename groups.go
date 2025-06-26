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
	"errors"
	"fmt"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
)

type GroupAPI interface {
	List(ctx context.Context) ([]v1.Group, error)
	Create(ctx context.Context, request *v1.Group) (*v1.Group, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.Group, error)
	Update(ctx context.Context, request *v1.Group, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ GroupAPI = (*groupOp)(nil)

type groupOp struct {
	client *v1.Client
}

func NewGroupOp(client *v1.Client) GroupAPI {
	return &groupOp{client: client}
}

func (op *groupOp) List(ctx context.Context) ([]v1.Group, error) {
	res, err := op.client.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetGroupsOK:
		// ogenが直接デコードできないため、jxを使用して手動でデコード。将来的には修正される可能性あり
		d := jx.DecodeBytes(p.Apigw)
		groups := make([]v1.Group, 0)
		if err := d.Obj(func(d *jx.Decoder, key string) error {
			switch key {
			case "groups":
				if err := d.Arr(func(d *jx.Decoder) error {
					var group v1.Group
					if err := group.Decode(d); err != nil {
						return err
					}
					groups = append(groups, group)
					return nil
				}); err != nil {
					return err
				}
				return nil
			default:
				return d.Skip()
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to decode GetGroups response: %w", err)
		}
		return groups, nil
	case *v1.GetGroupsBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetGroupsInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *groupOp) Create(ctx context.Context, request *v1.Group) (*v1.Group, error) {
	res, err := op.client.AddGroup(ctx, request)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.AddGroupCreated:
		return &p.Apigw.Group.Value, nil
	case *v1.AddGroupBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.AddGroupConflict:
		return nil, errors.New(p.Message.Value)
	case *v1.AddGroupInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *groupOp) Read(ctx context.Context, id uuid.UUID) (*v1.Group, error) {
	res, err := op.client.GetGroup(ctx, v1.GetGroupParams{GroupId: id})
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetGroupOK:
		return &p.Apigw.Group.Value, nil
	case *v1.GetGroupBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetGroupNotFound:
		return nil, errors.New(p.Message.Value)
	case *v1.GetGroupInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *groupOp) Update(ctx context.Context, request *v1.Group, id uuid.UUID) error {
	res, err := op.client.UpdateGroup(ctx, request, v1.UpdateGroupParams{GroupId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateGroupNoContent:
		return nil
	case *v1.UpdateGroupBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateGroupNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateGroupInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}

func (op *groupOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteGroup(ctx, v1.DeleteGroupParams{GroupId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.DeleteGroupNoContent:
		return nil
	case *v1.DeleteGroupBadRequest:
		return errors.New(p.Message.Value)
	case *v1.DeleteGroupNotFound:
		return errors.New(p.Message.Value)
	case *v1.DeleteGroupUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.DeleteGroupInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}
