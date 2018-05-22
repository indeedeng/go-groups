package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: move the import blocks to test fixtures, this file is pretty unreadable

func TestParse_NoImports(t *testing.T) {
	_, rewritten := parse([]byte(""))
	require.False(t, rewritten)
}

func TestParse_ValidGroupingSortImports(t *testing.T) {
	str :=
		`
package foo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)
`

	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, str, string(result))
}

func TestParse_ExtraGroups(t *testing.T) {
	str :=
		`
package foo

import (
	"io"
	"strings"

	"strconv"
	"fmt"

	"os"
	"io/ioutil"
)
`

	expected :=
		`
package foo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)
`
	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, expected, string(result))
}

func TestParse_ExternalGroups(t *testing.T) {
	str :=
		`
package foo

import (
	"io"
	"strings"
	"strconv"
	"fmt"
	"os"
	"io/ioutil"

	"indeed.com/devops/foobar"

	"indeed.com/gophers/quz"

	"github.com/hashicorp/vault/x"

)
`

	expected :=
		`
package foo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/x"

	"indeed.com/devops/foobar"

	"indeed.com/gophers/quz"
)
`
	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, expected, string(result))
}

func TestParse_SortingExternalGroups(t *testing.T) {
	str :=
		`
package foo

import (
    "net/http"
    "testing"

    "indeed.com/devops/libmarvin/gin/ginhandler"
    "indeed.com/devops/libmarvin/spec"
    "indeed.com/devops/marvlet/api/v2"
 
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/suite"
)
`

	expected :=
		`
package foo

import (
    "net/http"
    "testing"

    "github.com/gin-gonic/gin"

    "github.com/stretchr/testify/suite"

    "indeed.com/devops/libmarvin/gin/ginhandler"
    "indeed.com/devops/libmarvin/spec"
    "indeed.com/devops/marvlet/api/v2"
)
`
	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, expected, string(result))
}

func TestParse_PrefixedGroups(t *testing.T) {
	str :=
		`
package foo

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"indeed.com/gophers/libarc"
	"indeed.com/gophers/libops"
	"indeed.com/gophers/rlog"

	"indeed.com/devops/libmarvin/spec"
	"indeed.com/devops/marvlet/backend"
	"indeed.com/devops/marvlet/backend/indeed"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)
`
	expected :=
		`
package foo

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"indeed.com/devops/libmarvin/spec"
	"indeed.com/devops/marvlet/backend"
	"indeed.com/devops/marvlet/backend/indeed"

	"indeed.com/gophers/libarc"
	"indeed.com/gophers/libops"
	"indeed.com/gophers/rlog"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)
`
	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, expected, string(result))
}

func TestParse_SubdomainImports(t *testing.T) {
	str :=
		`
package foo

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/shoenig/petrify/v4"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)
`

	expected :=
		`
package foo

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/shoenig/petrify/v4"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"google.golang.org/grpc"
)
`
	result, rewritten := parse([]byte(str))

	require.True(t, rewritten)
	require.Equal(t, expected, string(result))
}
