package logging

import (
	"io"

	"gopkg.in/Graylog2/go-gelf.v1/gelf"
)

// UseGraylogOutput enables Graylog destination to the logger, removing the file output
func (l *Log) UseGraylogOutput(serverAddress string) error {
	graylog, err := gelf.NewWriter(serverAddress)
	if err != nil {
		return err
	}
	l.multiplexToStdOutAnd(graylog)
	return nil
}

// WithStubbedGraylogOutput is a test function to mock a Greylog server
func (l *Log) WithStubbedGraylogOutput(remoteServer io.Writer) error {
	l.logger.SetOutput(remoteServer)
	return nil
}
