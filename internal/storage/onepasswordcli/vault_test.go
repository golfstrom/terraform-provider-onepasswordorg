package onepasswordcli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli/onepasswordclimock"
)

func TestRepositoryCreateVault(t *testing.T) {
	tests := map[string]struct {
		vault    model.Vault
		mock     func(m *onepasswordclimock.OpCli)
		expVault *model.Vault
		expErr   bool
	}{
		"Creating a vault correctly, should return the data with the ID.": {
			vault: model.Vault{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get vault test-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("vault doesn't exist"))

				expCmd = `create vault test-00 --description Test00`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00","createdAt":"2022-03-14T07:48:26.179385832+01:00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expVault: &model.Vault{
				ID:          "1234567890",
				Name:        "test-00",
				Description: "Test00",
			},
		},

		"Creating a vault that already exists, should  fail.": {
			vault: model.Vault{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get vault test-00`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00","createdAt":"2022-03-14T07:48:26.179385832+01:00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the create op CLI action, should fail.": {
			vault: model.Vault{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get vault test-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("vault doesn't exist"))

				expCmd = `create vault test-00 --description Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotVault, err := repo.CreateVault(context.TODO(), test.vault)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expVault, gotVault)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetVaultByID(t *testing.T) {
	tests := map[string]struct {
		id       string
		mock     func(m *onepasswordclimock.OpCli)
		expVault *model.Vault
		expErr   bool
	}{
		"Getting a vault correctly, should return the vault data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get vault test-id`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expVault: &model.Vault{
				ID:          "1234567890",
				Name:        "test-00",
				Description: "Test00",
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get vault test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotVault, err := repo.GetVaultByID(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expVault, gotVault)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryEnsureVault(t *testing.T) {
	tests := map[string]struct {
		vault    model.Vault
		mock     func(m *onepasswordclimock.OpCli)
		expVault *model.Vault
		expErr   bool
	}{
		"Updating a vault correctly, should update the user data.": {
			vault: model.Vault{ID: "test-id", Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `edit vault test-id --description Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
			expVault: &model.Vault{ID: "test-id", Name: "test-00", Description: "Test00"},
		},

		"Having an error while calling the op CLI, should fail.": {
			vault: model.Vault{ID: "test-id", Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `edit vault test-id --description Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotVault, err := repo.EnsureVault(context.TODO(), test.vault)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expVault, gotVault)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteVault(t *testing.T) {
	tests := map[string]struct {
		id     string
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Delete a vault correctly, should return the vault data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `delete vault test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `delete vault test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			err = repo.DeleteVault(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
