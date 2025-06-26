// Copyright 2025- The sacloud/kms-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apigw_test

import (
	"context"
	"os"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET")(t)

	client, err := apigw.NewClient()
	require.Nil(t, err)

	ctx := context.Background()
	subscriptionOp := apigw.NewSubscriptionOp(client)

	subscriptions, err := subscriptionOp.ListPlans(ctx)
	require.Nil(t, err)
	require.Greater(t, len(subscriptions), 0)

	if os.Getenv("ENABLE_APIGW_SUBSCRIPTION_TEST") != "1" {
		return
	}

	err = subscriptionOp.Create(ctx, subscriptions[0].ID.Value)
	require.Nil(t, err)

	status, err := subscriptionOp.Read(ctx)
	assert.Nil(t, err)
	assert.Equal(t, string(status.Type), "Subscribed")

	err = subscriptionOp.Delete(ctx)
	assert.Nil(t, err)
}
