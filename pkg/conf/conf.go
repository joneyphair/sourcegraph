package conf

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sourcegraph/jsonx"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/globals"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/globals/confserver"
	"github.com/sourcegraph/sourcegraph/pkg/conf/store"
)

type configurationMode int

const (
	// The user of pkg/conf reads and writes to the configuration file.
	// This should only ever be used by frontend.
	modeServer configurationMode = iota

	// The user of pkg/conf only reads the configuration file.
	modeClient

	// The user of pkg/conf is a test case.
	modeTest
)

func getMode() configurationMode {
	mode := os.Getenv("CONFIGURATION_MODE")

	switch mode {
	case "server":
		return modeServer
	case "client":
		return modeClient
	default:
		// Detect 'go test' and default to test mode in that case.
		p, err := os.Executable()
		if err == nil && filepath.Ext(p) == ".test" {
			return modeTest
		}

		// Otherwise we default to client mode, so that most services need not
		// specify CONFIGURATION_MODE=client explicitly.
		return modeClient
	}
}

func init() {
	clientBasicStore := store.NewBasicStore()
	clientCoreStore := store.NewCoreStore()

	defaultClient = &client{
		basicStore:   clientBasicStore,
		coreStore:    clientCoreStore,
		basicFetcher: httpBasicFetcher{},
		coreFetcher:  httpCoreFetcher{},
	}

	mode := getMode()

	// Don't kickoff the background updaters for the client/server
	// when running test cases.
	if mode == modeTest {
		// Seed the client store with a dummy configuration for test cases.
		dummyConfig := `
		{
			// This is an empty configuration to run test cases.
		}`

		_, err := clientBasicStore.MaybeUpdate(dummyConfig)

		if err != nil {
			log.Fatalf("received error when setting up the basic store for the default client during test, err :%s", err)
		}

		_, err = clientCoreStore.MaybeUpdate(dummyConfig)

		if err != nil {
			log.Fatalf("received error when setting up the basic store for the default client during test, err :%s", err)
		}

		return
	}

	// If the caller of pkg/conf is the frontend service, instantiate the DefaultServerFrontendOnly
	// and install the passthrough fetcher for defaultClient in order to avoid deadlock issues.
	if mode == modeServer {
		globals.ConfigurationServerFrontendOnly = confserver.NewServer(os.Getenv("SOURCEGRAPH_CONFIG_FILE"), os.Getenv("SOURCEGRAPH_CONFIG_CORE_FILE"))

		globals.ConfigurationServerFrontendOnly.Start()
		defaultClient.basicFetcher = passthroughBasicFetcherFrontendOnly{}
		defaultClient.coreFetcher = passthroughCoreFetcherFrontendOnly{}
	}

	go defaultClient.continuouslyUpdate()
}

// FormatOptions is the default format options that should be used for jsonx
// edit computation.
var FormatOptions = jsonx.FormatOptions{InsertSpaces: true, TabSize: 2, EOL: "\n"}
