module "github.com/blueimp/passphrase/appengine"

require (
	"github.com/blueimp/passphrase" v1.0.0
	"github.com/golang/protobuf" v1.0.0
	"golang.org/x/net" v0.0.0-20180320002117-6078986fec03
	"google.golang.org/appengine" v1.0.0
)

replace "github.com/blueimp/passphrase" v1.0.0 => "../"
