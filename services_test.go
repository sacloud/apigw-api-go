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

package apigw_test

import (
	"context"
	"os"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET", "SAKURACLOUD_TEST_HOST")(t)

	client, err := apigw.NewClient()
	require.NoError(t, err)

	ctx := context.Background()
	serviceOp := apigw.NewServiceOp(client)

	// Create a service for testing
	serviceReq := v1.ServiceDetail{
		Name:     "test-service",
		Host:     os.Getenv("SAKURACLOUD_TEST_HOST"),
		Port:     v1.NewOptInt(80),
		Protocol: "http",
	}
	created, err := serviceOp.Create(ctx, &serviceReq)
	require.NoError(t, err)

	serviceReq.Name = "test-service-updated"
	err = serviceOp.Update(ctx, &serviceReq, created.ID.Value)
	assert.NoError(t, err)

	updated, err := serviceOp.Read(ctx, created.ID.Value)
	assert.NoError(t, err)
	assert.Equal(t, "test-service-updated", string(updated.Name))

	services, err := serviceOp.List(ctx)
	assert.NoError(t, err)
	assert.Greater(t, len(services), 0)

	err = serviceOp.Delete(ctx, created.ID.Value)
	assert.NoError(t, err)
}
