package auth

import "google.golang.org/grpc/metadata"

const (
	MetadataKeyID  = "id"
	MetadataKeyKey = "key"
)

func ToMetadata(id, key string) map[string]string {
	return map[string]string{
		MetadataKeyID:  id,
		MetadataKeyKey: key,
	}
}

func FromMetadata(metadata metadata.MD) (string, string, bool) {
	idSlice, ok := metadata[MetadataKeyID]
	if !ok || len(idSlice) != 1 {
		return "", "", false
	}

	keySlice, ok := metadata[MetadataKeyKey]
	if !ok || len(keySlice) != 1 {
		return "", "", false
	}

	return idSlice[0], keySlice[0], true
}
