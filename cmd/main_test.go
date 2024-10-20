package main

import (
	"bytes"
	"errors"
	"flag"
	"heartbeat/cmd/version"
	"io/fs"
	"log"
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(t *testing.T) {
	t.Run("when error not happened", func(*testing.T) {
		os.Args = []string{os.Args[0], "-config=../fixtures/config.yml"}
		signals <- syscall.SIGINT

		assert.NotPanics(t, main)
	})

	t.Run("when error happened", func(t *testing.T) {
		defer func() { logFatalf = log.Fatalf }()
		os.Args = []string{os.Args[0], "-config=config.yml"}
		logMock := new(logMock)
		logFatalf = logMock.Fatalf
		expectedError := &fs.PathError{
			Op:   "open",
			Path: "config.yml",
			Err:  errors.New("no such file or directory"),
		}
		logMock.On("Fatalf", "%s\n", mock.MatchedBy(func(arg interface{}) bool {
			err, ok := arg.([]interface{})[0].(*fs.PathError)
			return ok &&
				err.Op == expectedError.Op &&
				err.Path == expectedError.Path &&
				err.Err.Error() == expectedError.Err.Error()
		})).Once()
		main()

		logMock.AssertExpectations(t)
	})
}

func TestRun(t *testing.T) {
	path, config := "some-path-to-the-program", "-config=../fixtures/config_without_instances.yml"

	t.Run("when command line argument error", func(t *testing.T) {
		assert.Error(t, run([]string{path, "-not_existing_flag=42"}, flag.ContinueOnError))
	})

	t.Run("when configuration error", func(t *testing.T) {
		assert.Error(t, run([]string{path, "-config=wrong_path"}))
	})

	t.Run("when server starting error", func(t *testing.T) {
		listener, _ := net.Listen("tcp", ":8080")
		defer listener.Close()

		assert.Error(t, run([]string{path, config}))
	})

	t.Run("when server was started successfully, interrupt signal (exit 2) received", func(t *testing.T) {
		signals <- syscall.SIGINT

		assert.NoError(t, run([]string{path, config}))
	})

	t.Run("when server was started successfully, quit signal (exit 3) received", func(t *testing.T) {
		signals <- syscall.SIGQUIT

		assert.NoError(t, run([]string{path, config}))
	})

	t.Run("when server was started successfully, terminated signal (exit 15) received", func(t *testing.T) {
		signals <- syscall.SIGTERM

		assert.NoError(t, run([]string{path, config}))
	})

	t.Run("when version flag passed", func(t *testing.T) {
		assert.NoError(t, run([]string{path, "-v"}))
	})
}

func TestAttrFromCommandLine(t *testing.T) {
	t.Run("when known flags found creates pointer to ConfigurationAttr based on passed command line arguments", func(t *testing.T) {
		config := "path/to/config.yml"
		ver, configPath, err := attrFromCommandLine(
			[]string{
				"some-path-to-the-program",
				"-v",
				"-config=" + config,
			},
		)

		assert.True(t, ver)
		assert.Equal(t, config, configPath)
		assert.NoError(t, err)
	})

	t.Run("when unknown flags found sends exit signal", func(t *testing.T) {
		ver, configPath, err := attrFromCommandLine([]string{"some-path-to-the-program", "-notKnownFlag"}, flag.ContinueOnError)

		assert.False(t, ver)
		assert.Empty(t, configPath)
		assert.Error(t, err)
	})
}

func TestPrintVersionData(t *testing.T) {
	t.Run("composes version data", func(t *testing.T) {
		bytesBuffer := new(bytes.Buffer)
		printVersionData(bytesBuffer)
		ver := "heartbeat: " + version.Version + "\n"
		commit := "commit: " + version.GitCommit + "\n"
		builtAt := "built at: " + version.BuildTime + "\n"
		versionData := ver + commit + builtAt

		assert.Equal(t, versionData, bytesBuffer.String())
	})
}
