package server

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)


 var fileHashes = map[string]string{}

// HandleAsset ree
func HandleAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "assets/" + vars["path"]
	hash := vars["hash"]

	log.Infof("Serving `%s` @ %s", path, hash)



	file, err := os.Open(path);
	if err != nil { 
		log.Info("Completed [Not found]")
		http.Error(w, "Asset not found", http.StatusNotFound)
		return 
	}

	defer func() { 
		if err = file.Close(); err != nil { 
			log.Error("Could not close file", err)
		}
	}() 




	checkFingerprint := func(expected string) bool {
		if(hash != expected) { 
			log.Info("Completed [Gone]")
			http.Error(w, "The fingerprint in this link no longer matches the requested file; it must have changed.  Try reloading the page that sent you here?", http.StatusGone)
			return false 
		}

		return true 
	}



	// If there's a cached digest, we save ourself
	cachedDigest, hasCachedDigest := fileHashes[path] 
	if hasCachedDigest { 
		log.Debug("Using cached digest")
		if !checkFingerprint(cachedDigest) { return } 
	}
	// Now hasCachedDigest implies a valid cache 


	reader := bufio.NewReader(file) 
	hasher := sha256.New()

	BUFSIZE := 256 
	buffer := make([]byte, BUFSIZE)
	var bobj bytes.Buffer

	for { 

		n, err := reader.Read(buffer)

		if n == 0 { break }
		if err != nil { 
			log.Error("Error reading buffer", err) 
			return 
		}



		if hasCachedDigest { 
			w.Write( buffer[0:n] ) 
		} else { 
			data := buffer[0:n]
			hasher.Write( data )
			bobj.Write( data )
		}

	}


	if hasCachedDigest { 
		// We've already verified the digest
		// And served the file, b/c of the logic above 
		log.Info("Completed [OK] (HashCache)")
		return 
	}

	digest := fmt.Sprintf("%x", hasher.Sum(nil) )
	// [OPTIMIZATION] Better hash generation  ^ - https://stackoverflow.com/a/66116854 - and we should probably be storing them this way too  

	log.Debug("Caching digest: ", digest) 
	fileHashes[path] = digest  
	if !checkFingerprint( digest ) { return } 

	w.Write(bobj.Bytes())
	log.Info("Completed [OK] (Fingerprinted)")
}
