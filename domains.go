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

type DomainAPI interface {
	List(ctx context.Context) ([]v1.Domain, error)
	Create(ctx context.Context, request *v1.Domain) (*v1.Domain, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, request *v1.DomainPUT, id uuid.UUID) error
}

var _ DomainAPI = (*domainOp)(nil)

type domainOp struct {
	client *v1.Client
}

func NewDomainOp(client *v1.Client) DomainAPI {
	return &domainOp{client: client}
}

func (op *domainOp) List(ctx context.Context) ([]v1.Domain, error) {
	res, err := op.client.GetDomains(ctx)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetDomainsOK:
		// ogenが直接デコードできないため、jxを使用して手動でデコード。将来的には修正される可能性あり
		d := jx.DecodeBytes(p.Apigw)
		domains := make([]v1.Domain, 0)
		if err := d.Obj(func(d *jx.Decoder, key string) error {
			switch key {
			case "domains":
				if err := d.Arr(func(d *jx.Decoder) error {
					var domain v1.Domain
					if err := domain.Decode(d); err != nil {
						return err
					}
					domains = append(domains, domain)
					return nil
				}); err != nil {
					return err
				}
				return nil
			default:
				return d.Skip()
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to decode GetDomains response: %w", err)
		}
		return domains, nil
	case *v1.GetDomainsBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetDomainsUnauthorized:
		return nil, errors.New(p.Message.Value)
	case *v1.GetDomainsInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *domainOp) Create(ctx context.Context, request *v1.Domain) (*v1.Domain, error) {
	res, err := op.client.AddDomain(ctx, request)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.AddDomainCreated:
		return &p.Apigw.Domain.Value, nil
	case *v1.AddDomainBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.AddDomainConflict:
		return nil, errors.New(p.Message.Value)
	case *v1.AddDomainUnauthorized:
		return nil, errors.New(p.Message.Value)
	case *v1.AddDomainInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *domainOp) Update(ctx context.Context, request *v1.DomainPUT, id uuid.UUID) error {
	res, err := op.client.UpdateDomain(ctx, request, v1.UpdateDomainParams{DomainId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateDomainNoContent:
		return nil
	case *v1.UpdateDomainBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateDomainNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateDomainConflict:
		return errors.New(p.Message.Value)
	case *v1.UpdateDomainUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.UpdateDomainInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}

func (op *domainOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteDomain(ctx, v1.DeleteDomainParams{DomainId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.DeleteDomainNoContent:
		return nil
	case *v1.DeleteDomainBadRequest:
		return errors.New(p.Message.Value)
	case *v1.DeleteDomainNotFound:
		return errors.New(p.Message.Value)
	case *v1.DeleteDomainUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.DeleteDomainInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}
