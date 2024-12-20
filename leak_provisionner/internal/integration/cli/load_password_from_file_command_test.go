package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PasswordList = []string{
	"p@ssw0rd!",
	"Passw0rd#123",
	"123456@bcd",
	"!Q2w#E4r",
	"test1234$",
	"9ublic$foo>",
	"%admin$p@ss",
}

func TestLoadPasswordsFromFileCommand(t *testing.T) {
	rootCmd := InitTestCli()

	currentDir, err := os.Getwd()
	assert.NoError(t, err, "could not get current directory")

	filePath := filepath.Join(currentDir, "../files/passwords.txt")

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-f", filePath})
	err = rootCmd.Execute()

	assert.NoError(t, err, "error executing command")

	for _, password := range PasswordList {
		assert.True(t, PasswordGrpcServerMock.DB[password], "password %s not found in database", password)
	}
}

func TestUnexistingFile(t *testing.T) {
	rootCmd := InitTestCli()

	expected := "no such file or directory"

	currentDir, err := os.Getwd()
	assert.NoError(t, err, "could not get current directory")

	filePath := filepath.Join(currentDir, "../files/unexisting_file.txt")

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-f", filePath})

	err = rootCmd.Execute()
	assert.NotNil(t, err, "expected an error for non-existing file")
	assert.Contains(t, err.Error(), expected, "expected error message to contain: %s", expected)
}

func TestFileWithSomeInvalidUTF8Characters(t *testing.T) {
	rootCmd := InitTestCli()

	currentDir, err := os.Getwd()
	assert.NoError(t, err, "could not get current directory")

	filePath := filepath.Join(currentDir, "../files/invalid_passwords.txt")

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-f", filePath})

	rootCmd.Execute()

	assert.NotEqual(t, 0, len(PasswordGrpcServerMock.DB), "passwords were skipped")
}

func TestEmptyPasswordFile(t *testing.T) {
	rootCmd := InitTestCli()

	PasswordGrpcServerMock.DB = make(map[string]bool)

	currentDir, err := os.Getwd()
	assert.NoError(t, err, "could not get current directory")

	filePath := filepath.Join(currentDir, "../files/empty_passwords.txt")

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-f", filePath})

	err = rootCmd.Execute()
	assert.NoError(t, err, "error executing command with empty file")

	assert.Equal(t, 0, len(PasswordGrpcServerMock.DB), "database should be empty")
}

func TestInvalidCommandArgs(t *testing.T) {
	rootCmd := InitTestCli()

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-x", "wrongFilePath"})

	err := rootCmd.Execute()
	assert.NotNil(t, err, "expected an error with invalid command args")
}

func TestLargePasswordFile(t *testing.T) {
	rootCmd := InitTestCli()

	currentDir, err := os.Getwd()
	assert.NoError(t, err, "could not get current directory")

	filePath := filepath.Join(currentDir, "../files/large_passwords.txt")

	rootCmd.SetArgs([]string{"load_passwords_from_file", "-f", filePath})

	err = rootCmd.Execute()
	assert.NoError(t, err, "error executing command with a large file")
}
