package server

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog/log"
)


 var fileHashes = map[string]string{}

// HandleAsset ree
func HandleAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "assets/" + vars["path"]
	hash := vars["hash"]


	logg := log.With().Str("path", path).Str("hash", hash).Logger()


	file, err := os.Open(path);
	if err != nil { 
		logg.Info().Msg("Completed [Not found]")
		http.Error(w, "Asset not found", http.StatusNotFound)
		return 
	}

	defer func() { 
		if err = file.Close(); err != nil { 
			logg.Warn().Err(err).Msg("Could not close file")
		}
	}() 




	checkFingerprint := func(expected string) bool {
		if(hash != expected) { 
			logg.Info().Msg("Completed [Gone]")
			http.Error(w, "The fingerprint in this link no longer matches the requested file; it must have changed.  Try reloading the page that sent you here?", http.StatusGone)
			return false 
		}

		return true 
	}



	// If there's a cached digest, we save ourself
	cachedDigest, hasCachedDigest := fileHashes[path] 
	if hasCachedDigest { 
		logg.Debug().Msg("Using cached digest")
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
			logg.Err(err).Msg("reading buffer") 
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
		logg.Info().Msg("Completed [OK] (HashCache)")
		return 
	}

	digest := fmt.Sprintf("%x", hasher.Sum(nil) )
	// [OPTIMIZATION] Better hash generation  ^ - https://stackoverflow.com/a/66116854 - and we should probably be storing them this way too  

	logg.Debug().Msgf("Caching digest: %s", digest) 
	fileHashes[path] = digest  
	if !checkFingerprint( digest ) { return } 

	w.Write(bobj.Bytes())
	logg.Info().Msg("Completed [OK] (Fingerprinted)")
}
