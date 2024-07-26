package folders

import (
	"encoding/base64" // Used for encoding/decoding the cursor
	"errors"
	"fmt"
	"strconv" // Used for converting strings to ints
	"strings" // Used for string manipulation

	"github.com/gofrs/uuid"
)

type PaginatedFetchReq struct {
	Cursor string
	Limit  int
	OrgID  uuid.UUID
}

type PaginatedFetchRes struct {
	Folders    []*Folder
	NextCursor string
}

func GetAllFoldersPaginated(req *PaginatedFetchReq) (*PaginatedFetchRes, error) {
	// Checks if the request is valid
	if req == nil {
		return nil, errors.New("invalid request: request cannot be nil")
	}

	// Checks if the organization ID is valid
	if req.OrgID == uuid.Nil {
		return nil, errors.New("orgID cannot be nil")
	}

	// Checks if the limit is valid
	if req.Limit <= 0 {
		req.Limit = 10 //default limit
		return nil, errors.New("limit must be greater than 0")
	}

	// Checks if the limit is not too large
	if req.Limit > 100 {
		req.Limit = 100 //default limit
		return nil, errors.New("limit of maximum 100")
	}

	// Decodes the cursor
	startIdx := 0
	if req.Cursor != "" {
		var err error
		startIdx, err = DecodeNextCursor(req.Cursor)
		if err != nil {
			return nil, err
		}
	}

	// Fetches all folders of the organization
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// Determines the end index for fetching the folders
	endIdx := startIdx + req.Limit
	if endIdx > len(folders) {
		endIdx = len(folders)
	}

	// Encodes the next cursor
	nextCursor := ""
	if endIdx != len(folders) {
		nextCursor = EncodeNextCursor(endIdx)
	}

	// Returns the fetched folders and the next cursor
	return &PaginatedFetchRes{Folders: folders[startIdx:endIdx], NextCursor: nextCursor}, nil
}

// Encodes the next cursor.
func EncodeNextCursor(endIdx int) string {
	return base64.StdEncoding.EncodeToString([]byte("next_cursor:" + strconv.Itoa(endIdx)))
}

// Decodes the given cursor
func DecodeNextCursor(encodedCursor string) (int, error) {
	if encodedCursor == "" {
		return 0, nil
	}

	decodedCursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return 0, err
	}

	// Splits the decoded cursor into parts
	parts := strings.Split(string(decodedCursor), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid cursor format")
	}

	// Converts the index part of the cursor into an int
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return index, nil
}
