// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Code generated by operationgen. DO NOT EDIT.

package operation

import (
	"context"
	"errors"

	"github.com/dollarkillerx/mongo/mongo-driver/event"
	"github.com/dollarkillerx/mongo/mongo-driver/mongo/writeconcern"
	"github.com/dollarkillerx/mongo/mongo-driver/x/bsonx/bsoncore"
	"github.com/dollarkillerx/mongo/mongo-driver/x/mongo/driver"
	"github.com/dollarkillerx/mongo/mongo-driver/x/mongo/driver/description"
	"github.com/dollarkillerx/mongo/mongo-driver/x/mongo/driver/session"
)

// AbortTransaction performs an abortTransaction operation.
type AbortTransaction struct {
	recoveryToken bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
}

// NewAbortTransaction constructs and returns a new AbortTransaction.
func NewAbortTransaction() *AbortTransaction {
	return &AbortTransaction{}
}

func (at *AbortTransaction) processResponse(response bsoncore.Document, srvr driver.Server, desc description.Server) error {
	var err error
	return err
}

// Execute runs this operations and returns an error if the operaiton did not execute successfully.
func (at *AbortTransaction) Execute(ctx context.Context) error {
	if at.deployment == nil {
		return errors.New("the AbortTransaction operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         at.command,
		ProcessResponseFn: at.processResponse,
		RetryMode:         at.retry,
		Type:              driver.Write,
		Client:            at.session,
		Clock:             at.clock,
		CommandMonitor:    at.monitor,
		Database:          at.database,
		Deployment:        at.deployment,
		Selector:          at.selector,
		WriteConcern:      at.writeConcern,
	}.Execute(ctx, nil)

}

func (at *AbortTransaction) command(dst []byte, desc description.SelectedServer) ([]byte, error) {

	dst = bsoncore.AppendInt32Element(dst, "abortTransaction", 1)
	if at.recoveryToken != nil {
		dst = bsoncore.AppendDocumentElement(dst, "recoveryToken", at.recoveryToken)
	}
	return dst, nil
}

// RecoveryToken sets the recovery token to use when committing or aborting a sharded transaction.
func (at *AbortTransaction) RecoveryToken(recoveryToken bsoncore.Document) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.recoveryToken = recoveryToken
	return at
}

// Session sets the session for this operation.
func (at *AbortTransaction) Session(session *session.Client) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.session = session
	return at
}

// ClusterClock sets the cluster clock for this operation.
func (at *AbortTransaction) ClusterClock(clock *session.ClusterClock) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.clock = clock
	return at
}

// Collection sets the collection that this command will run against.
func (at *AbortTransaction) Collection(collection string) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.collection = collection
	return at
}

// CommandMonitor sets the monitor to use for APM events.
func (at *AbortTransaction) CommandMonitor(monitor *event.CommandMonitor) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.monitor = monitor
	return at
}

// Database sets the database to run this operation against.
func (at *AbortTransaction) Database(database string) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.database = database
	return at
}

// Deployment sets the deployment to use for this operation.
func (at *AbortTransaction) Deployment(deployment driver.Deployment) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.deployment = deployment
	return at
}

// ServerSelector sets the selector used to retrieve a server.
func (at *AbortTransaction) ServerSelector(selector description.ServerSelector) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.selector = selector
	return at
}

// WriteConcern sets the write concern for this operation.
func (at *AbortTransaction) WriteConcern(writeConcern *writeconcern.WriteConcern) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.writeConcern = writeConcern
	return at
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (at *AbortTransaction) Retry(retry driver.RetryMode) *AbortTransaction {
	if at == nil {
		at = new(AbortTransaction)
	}

	at.retry = &retry
	return at
}
