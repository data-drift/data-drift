package helpers

import (
	"encoding/json"
	"io/ioutil"
)

func WriteMetadataToFile(metadata interface{}, filename string) error {
	// Convert metadata to JSON-encoded byte array
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Write byte array to file
	err = ioutil.WriteFile(filename, metadataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
