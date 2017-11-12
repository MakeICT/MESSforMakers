package main

// import (
// 	"github.com/onsi/ginkgo"
// 	"github.com/onsi/gomega"
// 	"github.com/sclevine/agouti/matchers"
// )

// var _ = ginkgo.Describe("Website", func() {
// 	ginkgo.It("displays hello world", func() {
// 		page := getPage()
// 		defer page.Destroy()

// 		gomega.Expect(page.Navigate("https://golang-with-chrome-skippotter.c9users.io:8080/")).To(gomega.Succeed())
// 		gomega.Eventually(page).Should(matchers.HaveURL("https://golang-with-chrome-skippotter.c9users.io:8080/"))
// 		gomega.Eventually(page.Find("h1")).Should(matchers.HaveText("You got the root handler !"))
// 	})
// })
