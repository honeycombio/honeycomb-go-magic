package honeycomb

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"time"

	libhoney "github.com/honeycombio/libhoney-go"
)

func InstrumentHandleFunc(hf func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO find out if we're a sub-handler and don't stomp the parent
		// event, or at least get parent/child IDs and intentionally send a
		// subevent or something
		start := time.Now()
		ev := existingEventFromContext(r.Context())
		if ev == nil {
			ev = libhoney.NewEvent()
			// put the event on the context for everybody downsteam to use
			r = r.WithContext(context.WithValue(r.Context(), honeyEventContextKey, ev))
		}
		// add some common fields from the request to our event
		addRequestProps(r, ev)
		// replace the writer with our wrapper to catch the status code
		wrappedWriter := &hnyResponseWriter{ResponseWriter: w}
		// add the name of the handler func we're about to invoke
		ev.AddField("handler_func_name", runtime.FuncForPC(reflect.ValueOf(hf).Pointer()).Name())

		hf(wrappedWriter, r)
		if wrappedWriter.status == 0 {
			wrappedWriter.status = 200
		}
		ev.AddField("response.status_code", wrappedWriter.status)
		ev.AddField("duration_ms", float64(time.Since(start))/float64(time.Millisecond))
		ev.Send()
	}
}

func InstrumentMuxHandler(mux *http.ServeMux) http.Handler {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		// TODO find out if we're a sub-handler and don't stomp the parent
		// event, or at least get parent/child IDs and intentionally send a
		// subevent or something
		start := time.Now()
		ev := existingEventFromContext(r.Context())
		if ev == nil {
			ev = libhoney.NewEvent()
			// put the event on the context for everybody downsteam to use
			r = r.WithContext(context.WithValue(r.Context(), honeyEventContextKey, ev))
		}
		// add some common fields from the request to our event
		addRequestProps(r, ev)
		// replace the writer with our wrapper to catch the status code
		wrappedWriter := &hnyResponseWriter{ResponseWriter: w}
		handler, pat := mux.Handler(r)
		ev.AddField("handlerPattern", pat)
		// get Handler type and name
		hType := reflect.TypeOf(handler)
		ev.AddField("handlerType", hType.String())
		ev.AddField("handler_name", runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		handler.ServeHTTP(wrappedWriter, r)
		if wrappedWriter.status == 0 {
			wrappedWriter.status = 200
		}
		ev.AddField("response.status_code", wrappedWriter.status)
		ev.AddField("duration_ms", float64(time.Since(start))/float64(time.Millisecond))
		ev.Send()
	}
	return http.HandlerFunc(wrappedHandler)
}

func InstrumentHandler(handler http.Handler) http.Handler {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// TODO find out if we're a sub-handler and don't stomp the parent
		// event, or at least get parent/child IDs and intentionally send a
		// subevent or something
		ev := existingEventFromContext(r.Context())
		if ev == nil {
			ev = libhoney.NewEvent()
			// put the event on the context for everybody downsteam to use
			r = r.WithContext(context.WithValue(r.Context(), honeyEventContextKey, ev))
		}
		// add some common fields from the request to our event
		addRequestProps(r, ev)
		// replace the writer with our wrapper to catch the status code
		wrappedWriter := &hnyResponseWriter{ResponseWriter: w}
		// add the name of the handler func we're about to invoke
		ev.AddField("handler_name", runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		handler.ServeHTTP(wrappedWriter, r)
		if wrappedWriter.status == 0 {
			wrappedWriter.status = 200
		}
		ev.AddField("response.status_code", wrappedWriter.status)
		ev.AddField("duration_ms", float64(time.Since(start))/float64(time.Millisecond))
		ev.Send()
	}
	return http.HandlerFunc(wrappedHandler)
}