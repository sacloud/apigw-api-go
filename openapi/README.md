APIゲートウェイのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/apigw/

apigw-api-goではここで公開されている定義からogenが未サポートの機能を削除した定義を利用しています


## OpenAPI定義のdiff

以下の問題に対処するための暫定的な修正:

- 現状ogenが複雑な`anyOf`を処理できないケース
- 現状ogenが`array`に対する`default`を処理できないケース
- OpenAPI定義のtypo群

```diff
diff --git a/openapi/openapi.json b/openapi/openapi.json
index 72f29b0..e75bc4c 100644
--- a/openapi/openapi.json
+++ b/openapi/openapi.json
@@ -1126,7 +1126,7 @@
                     "apigw": {
                       "type": "object",
                       "properties": {
-                        "userAuthorization": {
+                        "userAuthentication": {
                           "$ref": "#/components/schemas/UserAuthentication"
                         }
                       }
@@ -1484,7 +1484,22 @@
             "content": {
               "application/json": {
                 "schema": {
-                  "$ref": "#/components/schemas/DomainDTO"
+                  "type": "object",
+                  "required": [
+                    "apigw"
+                  ],
+                  "properties": {
+                    "apigw": {
+                      "properties": {
+                        "domains": {
+                          "type": "array",
+                          "items": {
+                            "$ref": "#/components/schemas/DomainDTO"
+                          }
+                        }
+                      }
+                    }
+                  }
                 }
               }
             }
@@ -1649,7 +1664,22 @@
             "content": {
               "application/json": {
                 "schema": {
-                  "$ref": "#/components/schemas/CertificateDTO"
+                  "type": "object",
+                  "required": [
+                    "apigw"
+                  ],
+                  "properties": {
+                    "apigw": {
+                      "properties": {
+                        "certificates": {
+                          "type": "array",
+                          "items": {
+                            "$ref": "#/components/schemas/CertificateDTO"
+                          }
+                        }
+                      }
+                    }
+                  }
                 }
               }
             }
@@ -2526,17 +2556,6 @@
                   "CONNECT",
                   "TRACE"
                 ],
-                "default": [
-                  "GET",
-                  "POST",
-                  "PUT",
-                  "DELETE",
-                  "PATCH",
-                  "OPTIONS",
-                  "HEAD",
-                  "CONNECT",
-                  "TRACE"
-                ],
                 "description": "CORS許可メソッド<br>未指定の場合は全メソッドを許可"
               },
               "accessControlAllowOrigins": {
@@ -2689,17 +2708,6 @@
                   "GET",
                   "POST"
                 ],
-                "default": [
-                  "GET",
-                  "POST",
-                  "PUT",
-                  "DELETE",
-                  "PATCH",
-                  "OPTIONS",
-                  "HEAD",
-                  "CONNECT",
-                  "TRACE"
-                ],
                 "description": "RouteにアクセスするためのHTTPメソッド<br>未指定の場合は全メソッドを許可"
               },
               "httpsRedirectStatusCode": {
@@ -3043,32 +3051,6 @@
                 ]
               }
             }
-          },
-          {
-            "anyOf": [
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "rsa"
-                ]
-              },
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "ecdsa"
-                ]
-              },
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "rsa",
-                  "ecdsa"
-                ]
-              }
-            ]
           }
         ]
       },

```

## 生成されたコードのdiff

以下の問題に対処するための暫定的な修正:

- APIの裏側が`{}`を`[]`に自動変換する影響でCertificateのレスポンスのパースに失敗する不具合
- ogenが現状`writeOnly`を認識しない不具合に対するrequiredチェックの無効化
- AddServiceのステータスコードが201ではなく200になっている不具合

```diff
diff --git a/apis/v1/oas_json_gen.go b/apis/v1/oas_json_gen.go
index e5599b6..8a71ded 100644
--- a/apis/v1/oas_json_gen.go
+++ b/apis/v1/oas_json_gen.go
@@ -2518,7 +2518,7 @@ func (s *BasicAuth) Decode(d *jx.Decoder) error {
 	// Validate required fields.
 	var failures []validate.FieldError
 	for i, mask := range [1]uint8{
-		0b00011000,
+		0b00001000,
 	} {
 		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
 			// Mask only required fields and check equality to mask using XOR.
@@ -2790,7 +2790,13 @@ func (s *CertificateDetails) Decode(d *jx.Decoder) error {
 		}
 		return nil
 	}); err != nil {
-		return errors.Wrap(err, "decode CertificateDetails")
+		// ecdsaを指定していない場合には{}が返ってくるはずだが、現状APIの裏側で変換する処理が入ってしまい[]が返ってくるので、
+		// 修正されるまでそれを無視する
+		if errArr := d.Arr(func(d *jx.Decoder) error { return nil }); errArr == nil {
+			return nil
+		} else {
+			return errors.Wrap(err, "decode CertificateDetails")
+		}
 	}
 
 	return nil
@@ -8799,7 +8805,7 @@ func (s *HmacAuth) Decode(d *jx.Decoder) error {
 	// Validate required fields.
 	var failures []validate.FieldError
 	for i, mask := range [1]uint8{
-		0b00011000,
+		0b00001000,
 	} {
 		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
 			// Mask only required fields and check equality to mask using XOR.
@@ -9018,7 +9024,7 @@ func (s *Jwt) Decode(d *jx.Decoder) error {
 	// Validate required fields.
 	var failures []validate.FieldError
 	for i, mask := range [1]uint8{
-		0b00111000,
+		0b00101000,
 	} {
 		if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
 			// Mask only required fields and check equality to mask using XOR.
diff --git a/apis/v1/oas_response_decoders_gen.go b/apis/v1/oas_response_decoders_gen.go
index da25f54..c4071b4 100644
--- a/apis/v1/oas_response_decoders_gen.go
+++ b/apis/v1/oas_response_decoders_gen.go
@@ -1097,7 +1097,7 @@ func decodeAddRouteResponse(resp *http.Response) (res AddRouteRes, _ error) {
 
 func decodeAddServiceResponse(resp *http.Response) (res AddServiceRes, _ error) {
 	switch resp.StatusCode {
-	case 200:
+	case 201:
 		// Code 200.
 		ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
 		if err != nil {

```