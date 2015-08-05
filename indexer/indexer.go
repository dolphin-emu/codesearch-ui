/*
 * Copyright 2015 Dolphin Emulator project. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"flag"
	"github.com/google/codesearch/index"
	"golang.org/x/net/context"
	esearch "kythe.io/kythe/go/services/search"
	esrchsrv "kythe.io/kythe/go/serving/search"
	xsrv "kythe.io/kythe/go/serving/xrefs"
	"kythe.io/kythe/go/storage/leveldb"
	"kythe.io/kythe/go/storage/table"
	"kythe.io/kythe/go/util/flagutil"
	srvpb "kythe.io/kythe/proto/serving_proto"
	spb "kythe.io/kythe/proto/storage_proto"
	"log"
)

var (
	outPath      = flag.String("out", "", "Output index file")
	servingTable = flag.String("serving_table", "", "LevelDB serving table")
)

func init() {
	flag.Usage = flagutil.SimpleUsage("Indexes Kythe LevelDB serving tables into a codesearch index",
		"--serving_table path --out out-file")
}

func filesInTable(ctx context.Context, esr esearch.Service) []string {
	req := &spb.SearchRequest{
		Fact: []*spb.SearchRequest_Fact{
			&spb.SearchRequest_Fact{
				Name:   "/kythe/node/kind",
				Value:  []byte("file"),
				Prefix: false,
			},
		},
	}
	repl, err := esr.Search(ctx, req)
	if err != nil {
		log.Fatalf("File search in serving tables failed: %v", err)
	}
	return repl.Ticket
}

func main() {
	flag.Parse()
	if *servingTable == "" {
		flagutil.UsageError("missing --serving_table")
	} else if *outPath == "" {
		flagutil.UsageError("missing --out")
	}

	ctx := context.Background()
	db, err := leveldb.Open(*servingTable, nil)
	if err != nil {
		log.Fatalf("Error opening db at %q: %v", *servingTable, err)
	}
	defer db.Close()
	tbl := &table.KVProto{db}
	esr := &esrchsrv.Table{&table.KVInverted{db}}

	ix := index.Create(*outPath)
	for _, ticket := range filesInTable(ctx, esr) {
		var fd srvpb.FileDecorations
		if err := tbl.Lookup(ctx, xsrv.DecorationsKey(ticket), &fd); err != nil {
			log.Fatalf("Error looking up decoration for %q: %v", ticket, err)
		}
		ix.Add(ticket, bytes.NewReader(fd.SourceText))
	}
	ix.Flush()
}
