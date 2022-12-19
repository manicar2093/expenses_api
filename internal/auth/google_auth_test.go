package auth_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoogleAuth", func() {

	Describe("Login", func() {
		It("takes google token to create token to login", func() {

		})

		When("user is not register yet", func() {
			It("creates a new user to db and creates a new login token", func() {

			})
		})
	})

	Describe("RefreshToken", func() {
		It("creates a new access token from given refresh token", func() {

		})
	})

})
