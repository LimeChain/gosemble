package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionOutcomeCommit(t *testing.T) {
	value := sc.U8(5)
	assert.Equal(t, TransactionOutcome(sc.NewVaryingData(TransactionOutcomeCommit, value)), NewTransactionOutcomeCommit(value))
}

func Test_NewTransactionOutcomeRollback(t *testing.T) {
	value := sc.U8(5)
	assert.Equal(t, TransactionOutcome(sc.NewVaryingData(TransactionOutcomeRollback, value)), NewTransactionOutcomeRollback(value))
}
