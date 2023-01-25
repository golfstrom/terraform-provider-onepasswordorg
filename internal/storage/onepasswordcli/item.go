package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r Repository) CreateItem(ctx context.Context, item model.Item) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	// Create Fields
	cmdArgs.ItemArg().ProvisionArg().EditFieldFlag("title", item.Title).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}

func (r Repository) GetItemByID(ctx context.Context, id string) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().GetArg().RawStrArg(id).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}
func (r Repository) GetItemByTitle(ctx context.Context, vaultID string, title string) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().GetArg().RawStrArg(title).FormatJSONFlag().VaultFlag(vaultID)

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	fmt.Println(stdout)
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opItem{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotItem := mapOpToModelItem(ou)

	return &gotItem, nil
}

func (r Repository) EnsureItem(ctx context.Context, item model.Item) (*model.Item, error) {
	cmdArgs := &onePasswordCliCmd{}
	// TODO: update fields
	cmdArgs.ItemArg().EditArg().RawStrArg(item.ID).EditFieldFlag("title", item.Title)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &item, nil
}

func (r Repository) DeleteItem(ctx context.Context, id string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ItemArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

type opVaultItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type opItemField struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
	Label   string `json:"label"`
	Value   string `json:"value"`
}

type opItem struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Vault  opVaultItem   `json:"vault"`
	Fields []opItemField `json:"fields"`
}

func mapOpToModelItem(u opItem) model.Item {
	return model.Item{
		ID:      u.ID,
		Title:   u.Title,
		VaultID: u.Vault.ID,
		Fields:  mapOpToModelItemFields(u.Fields),
	}
}

func mapOpToModelItemFields(opFields []opItemField) []model.Field {
	fields := []model.Field{}

	for _, opField := range opFields {
		field := model.Field{
			ID:      opField.ID,
			Type:    opField.Type,
			Purpose: opField.Purpose,
			Label:   opField.Label,
			Value:   opField.Value,
		}
		fields = append(fields, field)
	}

	return fields
}
