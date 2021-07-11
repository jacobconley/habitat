package e2e

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"gopkg.in/h2non/baloo.v3"
)



func ReadAsset(t *testing.T, asset string) ([]byte, string) { 

	file, err := os.Open("../test-fixtures/userland/assets/" + asset)

	if err != nil { 
		t.Log("Error opening file", err) 
		t.Fail()
	}

	var buffer bytes.Buffer 
	hasher := sha256.New()
	
	reader := bufio.NewReader(file)
	rawbuf := make([]byte, 256) 

	for { 
		n, err := reader.Read(rawbuf)

		if n == 0 { break } 

		if err != nil { 
			t.Log("Read err", err) 
			t.Fail()
		}


		buffer.Write(rawbuf[0:n])
		hasher.Write(rawbuf[0:n])
	}


	return buffer.Bytes(), fmt.Sprintf("%x", hasher.Sum(nil))
}






var test = baloo.New(":3000")

func TestAsset404(t *testing.T) { 
	test.Get("/!/assets/asset-that-dont-exist").
		Expect(t). 
		Status(404). 
		Done()
}


func TestAssetValid(t *testing.T) { 
	asset := "test.txt"
	bytes, digest := ReadAsset(t, asset) 

	url := fmt.Sprintf("/!/assets/%s:%s", asset, digest)
	t.Log(url)

	test.Get(url).
		Expect(t). 
		Status(200). 
		BodyEquals( string(bytes) ).
		Done()
}
