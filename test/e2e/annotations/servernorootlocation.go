/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package annotations

import (
	"net/http"
	"strings"

	"github.com/onsi/ginkgo"

	"k8s.io/ingress-nginx/test/e2e/framework"
)

var _ = framework.DescribeAnnotation("server-no-root-location", func() {
	f := framework.NewDefaultFramework("servernorootlocation")

	ginkgo.BeforeEach(func() {
		f.NewEchoDeployment()
	})

	ginkgo.It(`when using the value (true) and enabling in the annotations`, func() {
		host := "servernorootlocation.foo.com"
		annotations := map[string]string{
			"nginx.ingress.kubernetes.io/server-no-root-location": "true",
			"nginx.ingress.kubernetes.io/server-snippet": `
				location / {
          			return 200;
      			}`,
		}

		ing := framework.NewSingleIngress(host, "/api", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, `location /`) &&
					strings.Contains(server, `return 200;`)
			})

		f.HTTPTestClient().
			GET("/").
			WithHeader("Host", host).
			Expect().
			Status(http.StatusOK)
	})
})