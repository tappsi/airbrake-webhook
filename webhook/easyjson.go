package webhook

import jlexer  "github.com/mailru/easyjson/jlexer"
import jwriter "github.com/mailru/easyjson/jwriter"

func (v Notification) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson_encode(w, v)
}

func (v *Notification) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson_decode(l, v)
}

func easyjson_encode(out *jwriter.Writer, in Notification) {

	out.RawByte('{')

	out.String("service")
	out.RawByte(':')
	out.String(in.service)
	out.RawByte(',')

	out.String("recipient")
	out.RawString(":[")
	for index, value := range in.recipients {
		if index != 0 {
			out.RawByte(',')
		}
		out.String(value)
	}
	out.RawString("],")

	out.String("message")
	out.RawByte(':')
	out.String(in.message)

	out.RawByte('}')

}

func easyjson_decode(in *jlexer.Lexer, out *Notification) {
	if in.IsNull() {
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
}
