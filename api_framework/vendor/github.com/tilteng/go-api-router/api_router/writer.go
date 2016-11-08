package api_router

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
	SetStatus(int)
	WriteStatusHeader()
	Status() int
	Size() int
}

type baseResponseWriter struct {
	http.ResponseWriter
	defaultStatus int
	statusWritten bool
	status        int
	size          int
}

func (self *baseResponseWriter) writeStatusHeader() {
	if self.status == 0 {
		self.status = self.defaultStatus
	}
	self.ResponseWriter.WriteHeader(self.status)
	self.statusWritten = true
}

func (self *baseResponseWriter) Write(b []byte) (int, error) {
	if !self.statusWritten {
		self.writeStatusHeader()
	}
	size, err := self.ResponseWriter.Write(b)
	self.size += size
	return size, err
}

func (self *baseResponseWriter) WriteStatusHeader() {
	if !self.statusWritten {
		self.writeStatusHeader()
	}
}

func (self *baseResponseWriter) WriteHeader(s int) {
	if !self.statusWritten {
		self.status = s
		self.writeStatusHeader()
	}
}

func (self *baseResponseWriter) Status() int {
	return self.status
}

func (self *baseResponseWriter) Size() int {
	return self.size
}

func (self *baseResponseWriter) SetStatus(status int) {
	if !self.statusWritten {
		self.status = status
	}
}

type flushWriter struct {
	ResponseWriter
	http.Flusher
}

type hijackWriter struct {
	ResponseWriter
	http.Hijacker
}

type closeNotifyWriter struct {
	ResponseWriter
	http.CloseNotifier
}

type hijackFlushWriter struct {
	ResponseWriter
	http.Flusher
	http.Hijacker
}

type hijackCloseNotifyWriter struct {
	ResponseWriter
	http.Hijacker
	http.CloseNotifier
}

type closeNotifyFlushWriter struct {
	ResponseWriter
	http.CloseNotifier
	http.Flusher
}

type allTheThingsWriter struct {
	ResponseWriter
	http.Flusher
	http.Hijacker
	http.CloseNotifier
}

func newResponseWriter(w http.ResponseWriter, default_status int) ResponseWriter {
	base_writer := &baseResponseWriter{
		ResponseWriter: w,
		defaultStatus:  default_status,
	}

	flusher, flusher_ok := w.(http.Flusher)
	hijacker, hijacker_ok := w.(http.Hijacker)
	close_notifier, close_notifier_ok := w.(http.CloseNotifier)

	if flusher_ok {
		if hijacker_ok && close_notifier_ok {
			return allTheThingsWriter{
				base_writer,
				flusher,
				hijacker,
				close_notifier,
			}
		}

		if close_notifier_ok {
			return closeNotifyFlushWriter{
				base_writer,
				close_notifier,
				flusher,
			}
		}

		if hijacker_ok {
			return hijackFlushWriter{
				base_writer,
				flusher,
				hijacker,
			}
		}

		return flushWriter{
			base_writer,
			flusher,
		}
	}

	if close_notifier_ok {
		if hijacker_ok {
			return hijackCloseNotifyWriter{
				base_writer,
				hijacker,
				close_notifier,
			}
		}

		return closeNotifyWriter{
			base_writer,
			close_notifier,
		}
	}

	if hijacker_ok {
		return hijackWriter{
			base_writer,
			hijacker,
		}
	}

	return base_writer
}
