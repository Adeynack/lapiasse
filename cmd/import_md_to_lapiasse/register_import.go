package main

import (
	"context"
	"fmt"

	"adeynack.net/lapiasse/cmd/import_md_to_lapiasse/moneydance"
	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/applog"
	"github.com/samber/lo"
)

type registerImport struct {
	apiClient *api.ClientWithResponses
	book      api.Book
	md        *moneydance.Export

	mdAccountsByParentId map[string][]*moneydance.Account
}

func (ri *registerImport) run(ctx context.Context) error {
	ri.mdAccountsByParentId = lo.GroupBy(ri.md.AllItems.Accounts, func(account *moneydance.Account) string {
		return account.ParentId
	})

	mdRootAccounts := ri.mdAccountsByParentId[""]
	if len(mdRootAccounts) != 1 {
		return fmt.Errorf("expecting exactly one root account, got %d", len(mdRootAccounts))
	}
	mdRootAccountId := mdRootAccounts[0].Id

	mdFirstLevelAccountsPerType := lo.GroupBy(ri.mdAccountsByParentId[mdRootAccountId], func(a *moneydance.Account) moneydance.AccountType {
		return a.Type
	})

	// Import income categories first ("i")
	applog.Info(ctx, "Importing income categories")
	if err := ri.importRegisterBatch(ctx, mdFirstLevelAccountsPerType[moneydance.AccountTypeIncome]); err != nil {
		return fmt.Errorf("importing income categories: %w", err)
	}
	delete(mdFirstLevelAccountsPerType, moneydance.AccountTypeIncome)

	// Then import expense categories ("e")
	applog.Info(ctx, "Importing expense categories")
	if err := ri.importRegisterBatch(ctx, mdFirstLevelAccountsPerType[moneydance.AccountTypeExpense]); err != nil {
		return fmt.Errorf("importing expense categories: %w", err)
	}
	delete(mdFirstLevelAccountsPerType, moneydance.AccountTypeExpense)

	// Finally import the other accounts
	for acctType, mdAccounts := range mdFirstLevelAccountsPerType {
		applog.Info(ctx, fmt.Sprintf("Importing %s accounts", acctType))
		if err := ri.importRegisterBatch(ctx, mdAccounts); err != nil {
			return fmt.Errorf("importing %s accounts: %w", acctType, err)
		}
	}

	return nil
}

func (ri *registerImport) importRegisterBatch(ctx context.Context, mdAccounts []*moneydance.Account) error {
	for _, mdAccount := range mdAccounts {
		if err := ri.importRegisterRecursively(ctx, mdAccount, nil); err != nil {
			return err
		}
	}

	return nil
}

func (ri *registerImport) importRegisterRecursively(
	ctx context.Context,
	mdAccount *moneydance.Account,
	parentRegister *api.Register,
) error {
	_ = ctx
	_ = mdAccount
	_ = parentRegister

	return nil
}
