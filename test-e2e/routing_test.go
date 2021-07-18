package e2e

import "testing"

func Test404(t *testing.T) {
	client.Get("/thang-that-dont-exist").
		Expect(t).
		Status(404).
		Done()
}


func TestGet(t *testing.T) { 
	client.Get("/test-get").
		Expect(t).
		Status(200).
		BodyEquals( "succ" ).
		Done()
}