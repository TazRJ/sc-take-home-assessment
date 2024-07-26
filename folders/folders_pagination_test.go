package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllFoldersPaginated(t *testing.T) {

	// Sets up the org id
	orgID := uuid.FromStringOrNil(folders.DefaultOrgID)

	// Defines the test cases
	tests := []struct {
		name    string
		req     *folders.PaginatedFetchReq
		wantErr bool
		setup   func() *folders.PaginatedFetchReq
		check   func(t *testing.T, res *folders.PaginatedFetchRes, err error)
	}{
		{
			name:    "Nil Request Error",
			req:     nil,
			wantErr: true,
		},
		{
			name: "Negative Limit Error",
			req: &folders.PaginatedFetchReq{
				OrgID:  orgID,
				Limit:  -1,
				Cursor: "",
			},
			wantErr: true,
		},
		{
			name: "Zero Limit Error",
			req: &folders.PaginatedFetchReq{
				OrgID:  orgID,
				Limit:  0,
				Cursor: "",
			},
			wantErr: true,
		},
		{
			name: "Exceeding Maximum Limit Error",
			req: &folders.PaginatedFetchReq{
				OrgID: orgID,
				Limit: 150,
			},
			wantErr: true,
		},
		{
			name: "Nil OrgID Error",
			req: &folders.PaginatedFetchReq{
				OrgID:  uuid.Nil,
				Limit:  5,
				Cursor: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid Cursor Token Error",
			req: &folders.PaginatedFetchReq{
				OrgID:  orgID,
				Limit:  5,
				Cursor: "invalidToken",
			},
			wantErr: true,
		},
		{
			name: "Fetch Empty Page",
			setup: func() *folders.PaginatedFetchReq {
				emptyOrgID := uuid.Must(uuid.NewV4())
				return &folders.PaginatedFetchReq{
					OrgID: emptyOrgID,
					Limit: 10,
				}
			},
			check: func(t *testing.T, res *folders.PaginatedFetchRes, err error) {
				assert.NoError(t, err)
				assert.Empty(t, res.Folders)
				assert.Empty(t, res.NextCursor)
			},
		},
		{
			name: "Fetch First 5 Folders",
			req: &folders.PaginatedFetchReq{
				OrgID:  orgID,
				Limit:  5,
				Cursor: "",
			},
			check: func(t *testing.T, res *folders.PaginatedFetchRes, err error) {
				assert.NoError(t, err)
				expected, _ := folders.FetchAllFoldersByOrgID(orgID)
				assert.Equal(t, expected[0:5], res.Folders)
			},
		},
		{
			name: "Fetch Next 5 Folders",
			setup: func() *folders.PaginatedFetchReq {
				firstBatch, _ := folders.GetAllFoldersPaginated(&folders.PaginatedFetchReq{
					OrgID:  orgID,
					Limit:  5,
					Cursor: "",
				})
				return &folders.PaginatedFetchReq{
					OrgID:  orgID,
					Limit:  5,
					Cursor: firstBatch.NextCursor,
				}
			},
			check: func(t *testing.T, res *folders.PaginatedFetchRes, err error) {
				expected, _ := folders.FetchAllFoldersByOrgID(orgID)
				assert.NoError(t, err)
				assert.Equal(t, expected[5:10], res.Folders)
			},
		},
		{
			name: "Fetch Near End of Folder List",
			setup: func() *folders.PaginatedFetchReq {
				expected, _ := folders.FetchAllFoldersByOrgID(orgID)
				nextCursor := folders.EncodeNextCursor(len(expected) - 3)
				return &folders.PaginatedFetchReq{
					OrgID:  orgID,
					Limit:  5,
					Cursor: nextCursor,
				}
			},
			check: func(t *testing.T, res *folders.PaginatedFetchRes, err error) {
				expected, _ := folders.FetchAllFoldersByOrgID(orgID)
				assert.NoError(t, err)
				assert.Equal(t, len(expected[len(expected)-3:]), len(res.Folders))
				assert.Equal(t, expected[len(expected)-3:], res.Folders)
			},
		},
	}

	// Loops over the test cases in the tests array and runs it as subtests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *folders.PaginatedFetchReq
			if tt.setup != nil {
				req = tt.setup()
			} else {
				req = tt.req
			}

			res, err := folders.GetAllFoldersPaginated(req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, res, err)
				}
			}
		})
	}
}

func Test_EncodeNextCursor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		original int
		wantErr  bool
		check    func(t *testing.T, decoded int, err error)
	}{
		{
			name:     "Encode and Decode Valid Cursor",
			original: 5,
			check: func(t *testing.T, decoded int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 5, decoded)
			},
		},
		{
			name:  "Decode Empty Cursor",
			input: "",
			check: func(t *testing.T, decoded int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 0, decoded)
			},
		},
		{
			name:    "Decode Invalid Cursor Error",
			input:   "invalid_cursor",
			wantErr: true,
			check: func(t *testing.T, decoded int, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "illegal base64 data")
			},
		},
		{
			name:    "Decode Invalid Base64 Error",
			input:   "ThisIsNotBase64!",
			wantErr: true,
			check: func(t *testing.T, decoded int, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "illegal base64 data")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded int
			var err error

			if tt.original != 0 {
				encoded := folders.EncodeNextCursor(tt.original)
				decoded, err = folders.DecodeNextCursor(encoded)
			} else {
				decoded, err = folders.DecodeNextCursor(tt.input)
			}

			if tt.check != nil {
				tt.check(t, decoded, err)
			}
		})
	}
}
