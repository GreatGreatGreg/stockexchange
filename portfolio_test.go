package stockexchange_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/stockexchange"
)

var _ = Describe("Portfolio", func() {
	var (
		stockOne  *stockexchange.Stock
		stockTwo  *stockexchange.Stock
		portfolio *stockexchange.Portfolio
	)

	BeforeEach(func() {
		stockOne = &stockexchange.Stock{
			Symbol:   "B",
			Name:     "Bengaza",
			AskPrice: 5,
			BidPrice: 10,
		}

		stockTwo = &stockexchange.Stock{
			Symbol:   "G",
			Name:     "Google",
			AskPrice: 3,
			BidPrice: 20,
		}

		portfolio = &stockexchange.Portfolio{
			Balance: 100,
		}
	})

	Describe("Buy", func() {
		var originalBalance float32

		BeforeEach(func() {
			originalBalance = portfolio.Balance
		})

		It("can buy stock", func() {
			Expect(portfolio.Buy(stockOne, 10)).To(Succeed())
			Expect(portfolio.Shares).To(HaveLen(1))

			share := portfolio.Shares[0]
			Expect(share.Symbol).To(Equal(stockOne.Symbol))
			Expect(share.Name).To(Equal(stockOne.Name))
			Expect(share.Price).To(Equal(stockOne.AskPrice))
			Expect(share.Quantity).To(Equal(10))

			Expect(portfolio.Balance).To(Equal(originalBalance - 10*stockOne.AskPrice))
		})

		It("can buy more that one stock", func() {
			Expect(portfolio.Buy(stockOne, 10)).To(Succeed())
			Expect(portfolio.Buy(stockTwo, 10)).To(Succeed())
			Expect(portfolio.Shares).To(HaveLen(2))
		})

		Context("when you buy the same stock", func() {
			It("increase the quantity of the share", func() {
				Expect(portfolio.Buy(stockOne, 10)).To(Succeed())
				Expect(portfolio.Buy(stockOne, 10)).To(Succeed())
				Expect(portfolio.Shares).To(HaveLen(1))

				share := portfolio.Shares[0]
				Expect(share.Symbol).To(Equal(stockOne.Symbol))
				Expect(share.Name).To(Equal(stockOne.Name))
				Expect(share.Price).To(Equal(stockOne.AskPrice))
				Expect(share.Quantity).To(Equal(20))

				Expect(portfolio.Balance).To(Equal(originalBalance - 20*stockOne.AskPrice))
			})
		})

		Context("when you do not have enough money", func() {
			It("returns an error", func() {
				Expect(portfolio.Buy(stockOne, 1000)).To(MatchError("Insufficient funds"))
			})
		})
	})
})
