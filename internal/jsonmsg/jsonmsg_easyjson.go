// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsonmsg

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg(in *jlexer.Lexer, out *Output) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "dataArray":
			if in.IsNull() {
				in.Skip()
				out.Texts = nil
			} else {
				in.Delim('[')
				if out.Texts == nil {
					if !in.IsDelim(']') {
						out.Texts = make([]string, 0, 4)
					} else {
						out.Texts = []string{}
					}
				} else {
					out.Texts = (out.Texts)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Texts = append(out.Texts, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg(out *jwriter.Writer, in Output) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"dataArray\":"
		out.RawString(prefix[1:])
		if in.Texts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Texts {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Output) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Output) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Output) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Output) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg(l, v)
}
func easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg1(in *jlexer.Lexer, out *Input) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			out.Text = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg1(out *jwriter.Writer, in Input) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Input) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Input) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32b970dbEncodeGithubComErupshisKodeGitInternalJsonmsg1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Input) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Input) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32b970dbDecodeGithubComErupshisKodeGitInternalJsonmsg1(l, v)
}
