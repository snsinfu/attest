package command

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCmd_cat(t *testing.T) {
	cat, err := Run([]string{"cat"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer cat.Close()

	sequence := []byte{
		109, 112, 199, 223, 57, 115, 237, 11, 168, 210,
		219, 63, 249, 235, 19, 164, 157, 153, 5, 104,
	}

	if _, err := cat.Write(sequence); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := cat.WriteEnd(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	output, err := ioutil.ReadAll(cat)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(output, sequence) {
		t.Fatalf("unexpected output: actual %v expect %v", output, sequence)
	}

	if err := cat.Wait(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	s := cat.Status()
	if ec := s.ExitCode(); ec != 0 {
		t.Fatalf("unexpected exit code: actual %d expect %d", ec, 0)
	}
}
