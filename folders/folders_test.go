package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllFolders(t *testing.T) {
	// An array of test cases is defined.
	tests := []struct {
		name        string
		orgID       uuid.UUID
		expectError bool
		expectEmpty bool
	}{
		// Test cases:.
		{
			name:        "returns error when orgID is nil",
			orgID:       uuid.Nil,
			expectError: true,
			expectEmpty: true,
		},
		{
			name:        "returns no folders when orgID does not match",
			orgID:       uuid.Must(uuid.NewV4()),
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "returns no folders when orgID is valid but no folders exist",
			orgID:       uuid.Must(uuid.NewV4()),
			expectError: false,
			expectEmpty: true,
		},
		{
			name:        "returns all folders when orgID matches",
			orgID:       uuid.FromStringOrNil(folders.DefaultOrgID),
			expectError: false,
			expectEmpty: false,
		},
		{
			name:        "returns multiple folders when orgID is valid",
			orgID:       uuid.FromStringOrNil(folders.DefaultOrgID),
			expectError: false,
			expectEmpty: false,
		},
	}

	// The test function loops over each test case in the tests array and runs it as subtests.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := folders.GetAllFolders(&folders.FetchFolderRequest{OrgID: tt.orgID})

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tt.expectEmpty {
					assert.Empty(t, res.Folders)
				} else {
					assert.NotEmpty(t, res.Folders)
					for _, folder := range res.Folders {
						assert.Equal(t, tt.orgID, folder.OrgId)
					}
				}
			}
		})
	}
}
