package ze

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/rdb"
)

// Application request that holds DB connection, transaction, checker,
// transaction queries, request start time, and logs
type Request struct {
	Task
	DB      *sql.DB
	DBTx    *sql.Tx
	Checker rdb.QueryResultChecker
	Status  int
	start   DateTime
	logs    []string
	txSteps []rdb.Query
}

// Contains the Action and Item of the task
type Task struct {
	Action string
	Item   string
}

// Combine the "Action-Item" of CoreTask
func (t Task) FullAction() string {
	return fmt.Sprintf("%s-%s", t.Action, t.Item)
}

// Creates a new Request
func NewRequest(name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	rq := &Request{
		DB:     dbConn,
		Status: OK200,
		start:  clock.DateTimeNow(),
		logs:   make([]string, 0),
	}
	if dbConn == nil {
		rq.Status = Err500
		return rq, errNoDBConnection
	}
	return rq, nil
}

// Combines the logs with newline
func (rq Request) Output() string {
	return strings.Join(rq.logs, "\n")
}

// Add log to request
func (rq *Request) AddLog(message string) {
	rq.logs = append(rq.logs, nowLog(message))
}

// Add formatted log to request
func (rq *Request) AddFmtLog(format string, args ...any) {
	rq.AddLog(fmt.Sprintf(format, args...))
}

// Add duration log to request
func (rq *Request) AddDurationLog(start time.Time) {
	rq.AddFmtLog("Time: %v", time.Since(start))
}

// Add error log to request
func (rq *Request) AddErrorLog(err error) {
	rq.AddFmtLog("Error: %s", err.Error())
}

// Add transaction step to request
func (rq *Request) AddTxStep(q rdb.Query) {
	rq.txSteps = append(rq.txSteps, q)
}

// Start database transaction
func (rq *Request) StartTransaction(numSteps int) error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		rq.Status = Err500
		return errNoDBConnection
	}
	dbtx, err := rq.DB.Begin()
	if err != nil {
		rq.AddLog("Failed to start transaction")
		rq.Status = Err500
		return err
	}
	rq.DBTx = dbtx
	rq.txSteps = make([]rdb.Query, 0)
	rq.Checker = rdb.AssertNothing // default result checker
	return nil
}

// Commit database transaction
func (rq *Request) CommitTransaction() error {
	if rq.DB == nil {
		rq.AddLog("No DB connection")
		rq.Status = Err500
		return errNoDBConnection
	}
	if rq.DBTx == nil {
		rq.AddLog("No DB transaction")
		rq.Status = Err500
		return errNoDBTx
	}
	err := rq.DBTx.Commit()
	if err != nil {
		for i, q := range rq.txSteps {
			rq.AddFmtLog("Query %d: %s", i+1, rdb.QueryString(q))
		}
		rq.AddLog("Failed to commit transaction")
		rq.Status = Err500
		return fmt.Errorf("dbtx commit error: %w", err)
	}
	return nil
}

// Creates a message log with current datetime as prefix
func nowLog(message string) string {
	return fmt.Sprintf("%s | %s", clock.DateTimeNow(), message)
}
