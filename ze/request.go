package ze

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/rdb"
)

// Contains the Action and Target of the Task
type Task struct {
	Action string
	Target string
}

// Application request that holds DB connection, transaction, checker,
// transaction queries, request start time, and logs
type Request struct {
	Task
	Name    string
	Params  dict.Object
	DB      *sql.DB
	DBTx    *sql.Tx
	Checker rdb.ResultChecker
	Status  int
	// Private fields
	start   DateTime
	txSteps []rdb.Query
	// Logs
	mu   sync.RWMutex
	logs []string
}

// Create new Request
func NewRequest(name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	rq := newRequest(name)
	if dbConn == nil {
		rq.Status = Err500
		return rq, errNoDBConnection
	}
	rq.DB = dbConn
	return rq, nil
}

// Create new Request at custom db
func NewRequestAt(key, name string, args ...any) (*Request, error) {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}
	rq := newRequest(name)
	conn, ok := dbConnMap[key]
	if !ok || conn == nil {
		rq.Status = Err500
		return rq, errNoDBConnection
	}
	rq.DB = conn
	return rq, nil
}

// Create a new request object
func newRequest(name string) *Request {
	return &Request{
		Name:   name,
		Params: make(dict.Object),
		Status: OK200,
		start:  clock.DateTimeNow(),
		logs:   make([]string, 0),
	}
}

// Create a subrequest for concurrent tasks
func (rq *Request) SubRequest() *Request {
	return &Request{
		Task:   rq.Task,
		Params: rq.Params,
		DB:     rq.DB,
		Status: OK200,
		logs:   make([]string, 0),
	}
}

// Combine logs with newline
func (rq *Request) Output() string {
	return strings.Join(rq.logs, "\n")
}

// Concurrent-safe merging of logs
func (rq *Request) MergeLogs(srq *Request) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.logs = append(rq.logs, srq.logs...)
}

// Add log to request
func (rq *Request) AddLog(message string) {
	message = fmt.Sprintf("%s | %s", clock.DateTimeNow(), message)
	rq.logs = append(rq.logs, message)
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
	rq.Checker = rdb.AssertNothing // default checker
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

// Combine "Action-Target" of Task
func (t Task) FullName() string {
	target := t.Target
	if strings.HasSuffix(target, "-%s") {
		parts := str.CleanSplit(target, "-")
		target = strings.Join(parts[:len(parts)-1], "-")
	}
	return fmt.Sprintf("%s-%s", t.Action, target)
}
