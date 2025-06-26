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

type ServiceAPI interface {
	List(ctx context.Context) ([]v1.ServiceDetail, error)
	Create(ctx context.Context, request *v1.ServiceDetail) (*v1.ServiceDetail, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.ServiceDetail, error)
	Update(ctx context.Context, request *v1.ServiceDetail, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ ServiceAPI = (*serviceOp)(nil)

type serviceOp struct {
	client *v1.Client
}

func NewServiceOp(client *v1.Client) ServiceAPI {
	return &serviceOp{client: client}
}

func (op *serviceOp) List(ctx context.Context) ([]v1.ServiceDetail, error) {
	res, err := op.client.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetServicesOK:
		// ogenが直接デコードできないため、jxを使用して手動でデコード。将来的には修正される可能性あり
		d := jx.DecodeBytes(p.Apigw)
		services := make([]v1.ServiceDetail, 0)
		if err := d.Obj(func(d *jx.Decoder, key string) error {
			switch key {
			case "services":
				if err := d.Arr(func(d *jx.Decoder) error {
					var service v1.ServiceDetail
					if err := service.Decode(d); err != nil {
						return err
					}
					services = append(services, service)
					return nil
				}); err != nil {
					return err
				}
				return nil
			default:
				return d.Skip()
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to decode GetServiceRoutes response: %w", err)
		}
		return services, nil
	case *v1.GetServicesBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetServicesUnauthorized:
		return nil, errors.New(p.Message.Value)
	case *v1.GetServicesInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *serviceOp) Create(ctx context.Context, request *v1.ServiceDetail) (*v1.ServiceDetail, error) {
	res, err := op.client.AddService(ctx, request)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.AddServiceOK:
		return &p.Apigw.Service.Value, nil
	case *v1.AddServiceBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.AddServiceConflict:
		return nil, errors.New(p.Message.Value)
	case *v1.AddServiceNotFound:
		return nil, errors.New(p.Message.Value)
	case *v1.AddServiceInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *serviceOp) Read(ctx context.Context, id uuid.UUID) (*v1.ServiceDetail, error) {
	res, err := op.client.GetServiceById(ctx, v1.GetServiceByIdParams{ServiceId: id})
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetServiceByIdOK:
		return &p.Apigw.Service.Value, nil
	case *v1.GetServiceByIdBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetServiceByIdNotFound:
		return nil, errors.New(p.Message.Value)
	case *v1.GetServiceByIdInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *serviceOp) Update(ctx context.Context, request *v1.ServiceDetail, id uuid.UUID) error {
	res, err := op.client.UpdateService(ctx, request, v1.UpdateServiceParams{ServiceId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateServiceNoContent:
		return nil
	case *v1.UpdateServiceBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateServiceNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateServiceInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}

func (op *serviceOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteService(ctx, v1.DeleteServiceParams{ServiceId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.DeleteServiceNoContent:
		return nil
	case *v1.DeleteServiceBadRequest:
		return errors.New(p.Message.Value)
	case *v1.DeleteServiceNotFound:
		return errors.New(p.Message.Value)
	case *v1.DeleteServiceUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.DeleteServiceInternalServerError:
		return errors.New(p.Message.Value)
	}

	return nil
}
