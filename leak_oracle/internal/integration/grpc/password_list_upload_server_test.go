package grpc_test

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/grpc_proto"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/value_object"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/integration/mock"
	"github.com/stretchr/testify/assert"
)

var (
	PasswordList = [][]byte{
		[]byte("p@ssw0rd!"),
		[]byte("Passw0rd#123"),
		[]byte("123456@bcd"),
		[]byte("!Q2w#E4r"),
		[]byte("test1234$"),
		[]byte("9ublic$foo>"),
		[]byte("%admin$p@ss"),
	}

	SecondDifferentPasswordList = [][]byte{
		[]byte("p@ssw0rd!BIS"),
		[]byte("Passw0rd#123BIS"),
		[]byte("123456@bcdBIS"),
		[]byte("!Q2w#E4rBIS"),
		[]byte("test1234$BIS"),
		[]byte("9ublic$foo>BIS"),
		[]byte("%admin$p@ssBIS"),
	}

	EmptyPasswordList = [][]byte{}
)

func TestMain(m *testing.M) {
	err := mock.SetupTestPostgresqlDB()
	if err != nil {
		log.Fatal(err)
	}
	postgresql.Init()

	m.Run()
}

func TestPasswordListUploadServer(t *testing.T) {
	conn, closer := InitTestServer()
	defer closer()

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	stream, err := client.UploadRawPasswordList(context.Background())
	assert.NoError(t, err, "error at stream creation")

	err = stream.Send(&grpc_proto.RawPasswordList{Passwords: PasswordList})
	assert.NoError(t, err, "error sending passwords")

	_, err = stream.CloseAndRecv()
	assert.NoError(t, err, "error uploading passwords")

	for _, password := range PasswordList {
		hashedPassword, _ := value_object.NewPasswordHashFromPassword(password)
		leakedHash, err := LeakedHashRepositoryImpl.Retrieve(hashedPassword)
		assert.NoError(t, err, "error retrieving leaked hash for password")
		assert.NotNil(t, leakedHash, "failed to retrieve leaked hash for password")
	}
}

func TestEmptyPasswordListUpload(t *testing.T) {
	conn, closer := InitTestServer()
	defer closer()

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	stream, err := client.UploadRawPasswordList(context.Background())
	assert.NoError(t, err, "error at stream creation")

	err = stream.Send(&grpc_proto.RawPasswordList{Passwords: EmptyPasswordList})
	assert.NoError(t, err, "error sending empty password list")

	_, err = stream.CloseAndRecv()
	assert.NoError(t, err, "error uploading passwords")

	length := LeakedHashRepositoryImpl.Len()
	assert.Equal(t, len(PasswordList), int(length), "expected no more hash db, but got some hashes")
}

func TestDuplicatedPasswords(t *testing.T) {
	conn, closer := InitTestServer()
	defer closer()

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	stream, err := client.UploadRawPasswordList(context.Background())
	assert.NoError(t, err, "error at stream creation")

	err = stream.Send(&grpc_proto.RawPasswordList{Passwords: PasswordList})
	assert.NoError(t, err, "error sending passwords")

	_, err = stream.CloseAndRecv()
	assert.NoError(t, err, "error uploading passwords")

	length := LeakedHashRepositoryImpl.Len()
	assert.Equal(t, len(PasswordList), int(length), "expected no hash append to db, but got some hashes")
}

func TestConcurrentUploads(t *testing.T) {
	conn, closer := InitTestServer()
	defer closer()

	client := grpc_proto.NewRawPasswordListUploadClient(conn)

	wg := sync.WaitGroup{}

	concurrentFn := func(passwords [][]byte) {
		defer wg.Done()
		stream, err := client.UploadRawPasswordList(context.Background())
		assert.NoError(t, err, "error at stream creation")

		err = stream.Send(&grpc_proto.RawPasswordList{Passwords: passwords})
		assert.NoError(t, err, "error sending passwords concurrently")

		_, err = stream.CloseAndRecv()
		assert.NoError(t, err, "error uploading passwords concurrently")
	}

	passwordList := [][][]byte{PasswordList, SecondDifferentPasswordList}
	DuplicatePasswordList := [][][]byte{PasswordList, SecondDifferentPasswordList}
	for _, passwords := range append(passwordList, DuplicatePasswordList...) {
		wg.Add(1)
		go concurrentFn(passwords)
	}

	wg.Wait()

	expectedLength := len(PasswordList) * len(passwordList)
	length := LeakedHashRepositoryImpl.Len()
	assert.GreaterOrEqual(t, int(length), expectedLength, "concurrent uploads should store all unique passwords")
}
