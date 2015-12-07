// A generic way of creating a url.Values object from any struct. Uses the
// go generate command from Go 1.4. To enable, install this command and
// then add the following comment anywhere in the .go file that contains
// your struct:
//
//     //go:generate urlvaluer $GOFILE
//
// Next time you run "go generate" for these files, it will generate a file
// with the suffix ".urlvaluer.go" in that same directory which contains
// the relevant code to generate url.Values from your structs.
//
// LICENSE:
// Copyright 2014 Brett Slatkin, Parker Moore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main
