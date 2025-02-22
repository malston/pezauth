package pezdispenser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezdispenser/service"
)

var _ = Describe("NewMongoCollectionWrapper", func() {
	Context("when called with a collection", func() {
		It("Should return me a non nil Persistence interface", func() {
			wrap := NewMongoCollectionWrapper(new(mockMongo))
			Ω(wrap).ShouldNot(BeNil())
		})
	})

	Context("calling .Upsert when collection yeilds no error", func() {
		It("Should return ", func() {
			wrap := NewMongoCollectionWrapper(new(mockMongo))
			Ω(wrap.Upsert(nil, nil)).ShouldNot(Equal(ErrCanNotAddOrgRec))
		})
	})

	Context("calling .Upsert when collection yeilds error", func() {
		It("Should return ", func() {
			wrap := NewMongoCollectionWrapper(&mockMongo{
				err: errors.New("my mock error"),
			})
			Ω(wrap.Upsert(nil, nil)).ShouldNot(BeNil())
			Ω(wrap.Upsert(nil, nil)).Should(Equal(ErrCanNotAddOrgRec))
		})
	})

	Context("calling .Remove when collection yeilds no error", func() {
		It("Should return ", func() {
			wrap := NewMongoCollectionWrapper(new(mockMongo))
			Ω(wrap.Remove(nil)).Should(BeNil())
		})
	})

	Context("calling .Remove when collection yeilds error", func() {
		It("Should return ", func() {
			controlErr := errors.New("my mock error")
			wrap := NewMongoCollectionWrapper(&mockMongo{
				err: controlErr,
			})
			Ω(wrap.Remove(nil)).ShouldNot(BeNil())
			Ω(wrap.Remove(nil)).Should(Equal(controlErr))
		})
	})

})
