package opax_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	Opax "github.com/w6d-io/opax"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Session", func() {
	Context("GET Session", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
			Opax.Opax = nil
		})
		It("succeeds to Get OPA Decision FromGRPCCtx to path << loopback >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			ctx := metadata.NewIncomingContext(ctx, metadata.MD{
				Opax.OpaDataName: []string{`{"path": "/v1/data/system/loopback", "input":{"foo": "bar"}}`},
			})
			decision, err := Opax.Opax.GetAuthorizationFromGRPCCtx(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision).ToNot(BeEmpty())
			Expect(decision).To(Equal("{\"result\":{\"foo\":\"bar\"}}"))
		})
		It("succeeds to Get OPA Decision FromGRPCCtx to path << str >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				Opax.OpaDataName: []string{`{"path": "/v1/data/system/str", "input":{"foo": "bar"}}`},
			})
			decision11, err := Opax.Opax.GetAuthorizationFromGRPCCtx(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision11).ToNot(BeEmpty())
			Expect(decision11).To(Equal("{\"result\":\"foo\"}"))
		})
		It("succeeds to Get OPA Decision FromGRPCCtx to path << main >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				Opax.OpaDataName: []string{`{"path": "/v1/data/system/main", "input":{"foo": "bar"}}`},
			})
			decision33, err := Opax.Opax.GetAuthorizationFromGRPCCtx(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision33).ToNot(BeEmpty())
			Expect(decision33).To(Equal("{\"result\":true}"))
		})
		It("Error to Get OPA Decision FromGRPCCtx to fake path", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				Opax.OpaDataName: []string{`{"path": "/foo/bar/foo", "input":{"foo": "bar"}}`},
			})
			decision22, err := Opax.Opax.GetAuthorizationFromGRPCCtx(ctx)
			Expect(err).To(HaveOccurred())
			Expect(decision22).To(BeEmpty())
		})
		It("succeeds to Get OPA Decision FromHTTP to path << str >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			var cfg interface{}

			json.Unmarshal([]byte(`{"path": "/v1/data/system/str", "input":{"foo": "bar"}}`), &cfg)
			decision1, err := Opax.Opax.GetAuthorizationFromHttp(ctx, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision1).ToNot(BeEmpty())
			Expect(decision1).To(Equal("{\"result\":\"foo\"}"))
		})
		It("succeeds to Get OPA Decision FromHTTP to path << main >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			var cfg interface{}

			json.Unmarshal([]byte(`{"path": "/v1/data/system/main", "input":{"foo": "bar"}}`), &cfg)
			decision5, err := Opax.Opax.GetAuthorizationFromHttp(ctx, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision5).ToNot(BeEmpty())
			Expect(decision5).To(Equal("{\"result\":true}"))
		})
		It("Error to Get OPA Decision FromHTTP to fake path", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			var cfg interface{}

			json.Unmarshal([]byte(`{"path": "/foo/bar/foo", "input":{"foo": "bar"}}`), &cfg)
			decision2, err := Opax.Opax.GetAuthorizationFromHttp(ctx, cfg)
			Expect(err).To(HaveOccurred())
			Expect(decision2).To(BeEmpty())
		})
		It("succeeds to Get OPA Decision FromHTTP to path << loopback >>", func() {
			Opax.SetOpaxDetails(false, "127.0.0.1", false, 8181)

			var cfg interface{}

			json.Unmarshal([]byte(`{"path": "/v1/data/system/loopback", "input":{"foo": "bar"}}`), &cfg)
			decision3, err := Opax.Opax.GetAuthorizationFromHttp(ctx, cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(decision3).ToNot(BeEmpty())
			Expect(decision3).To(Equal("{\"result\":{\"foo\":\"bar\"}}"))
		})
	})
})
