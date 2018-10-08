package stns

import (
	"testing"
	"google.golang.org/grpc/credentials"
	"github.com/jarcoal/httpmock"
	"bytes"
)

func Test_getPubKeyFromSTNS(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://127.0.0.1:1104/v1/users?name=test",
		httpmock.NewStringResponder(200, `[{"id":0,"name":"test","password":"","group_id":0,"directory":"","shell":"","gecos":"","keys":["ssh-rsa PublicKey"]}]`))

	tc := stnsTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion:  "1.0",
		},
		stnsAddress: "127.0.0.1",
		stnsPort: "1104",
	}

	key, err := tc.getPubKeyFromSTNS("test")

	if err != nil || !bytes.Equal(key, []byte("ssh-rsa PublicKey")) {
		t.Errorf("Failed to get PublicKey from STNS")
	}
}
