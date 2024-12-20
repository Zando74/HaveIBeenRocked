package grpc_controller

import (
	"io"
	"sync"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/usecase"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/logger"
)

var (
	log = logger.LoggerSingleton.GetInstance()
)

type PasswordListUploadServer struct {
	grpc_proto.UnimplementedRawPasswordListUploadServer
	ProcessPasswordBatchUseCase *usecase.ProcessPasswordBatch
}

func (s *PasswordListUploadServer) processBatchAsync(passwords [][]byte, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	err := s.ProcessPasswordBatchUseCase.Execute(passwords)
	if err != nil {
		errChan <- err
	}
}

func (s *PasswordListUploadServer) UploadRawPasswordList(stream grpc_proto.RawPasswordListUpload_UploadRawPasswordListServer) error {
	wg := sync.WaitGroup{}

	errChan := make(chan error)
	errorThrown := false

	go func() {
		for err := range errChan {
			log.Error(err)
			errorThrown = true
		}
	}()

	log.Info(logger.OpeningStream)
	for {
		log.Debug(logger.ReceivingBatchOfPasswords)
		passwords, err := stream.Recv()

		if err != nil {
			wg.Wait()

			close(errChan)

			if err == io.EOF {
				switch errorThrown {
				case true:
					log.Error(logger.UploadCompleteWithErrors)
					return stream.SendAndClose(&grpc_proto.Status{Success: false, Message: "Upload complete with errors"})
				default:
					log.Info(logger.UploadCompleteSuccessfully)
					return stream.SendAndClose(&grpc_proto.Status{Success: true, Message: "Upload complete"})
				}
			}
			log.Error(logger.UploadCompleteWithErrors)
			return stream.SendAndClose(&grpc_proto.Status{Success: false, Message: err.Error()})

		}

		passwordBytes := passwords.GetPasswords()

		if len(passwordBytes) > 0 {
			wg.Add(1)
			go s.processBatchAsync(passwordBytes, errChan, &wg)
		}

	}

}
