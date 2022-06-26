// link: https://stackoverflow.com/questions/72752383/how-to-mock-kafka-dependencies-and-writer-in-golang/72759821#72759821
package kafka

type Server struct {
	grpcServerPort int
	grpcServer     *grpc.Server
	writer         KafkaWriterWrapper
}

type KafkaWriterWrapper interface {
	Write(msg string) error // Suppose kafka.Writer has func Write(msg string) error
}

type KafkaWriterWrapperImpl struct {
	*kafka.Writer
}

type MockKafkaWriterWrapperImpl struct {
}

func (m *MockKafkaWriterWrapperImpl) Write(msg string) error {
	// logic
	return nil
}
