package integration_test

import (
	"net/http"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	var (
		args       []string
		session    *gexec.Session
		sessionErr error
	)

	JustBeforeEach(func() {
		session, sessionErr = runner.Start(args...)
	})

	AfterEach(func() {
		if session != nil {
			session.Kill()
		}
	})

	It("is listenting on HTTP port 9292", func() {
		Expect(sessionErr).NotTo(HaveOccurred())
		_, err := http.Get("http://127.0.0.1:9292")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when the addr is provided", func() {
		BeforeEach(func() {
			args = []string{"--addr=127.0.0.1:8080"}
		})

		It("is listenting on that HTTP port", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			_, err := http.Get("http://127.0.0.1:8080")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when the addr is wrong", func() {
		BeforeEach(func() {
			args = []string{"--addr=wrong_host_and_port"}
		})

		It("returns an error", func() {
			Expect(session).To(BeNil())
			Expect(sessionErr.Error()).To(ContainSubstring("The provided wrong_host_and_port addr is not correct"))
		})
	})
})
