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

type SubscriptionAPI interface {
	ListPlans(ctx context.Context) ([]v1.Plan, error)
	Create(ctx context.Context, id uuid.UUID) error
	Read(ctx context.Context) (*v1.SubscriptionStatusSum, error)
	Delete(ctx context.Context) error
}

var _ SubscriptionAPI = (*subscriptionOp)(nil)

type subscriptionOp struct {
	client *v1.Client
}

func NewSubscriptionOp(client *v1.Client) SubscriptionAPI {
	return &subscriptionOp{client: client}
}

func (op *subscriptionOp) ListPlans(ctx context.Context) ([]v1.Plan, error) {
	res, err := op.client.GetPlans(ctx)
	if err != nil {
		return nil, NewAPIError("Subscription.ListPlans", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetPlansOK:
		d := jx.DecodeBytes(p.Apigw)
		plans := make([]v1.Plan, 0)
		if err := d.Obj(func(d *jx.Decoder, key string) error {
			switch key {
			case "plans":
				if err := d.Arr(func(d *jx.Decoder) error {
					var plan v1.Plan
					if err := plan.Decode(d); err != nil {
						return err
					}
					plans = append(plans, plan)
					return nil
				}); err != nil {
					return err
				}
				return nil
			default:
				return d.Skip()
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to decode GetPlans response: %w", err)
		}
		return plans, nil
	case *v1.ErrorSchema:
		return nil, NewAPIError("Subscription.ListPlans", 400, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Subscription.ListPlans", 0, nil)
}

func (op *subscriptionOp) Create(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.Subscribe(ctx, &v1.SubscriptionOption{PlanId: v1.NewOptUUID(id)})
	if err != nil {
		return NewAPIError("Subscription.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.SubscribeNoContent:
		return nil
	case *v1.SubscribeBadRequest:
		return NewAPIError("Subscription.Create", 400, errors.New(p.Message.Value))
	case *v1.SubscribeUnauthorized:
		return NewAPIError("Subscription.Create", 401, errors.New(p.Message.Value))
	case *v1.SubscribeInternalServerError:
		return NewAPIError("Subscription.Create", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Subscription.Create", 0, nil)
}

func (op *subscriptionOp) Read(ctx context.Context) (*v1.SubscriptionStatusSum, error) {
	res, err := op.client.GetSubscription(ctx)
	if err != nil {
		return nil, NewAPIError("Subscription.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetSubscriptionOK:
		return &p.Apigw.Subscription.Value.OneOf, nil
	case *v1.GetSubscriptionUnauthorized:
		return nil, NewAPIError("Subscription.Read", 401, errors.New(p.Message.Value))
	case *v1.GetSubscriptionInternalServerError:
		return nil, NewAPIError("Subscription.Read", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Subscription.Read", 0, nil)
}

func (op *subscriptionOp) Delete(ctx context.Context) error {
	res, err := op.client.Unsubscribe(ctx)
	if err != nil {
		return NewAPIError("Subscription.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.UnsubscribeNoContent:
		return nil
	case *v1.UnsubscribeBadRequest:
		return NewAPIError("Subscription.Delete", 400, errors.New(p.Message.Value))
	case *v1.UnsubscribeNotFound:
		return NewAPIError("Subscription.Delete", 404, errors.New(p.Message.Value))
	case *v1.UnsubscribeInternalServerError:
		return NewAPIError("Subscription.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Subscription.Delete", 0, nil)
}
