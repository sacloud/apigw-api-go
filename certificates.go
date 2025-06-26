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

type CertificateAPI interface {
	List(ctx context.Context) ([]v1.Certificate, error)
	Create(ctx context.Context, request *v1.Certificate) (*v1.Certificate, error)
	Update(ctx context.Context, request *v1.Certificate, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ CertificateAPI = (*certificateOp)(nil)

type certificateOp struct {
	client *v1.Client
}

func NewCertificateOp(client *v1.Client) CertificateAPI {
	return &certificateOp{client: client}
}

func (op *certificateOp) List(ctx context.Context) ([]v1.Certificate, error) {
	res, err := op.client.GetCertificates(ctx)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.GetCertificatesOK:
		// ogenが直接デコードできないため、jxを使用して手動でデコード。将来的には修正される可能性あり
		d := jx.DecodeBytes(p.Apigw)
		certificates := make([]v1.Certificate, 0)
		if err := d.Obj(func(d *jx.Decoder, key string) error {
			switch key {
			case "certificates":
				if err := d.Arr(func(d *jx.Decoder) error {
					var certificate v1.Certificate
					if err := certificate.Decode(d); err != nil {
						return err
					}
					certificates = append(certificates, certificate)
					return nil
				}); err != nil {
					return err
				}
				return nil
			default:
				return d.Skip()
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to decode GetCertificates response: %w", err)
		}
		return certificates, nil
	case *v1.GetCertificatesBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.GetCertificatesUnauthorized:
		return nil, errors.New(p.Message.Value)
	case *v1.GetCertificatesInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *certificateOp) Create(ctx context.Context, request *v1.Certificate) (*v1.Certificate, error) {
	res, err := op.client.AddCertificate(ctx, request)
	if err != nil {
		return nil, err
	}

	switch p := res.(type) {
	case *v1.AddCertificateCreated:
		return &p.Apigw.Certificate.Value, nil
	case *v1.AddCertificateBadRequest:
		return nil, errors.New(p.Message.Value)
	case *v1.AddCertificateConflict:
		return nil, errors.New(p.Message.Value)
	case *v1.AddCertificateUnauthorized:
		return nil, errors.New(p.Message.Value)
	case *v1.AddCertificateInternalServerError:
		return nil, errors.New(p.Message.Value)
	}

	return nil, errors.New("unexpected response type")
}

func (op *certificateOp) Update(ctx context.Context, request *v1.Certificate, id uuid.UUID) error {
	res, err := op.client.UpdateCertificate(ctx, request, v1.UpdateCertificateParams{CertificateId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.UpdateCertificateNoContent:
		return nil
	case *v1.UpdateCertificateBadRequest:
		return errors.New(p.Message.Value)
	case *v1.UpdateCertificateNotFound:
		return errors.New(p.Message.Value)
	case *v1.UpdateCertificateConflict:
		return errors.New(p.Message.Value)
	case *v1.UpdateCertificateUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.UpdateCertificateInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}

func (op *certificateOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteCertificate(ctx, v1.DeleteCertificateParams{CertificateId: id})
	if err != nil {
		return err
	}

	switch p := res.(type) {
	case *v1.DeleteCertificateNoContent:
		return nil
	case *v1.DeleteCertificateBadRequest:
		return errors.New(p.Message.Value)
	case *v1.DeleteCertificateNotFound:
		return errors.New(p.Message.Value)
	case *v1.DeleteCertificateUnauthorized:
		return errors.New(p.Message.Value)
	case *v1.DeleteCertificateInternalServerError:
		return errors.New(p.Message.Value)
	}

	return errors.New("unexpected response type")
}
