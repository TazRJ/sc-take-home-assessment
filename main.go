package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
)

func main() {
	req := &folders.FetchFolderRequest{
		OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
	}

	res, err := folders.GetAllFolders(req)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	folders.PrettyPrint(res)

	// req := &folders.PaginatedFetchReq{
	// 	OrgID:  uuid.FromStringOrNil(folders.DefaultOrgID),
	// 	Limit:  20,
	// 	Cursor: "",
	// }

	// for {
	// 	res, err := folders.GetAllFoldersPaginated(req)
	// 	if err != nil {
	// 		// Error handling
	// 		fmt.Printf("%v", err)
	// 		return
	// 	}

	// 	folders.PrettyPrint(res.Folders)

	// 	// break loop if there aren't any more pages
	// 	if res.NextCursor == "" {
	// 		break
	// 	}

	// 	// Updates the cursor to fetch the next page afterwards
	// 	req.Cursor = res.NextCursor
	// }
}
