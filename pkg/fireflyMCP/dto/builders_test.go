package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccountBuilder(t *testing.T) {
	t.Run("builds account with all fields", func(t *testing.T) {
		notes := "test notes"
		account := NewAccountBuilder().
			WithId("123").
			WithActive(true).
			WithName("Test Account").
			WithNotes(&notes).
			WithType("asset").
			Build()

		assert.Equal(t, "123", account.GetId())
		assert.True(t, account.GetActive())
		assert.Equal(t, "Test Account", account.GetName())
		assert.NotNil(t, account.GetNotes())
		assert.Equal(t, "test notes", *account.GetNotes())
		assert.Equal(t, "asset", account.GetType())
	})

	t.Run("builds account with nil notes", func(t *testing.T) {
		account := NewAccountBuilder().
			WithId("456").
			WithActive(false).
			WithName("Another Account").
			WithNotes(nil).
			WithType("expense").
			Build()

		assert.Equal(t, "456", account.GetId())
		assert.False(t, account.GetActive())
		assert.Equal(t, "Another Account", account.GetName())
		assert.Nil(t, account.GetNotes())
		assert.Equal(t, "expense", account.GetType())
	})

	t.Run("immutability test", func(t *testing.T) {
		notes := "original"
		account := NewAccountBuilder().
			WithId("789").
			WithNotes(&notes).
			Build()

		// Changing the original variable shouldn't affect the account
		notes = "modified"
		assert.Equal(t, "original", *account.GetNotes())
	})
}

func TestBudgetBuilder(t *testing.T) {
	t.Run("builds budget with all fields", func(t *testing.T) {
		spent := NewSpentBuilder().
			WithSum("100.00").
			WithCurrencyCode("USD").
			Build()

		budget := NewBudgetBuilder().
			WithId("budget-1").
			WithActive(true).
			WithName("Monthly Budget").
			WithNotes("Budget notes").
			WithSpent(spent).
			Build()

		assert.Equal(t, "budget-1", budget.GetId())
		assert.True(t, budget.GetActive())
		assert.Equal(t, "Monthly Budget", budget.GetName())
		assert.Equal(t, "Budget notes", budget.GetNotes())
		assert.Equal(t, "100.00", budget.GetSpent().GetSum())
		assert.Equal(t, "USD", budget.GetSpent().GetCurrencyCode())
	})
}

func TestCategoryBuilder(t *testing.T) {
	t.Run("builds category with all fields", func(t *testing.T) {
		category := NewCategoryBuilder().
			WithId("cat-1").
			WithName("Groceries").
			WithNotes("Food expenses").
			Build()

		assert.Equal(t, "cat-1", category.GetId())
		assert.Equal(t, "Groceries", category.GetName())
		assert.Equal(t, "Food expenses", category.GetNotes())
	})
}

func TestTransactionBuilder(t *testing.T) {
	t.Run("builds transaction with all fields", func(t *testing.T) {
		budgetId := "budget-1"
		budgetName := "Monthly"
		categoryId := "cat-1"
		categoryName := "Food"
		notes := "Grocery shopping"
		
		transaction := NewTransactionBuilder().
			WithId("tx-1").
			WithAmount("50.00").
			WithBillId(nil).
			WithBillName(nil).
			WithBudgetId(&budgetId).
			WithBudgetName(&budgetName).
			WithCategoryId(&categoryId).
			WithCategoryName(&categoryName).
			WithCurrencyCode("EUR").
			WithDate(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)).
			WithDescription("Supermarket purchase").
			WithDestinationId("acc-2").
			WithDestinationName("Savings").
			WithDestinationType("asset").
			WithNotes(&notes).
			WithReconciled(true).
			WithSourceId("acc-1").
			WithSourceName("Checking").
			WithTags([]string{"food", "essentials"}).
			WithType("withdrawal").
			Build()

		assert.Equal(t, "tx-1", transaction.GetId())
		assert.Equal(t, "50.00", transaction.GetAmount())
		assert.Nil(t, transaction.GetBillId())
		assert.Nil(t, transaction.GetBillName())
		assert.NotNil(t, transaction.GetBudgetId())
		assert.Equal(t, "budget-1", *transaction.GetBudgetId())
		assert.NotNil(t, transaction.GetBudgetName())
		assert.Equal(t, "Monthly", *transaction.GetBudgetName())
		assert.NotNil(t, transaction.GetCategoryId())
		assert.Equal(t, "cat-1", *transaction.GetCategoryId())
		assert.NotNil(t, transaction.GetCategoryName())
		assert.Equal(t, "Food", *transaction.GetCategoryName())
		assert.Equal(t, "EUR", transaction.GetCurrencyCode())
		assert.Equal(t, time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), transaction.GetDate())
		assert.Equal(t, "Supermarket purchase", transaction.GetDescription())
		assert.Equal(t, "acc-2", transaction.GetDestinationId())
		assert.Equal(t, "Savings", transaction.GetDestinationName())
		assert.Equal(t, "asset", transaction.GetDestinationType())
		assert.NotNil(t, transaction.GetNotes())
		assert.Equal(t, "Grocery shopping", *transaction.GetNotes())
		assert.True(t, transaction.GetReconciled())
		assert.Equal(t, "acc-1", transaction.GetSourceId())
		assert.Equal(t, "Checking", transaction.GetSourceName())
		assert.Equal(t, []string{"food", "essentials"}, transaction.GetTags())
		assert.Equal(t, "withdrawal", transaction.GetType())
	})

	t.Run("tags immutability", func(t *testing.T) {
		tags := []string{"tag1", "tag2"}
		transaction := NewTransactionBuilder().
			WithTags(tags).
			Build()

		// Modifying original tags shouldn't affect transaction
		tags[0] = "modified"
		retrievedTags := transaction.GetTags()
		assert.Equal(t, "tag1", retrievedTags[0])

		// Modifying retrieved tags shouldn't affect transaction
		retrievedTags[1] = "changed"
		assert.Equal(t, "tag2", transaction.GetTags()[1])
	})
}

func TestTransactionGroupBuilder(t *testing.T) {
	t.Run("builds transaction group with transactions", func(t *testing.T) {
		tx1 := NewTransactionBuilder().
			WithId("tx-1").
			WithDescription("Transaction 1").
			Build()
		
		tx2 := NewTransactionBuilder().
			WithId("tx-2").
			WithDescription("Transaction 2").
			Build()

		group := NewTransactionGroupBuilder().
			WithId("group-1").
			WithGroupTitle("Split Transaction").
			WithTransactions([]Transaction{tx1, tx2}).
			Build()

		assert.Equal(t, "group-1", group.GetId())
		assert.Equal(t, "Split Transaction", group.GetGroupTitle())
		assert.Len(t, group.GetTransactions(), 2)
		assert.Equal(t, "tx-1", group.GetTransactions()[0].GetId())
		assert.Equal(t, "tx-2", group.GetTransactions()[1].GetId())
	})

	t.Run("add transaction individually", func(t *testing.T) {
		tx1 := NewTransactionBuilder().WithId("tx-1").Build()
		tx2 := NewTransactionBuilder().WithId("tx-2").Build()

		group := NewTransactionGroupBuilder().
			WithId("group-1").
			WithGroupTitle("Group").
			AddTransaction(tx1).
			AddTransaction(tx2).
			Build()

		assert.Len(t, group.GetTransactions(), 2)
	})

	t.Run("transactions immutability", func(t *testing.T) {
		tx := NewTransactionBuilder().WithId("tx-1").Build()
		transactions := []Transaction{tx}
		
		group := NewTransactionGroupBuilder().
			WithTransactions(transactions).
			Build()

		// Modifying original array shouldn't affect group
		transactions[0] = NewTransactionBuilder().WithId("modified").Build()
		assert.Equal(t, "tx-1", group.GetTransactions()[0].GetId())
	})
}

func TestPaginationBuilder(t *testing.T) {
	t.Run("builds pagination with all fields", func(t *testing.T) {
		pagination := NewPaginationBuilder().
			WithCount(10).
			WithTotal(100).
			WithCurrentPage(2).
			WithPerPage(10).
			WithTotalPages(10).
			Build()

		assert.Equal(t, 10, pagination.GetCount())
		assert.Equal(t, 100, pagination.GetTotal())
		assert.Equal(t, 2, pagination.GetCurrentPage())
		assert.Equal(t, 10, pagination.GetPerPage())
		assert.Equal(t, 10, pagination.GetTotalPages())
	})
}

func TestTagBuilder(t *testing.T) {
	t.Run("builds tag with description", func(t *testing.T) {
		desc := "Important tag"
		tag := NewTagBuilder().
			WithId("tag-1").
			WithTag("important").
			WithDescription(&desc).
			Build()

		assert.Equal(t, "tag-1", tag.GetId())
		assert.Equal(t, "important", tag.GetTag())
		assert.NotNil(t, tag.GetDescription())
		assert.Equal(t, "Important tag", *tag.GetDescription())
	})

	t.Run("builds tag without description", func(t *testing.T) {
		tag := NewTagBuilder().
			WithId("tag-2").
			WithTag("normal").
			WithDescription(nil).
			Build()

		assert.Equal(t, "tag-2", tag.GetId())
		assert.Equal(t, "normal", tag.GetTag())
		assert.Nil(t, tag.GetDescription())
	})
}

func TestBillBuilder(t *testing.T) {
	t.Run("builds bill with all fields", func(t *testing.T) {
		notes := "Monthly bill"
		nextMatch := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		paidDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		groupId := "group-1"
		journalId := "journal-1"
		
		paidDates := []PaidDate{
			{
				Date:                 &paidDate,
				TransactionGroupId:   &groupId,
				TransactionJournalId: &journalId,
			},
		}

		bill := NewBillBuilder().
			WithId("bill-1").
			WithActive(true).
			WithName("Internet").
			WithAmountMin("49.99").
			WithAmountMax("49.99").
			WithDate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
			WithRepeatFreq("monthly").
			WithSkip(0).
			WithCurrencyCode("USD").
			WithNotes(&notes).
			WithNextExpectedMatch(&nextMatch).
			WithPaidDates(paidDates).
			Build()

		assert.Equal(t, "bill-1", bill.GetId())
		assert.True(t, bill.GetActive())
		assert.Equal(t, "Internet", bill.GetName())
		assert.Equal(t, "49.99", bill.GetAmountMin())
		assert.Equal(t, "49.99", bill.GetAmountMax())
		assert.Equal(t, "monthly", bill.GetRepeatFreq())
		assert.Equal(t, 0, bill.GetSkip())
		assert.Equal(t, "USD", bill.GetCurrencyCode())
		assert.NotNil(t, bill.GetNotes())
		assert.Equal(t, "Monthly bill", *bill.GetNotes())
		assert.NotNil(t, bill.GetNextExpectedMatch())
		assert.Len(t, bill.GetPaidDates(), 1)
	})

	t.Run("add paid date individually", func(t *testing.T) {
		bill := NewBillBuilder().
			WithId("bill-1").
			AddPaidDate(PaidDate{}).
			AddPaidDate(PaidDate{}).
			Build()

		assert.Len(t, bill.GetPaidDates(), 2)
	})
}

func TestRecurrenceBuilder(t *testing.T) {
	t.Run("builds recurrence with all fields", func(t *testing.T) {
		firstDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		latestDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		repeatUntil := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		nrRepetitions := 12
		notes := "Monthly recurrence"

		repetitions := []RecurrenceRepetition{
			{
				Id:     "rep-1",
				Type:   "monthly",
				Moment: "1",
			},
		}

		transactions := []RecurrenceTransaction{
			{
				Id:          "rec-tx-1",
				Description: "Recurring payment",
				Amount:      "100.00",
			},
		}

		recurrence := NewRecurrenceBuilder().
			WithId("rec-1").
			WithType("withdrawal").
			WithTitle("Monthly Payment").
			WithDescription("Recurring monthly payment").
			WithFirstDate(firstDate).
			WithLatestDate(&latestDate).
			WithRepeatUntil(&repeatUntil).
			WithNrOfRepetitions(&nrRepetitions).
			WithApplyRules(true).
			WithActive(true).
			WithNotes(&notes).
			WithRepetitions(repetitions).
			WithTransactions(transactions).
			Build()

		assert.Equal(t, "rec-1", recurrence.GetId())
		assert.Equal(t, "withdrawal", recurrence.GetType())
		assert.Equal(t, "Monthly Payment", recurrence.GetTitle())
		assert.Equal(t, "Recurring monthly payment", recurrence.GetDescription())
		assert.Equal(t, firstDate, recurrence.GetFirstDate())
		assert.NotNil(t, recurrence.GetLatestDate())
		assert.Equal(t, latestDate, *recurrence.GetLatestDate())
		assert.NotNil(t, recurrence.GetRepeatUntil())
		assert.Equal(t, repeatUntil, *recurrence.GetRepeatUntil())
		assert.NotNil(t, recurrence.GetNrOfRepetitions())
		assert.Equal(t, 12, *recurrence.GetNrOfRepetitions())
		assert.True(t, recurrence.GetApplyRules())
		assert.True(t, recurrence.GetActive())
		assert.NotNil(t, recurrence.GetNotes())
		assert.Equal(t, "Monthly recurrence", *recurrence.GetNotes())
		assert.Len(t, recurrence.GetRepetitions(), 1)
		assert.Len(t, recurrence.GetTransactions(), 1)
	})
}

func TestSpentBuilder(t *testing.T) {
	t.Run("builds spent with all fields", func(t *testing.T) {
		spent := NewSpentBuilder().
			WithSum("250.50").
			WithCurrencyCode("EUR").
			Build()

		assert.Equal(t, "250.50", spent.GetSum())
		assert.Equal(t, "EUR", spent.GetCurrencyCode())
	})

	t.Run("helper function NewSpentFromValues", func(t *testing.T) {
		spent := NewSpentFromValues("100.00", "USD")

		assert.Equal(t, "100.00", spent.GetSum())
		assert.Equal(t, "USD", spent.GetCurrencyCode())
	})
}

func TestNewPaginationFromValues(t *testing.T) {
	pagination := NewPaginationFromValues(5, 50, 1, 10, 5)

	assert.Equal(t, 5, pagination.GetCount())
	assert.Equal(t, 50, pagination.GetTotal())
	assert.Equal(t, 1, pagination.GetCurrentPage())
	assert.Equal(t, 10, pagination.GetPerPage())
	assert.Equal(t, 5, pagination.GetTotalPages())
}