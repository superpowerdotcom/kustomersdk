package kustomersdk_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKustomersdk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kustomersdk Suite")
}
