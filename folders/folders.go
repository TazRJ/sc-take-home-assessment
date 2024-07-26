package folders

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Takes in a request to fetch folders and returns a response or an error

// func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {

// 1. Declares local variables:
// err -> error
// f1 -> var of type Folder
// fs -> slice of pointers to Folder

// var (
// 	err error
// 	f1  Folder
// 	fs  []*Folder
// )

// 2. Code snippet below calls FetchAllFoldersByOrgID with req.OrgID to fetch all folders associated with the given organisation ID. The result is stored in the var r. f -> slice of Folder

// f := []Folder{}
// r, _ := FetchAllFoldersByOrgID(req.OrgID)

// 3. The below code snippet iterates over the result r and appends the values to the slice f

// for k, v := range r {
// 	f = append(f, *v)
// }

// 4. Iterates over the slice f and appends pointers to each folder to the slice fp

// var fp []*Folder
// for k1, v1 := range f {
// 	fp = append(fp, &v1)
// }

// 5. Creates a new FetchFolderResponse with the slice of pointers to folders fp

// var ffr *FetchFolderResponse
// ffr = &FetchFolderResponse{Folders: fp}

// 6. Returns the response ffr and nil for the error
// return ffr, nil
//}

func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// Error handling: Check if the request is nil and return an error if it is.
	if req == nil {
		return nil, errors.New("received a nil request. Please provide a valid request")
	}

	if req.OrgID == uuid.Nil {
		return nil, errors.New("invalid orgID: Nil UUID")
	}
	// Fetch all folders associated with the given OrgID. This is more efficient because it directly fetches the folders without intermediate steps.
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// Return the fetched folders wrapped in a FetchFolderResponse, which is more efficient because it avoids unnecessary iterations and memory allocations.
	return &FetchFolderResponse{Folders: folders}, nil
}

// Preallocation of memory for the slice of folders can reduce the number of allocations, improve performance and reduce memory usage.
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	allFolders := GetSampleData()

	// Preallocate slice with an estimated size to avoid multiple allocations.
	orgFolders := make([]*Folder, 0, len(allFolders))

	for _, folder := range allFolders {
		if folder.OrgId == orgID {
			orgFolders = append(orgFolders, folder)
		}
	}
	return orgFolders, nil
}
