package scraper

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/status-im/keycard-go/hexutils"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/rpc"
	"ramdeuter.org/solscraper/query"
)

func ScrapeData(q query.Query) []string {
	solanaRpcClient := rpc.NewClient("https://solana-api.projectserum.com")
	txList, err := solanaRpcClient.GetSignaturesForAddress(context.TODO(), solana.MustPublicKeyFromBase58(q.Program),
		&rpc.GetSignaturesForAddressOpts{
			Limit: 100,
		})

	if err != nil {
		//log.Warnf("unable to retrieve confirmed transaction signatures for account: %s", err.Error())
		//return err
	}

	fmt.Printf("processing %d  transactions", len(txList))

	resp := make([]string, 0)

	for i := 0; i < len(txList); i++ {
		tx := txList[i]

		confirmedTx, err := solanaRpcClient.GetConfirmedTransaction(context.TODO(), tx.Signature)
		if err != nil || confirmedTx.Meta == nil ||
			confirmedTx.Transaction == nil || confirmedTx.Transaction.Message.AccountKeys == nil {
			//log.Errorf("unable to get confirmed transaction with signature %q: %v", tx.Signature, err)
			//return false, err
		} else if confirmedTx.Meta.Err != nil {
			//return true, err
		}

		processTx := true
		for _, filter := range q.Filters {
			value := parseValue(filter.Value, confirmedTx)
			if !match(value, filter.MatchType, filter.MatchValue) {
				processTx = false
				break
			}
		}

		if processTx {
			val := ""
			for _, exportFunc := range q.Export {
				val += exportFunc.Name + ":" + parseValue(exportFunc.Function, confirmedTx) + ", "
			}
			resp = append(resp, val)
		}

	}
	return resp
}

func parseValue(s string, tx rpc.TransactionWithMeta) string {
	values := strings.Split(s, ".")
	if values[0] == "instructions" {
		instIndex, _ := strconv.Atoi(values[1])
		inst := tx.Transaction.Message.Instructions[instIndex]
		switch values[2] {
		case "data":
			return strings.ToLower(hexutils.BytesToHex(inst.Data))
		case "accounts":
			acctIndex, _ := strconv.Atoi(values[3])
			acccIndexInTx := inst.Accounts[acctIndex]
			return tx.Transaction.Message.AccountKeys[acccIndexInTx].String()
		}
	}
	return ""
}

func match(value, matchType, matchValue string) bool {
	switch matchType {
	case "contains":
		if strings.Contains(value, matchValue) {
			return true
		}
		break
	case "eq":
		if strings.EqualFold(value, matchValue) {
			return true
		}
		break
	case "startsWith":
		if strings.HasPrefix(value, matchValue) {
			return true
		}
		break
	case "endsWith":
		if strings.HasSuffix(value, matchValue) {
			return true
		}
		break
	case "regex":
		if match, _ := regexp.MatchString(matchValue, value); match {
			return true
		}
		break
	}
	return false
}
