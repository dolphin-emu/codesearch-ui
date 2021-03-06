// Code generated by protoc-gen-go.
// source: codesearch/codesearch.proto
// DO NOT EDIT!

/*
Package codesearch is a generated protocol buffer package.

It is generated from these files:
	codesearch/codesearch.proto

It has these top-level messages:
	CodeSearchRequest
	Regexp
	CodeSearchReply
	Match
	Snippet
*/
package codesearch

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type CodeSearchRequest struct {
	Regexp     *Regexp `protobuf:"bytes,1,opt,name=regexp" json:"regexp,omitempty"`
	FileRegexp *Regexp `protobuf:"bytes,2,opt,name=file_regexp" json:"file_regexp,omitempty"`
	NumResults int32   `protobuf:"varint,3,opt,name=num_results" json:"num_results,omitempty"`
	Token      []byte  `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *CodeSearchRequest) Reset()         { *m = CodeSearchRequest{} }
func (m *CodeSearchRequest) String() string { return proto.CompactTextString(m) }
func (*CodeSearchRequest) ProtoMessage()    {}

func (m *CodeSearchRequest) GetRegexp() *Regexp {
	if m != nil {
		return m.Regexp
	}
	return nil
}

func (m *CodeSearchRequest) GetFileRegexp() *Regexp {
	if m != nil {
		return m.FileRegexp
	}
	return nil
}

type Regexp struct {
	Expr          string `protobuf:"bytes,1,opt,name=expr" json:"expr,omitempty"`
	CaseSensitive bool   `protobuf:"varint,2,opt,name=case_sensitive" json:"case_sensitive,omitempty"`
}

func (m *Regexp) Reset()         { *m = Regexp{} }
func (m *Regexp) String() string { return proto.CompactTextString(m) }
func (*Regexp) ProtoMessage()    {}

type CodeSearchReply struct {
	Match         []*Match `protobuf:"bytes,1,rep,name=match" json:"match,omitempty"`
	NextPageToken []byte   `protobuf:"bytes,2,opt,name=next_page_token,proto3" json:"next_page_token,omitempty"`
}

func (m *CodeSearchReply) Reset()         { *m = CodeSearchReply{} }
func (m *CodeSearchReply) String() string { return proto.CompactTextString(m) }
func (*CodeSearchReply) ProtoMessage()    {}

func (m *CodeSearchReply) GetMatch() []*Match {
	if m != nil {
		return m.Match
	}
	return nil
}

type Match struct {
	Filename string     `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
	Snippet  []*Snippet `protobuf:"bytes,2,rep,name=snippet" json:"snippet,omitempty"`
}

func (m *Match) Reset()         { *m = Match{} }
func (m *Match) String() string { return proto.CompactTextString(m) }
func (*Match) ProtoMessage()    {}

func (m *Match) GetSnippet() []*Snippet {
	if m != nil {
		return m.Snippet
	}
	return nil
}

type Snippet struct {
	Content    string `protobuf:"bytes,1,opt,name=content" json:"content,omitempty"`
	LineNumber int32  `protobuf:"varint,2,opt,name=line_number" json:"line_number,omitempty"`
}

func (m *Snippet) Reset()         { *m = Snippet{} }
func (m *Snippet) String() string { return proto.CompactTextString(m) }
func (*Snippet) ProtoMessage()    {}
