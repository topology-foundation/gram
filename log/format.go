// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package log

import (
	"bytes"
	"log/slog"
)

const (
	timeFormat  = "2006-01-02T15:04:05"
	termMsgJust = 40
)

func (h *Handler) format(buf []byte, r slog.Record, usecolor bool) []byte {
	var color = ""
	if usecolor {
		switch r.Level {
		case slog.LevelError:
			color = "\x1b[31m"
		case slog.LevelWarn:
			color = "\x1b[33m"
		case slog.LevelInfo:
			color = "\x1b[32m"
		case slog.LevelDebug:
			color = "\x1b[36m"
		}
	}

	if buf == nil {
		buf = make([]byte, 0, 30+termMsgJust)
	}
	b := bytes.NewBuffer(buf)

	if color != "" { // Start color
		b.WriteString(color)
		b.WriteString(LevelString(r.Level))
		b.WriteString("\x1b[0m")
	} else {
		b.WriteString(LevelString(r.Level))
	}
	b.WriteString("[")
	b.WriteString(r.Time.Format(timeFormat))
	b.WriteString("] ")
	b.WriteString(r.Message)

	h.formatAttributes(b, r, color)
	return b.Bytes()
}

func (h *Handler) formatAttributes(buf *bytes.Buffer, r slog.Record, color string) {
	writeAttr := func(attr slog.Attr) {
		buf.WriteByte(' ')

		if color != "" {
			buf.WriteString(color)
			buf.WriteString(attr.Key)
			buf.WriteString("\x1b[0m=")
		} else {
			buf.WriteString(attr.Key)
			buf.WriteByte('=')
		}

		buf.WriteString(attr.Value.String())
	}

	r.Attrs(func(attr slog.Attr) bool {
		writeAttr(attr)
		return true
	})

	buf.WriteByte('\n')
}
