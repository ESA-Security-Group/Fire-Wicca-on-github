// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

type testRequestlog struct {
	url string
}

func (t testRequestlog) Run() error {
	tok, err := get(address + t.url)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	time.Sleep(time.Minute)
	return t.readLogEntries(tok)
}

func (t testRequestlog) readLogEntries(tok string) error {
	iter := logadminClient.Entries(context.Background(),
		logadmin.Filter(fmt.Sprintf(`timestamp >= %q`, tok)),
		logadmin.Filter(fmt.Sprintf(`httpRequest.requestUrl = "%s%s"`, t.url, tok)),
	)
	_, err := iter.Next()
	if err == iterator.Done {
		return fmt.Errorf("no entry found for request log that matches %s", tok)
	}
	if err != nil {
		return err
	}
	return nil
}

func (t testRequestlog) String() string {
	return fmt.Sprintf("%T:%s", t, t.url)
}
