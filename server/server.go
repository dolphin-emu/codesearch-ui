/*
 * Copyright 2015 Google Inc. All rights reserved.
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

// Binary http_server exposes HTTP interfaces for the search, xrefs, and
// filetree services backed by either a combined serving table or a bare
// GraphStore.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dolphin-emu/codesearch-ui/server/codesearch"

	"kythe.io/kythe/go/services/filetree"
	"kythe.io/kythe/go/services/graphstore"
	esearch "kythe.io/kythe/go/services/search"
	"kythe.io/kythe/go/services/xrefs"
	ftsrv "kythe.io/kythe/go/serving/filetree"
	esrchsrv "kythe.io/kythe/go/serving/search"
	xsrv "kythe.io/kythe/go/serving/xrefs"
	"kythe.io/kythe/go/storage/gsutil"
	"kythe.io/kythe/go/storage/leveldb"
	"kythe.io/kythe/go/storage/table"
	xstore "kythe.io/kythe/go/storage/xrefs"
	"kythe.io/kythe/go/util/flagutil"

	"golang.org/x/net/context"

	_ "kythe.io/kythe/go/services/graphstore/proxy"
	_ "kythe.io/kythe/go/storage/leveldb"
)

var (
	gs           graphstore.Service
	servingTable = flag.String("serving_table", "", "LevelDB serving table")
	csIndex      = flag.String("csearch_index", "", "Codesearch index file")

	httpListeningAddr = flag.String("listen", "localhost:8080", "Listening address for HTTP server")
	publicResources   = flag.String("public_resources", "", "Path to directory of static resources to serve")
)

func init() {
	gsutil.Flag(&gs, "graphstore", "GraphStore to serve xrefs")
	flag.Usage = flagutil.SimpleUsage("Exposes HTTP interfaces for the search, xrefs, and filetree services",
		"(--graphstore spec | --serving_table path) [--csearch_index file] [--listen addr] [--public_resources dir]")
}

func main() {
	flag.Parse()
	if *servingTable == "" && gs == nil {
		flagutil.UsageError("missing either --serving_table or --graphstore")
	} else if *httpListeningAddr == "" {
		flagutil.UsageError("missing --listen argument")
	} else if *servingTable != "" && gs != nil {
		flagutil.UsageError("--serving_table and --graphstore are mutually exclusive")
	}

	var (
		xs  xrefs.Service
		ft  filetree.Service
		esr esearch.Service
		csr codesearch.Service
	)

	ctx := context.Background()
	if *servingTable != "" {
		db, err := leveldb.Open(*servingTable, nil)
		if err != nil {
			log.Fatalf("Error opening db at %q: %v", *servingTable, err)
		}
		defer db.Close()
		tbl := &table.KVProto{db}
		xs = &xsrv.Table{tbl}
		ft = &ftsrv.Table{tbl}
		esr = &esrchsrv.Table{&table.KVInverted{db}}

		if *csIndex != "" {
			csr = codesearch.New(*csIndex, xs)
		}
	} else {
		log.Println("WARNING: serving directly from a GraphStore can be slow; you may want to use a --serving_table")
		if f, ok := gs.(filetree.Service); ok {
			log.Printf("Using %T directly as filetree service", gs)
			ft = f
		} else {
			m := filetree.NewMap()
			if err := m.Populate(ctx, gs); err != nil {
				log.Fatalf("Error populating file tree from GraphStore: %v", err)
			}
			ft = m
		}

		if x, ok := gs.(xrefs.Service); ok {
			log.Printf("Using %T directly as xrefs service", gs)
			xs = x
		} else {
			if err := xstore.EnsureReverseEdges(ctx, gs); err != nil {
				log.Fatalf("Error ensuring reverse edges in GraphStore: %v", err)
			}
			xs = xstore.NewGraphStoreService(gs)
		}

		if s, ok := gs.(esearch.Service); ok {
			log.Printf("Using %T directly as entity search service", gs)
			esr = s
		}
	}

	if esr == nil {
		log.Println("Entity search API not supported")
	}
	if csr == nil {
		log.Println("Code search API not supported")
	}

	if *httpListeningAddr != "" {
		xrefs.RegisterHTTPHandlers(ctx, xs, http.DefaultServeMux)
		filetree.RegisterHTTPHandlers(ctx, ft, http.DefaultServeMux)
		if esr != nil {
			esearch.RegisterHTTPHandlers(ctx, esr, http.DefaultServeMux)
		}
		if csr != nil {
			codesearch.RegisterHTTPHandlers(ctx, csr, http.DefaultServeMux)
		}
		go startHTTP()
	}

	select {} // block forever
}

func startHTTP() {
	if *publicResources != "" {
		log.Println("Serving public resources at", *publicResources)
		if s, err := os.Stat(*publicResources); err != nil {
			log.Fatalf("ERROR: could not get FileInfo for %q: %v", *publicResources, err)
		} else if !s.IsDir() {
			log.Fatalf("ERROR: %q is not a directory", *publicResources)
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(*publicResources, filepath.Clean(r.URL.Path)))
		})
	}

	log.Printf("HTTP server listening on %q", *httpListeningAddr)
	log.Fatal(http.ListenAndServe(*httpListeningAddr, nil))
}
