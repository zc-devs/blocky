package stringcache_test

import (
	"github.com/0xERR0R/blocky/cache/stringcache"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"
)

var _ = Describe("In-Memory grouped cache", func() {
	var (
		cache   *stringcache.InMemoryGroupedCache
		factory stringcache.GroupFactory
	)

	Describe("Empty cache", func() {
		BeforeEach(func() {
			cache = stringcache.NewInMemoryGroupedStringCache()
		})
		When("empty cache was created", func() {
			It("should have element count of 0", func() {
				Expect(cache.ElementCount("someGroup")).Should(BeNumerically("==", 0))
			})

			It("should not find any string", func() {
				Expect(cache.Contains("searchString", []string{"someGroup"})).Should(BeEmpty())
			})
		})
		When("cache with one empty group", func() {
			BeforeEach(func() {
				factory = cache.Refresh("group1")
				factory.Finish()
			})

			It("should have element count of 0", func() {
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 0))
			})

			It("should not find any string", func() {
				Expect(cache.Contains("searchString", []string{"group1"})).Should(BeEmpty())
			})
		})
	})
	Describe("Cache creation", func() {
		When("cache with 1 group was created", func() {
			BeforeEach(func() {
				cache = stringcache.NewInMemoryGroupedStringCache()
				factory = cache.Refresh("group1")

				Expect(factory.AddEntry("string1")).Should(BeTrue())
				Expect(factory.AddEntry("string2")).Should(BeTrue())
			})

			It("cache should still have 0 element, since finish was not executed", func() {
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 0))
			})

			It("factory has 2 elements", func() {
				Expect(factory.Count()).Should(BeNumerically("==", 2))
			})

			It("should have element count of 2", func() {
				factory.Finish()
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 2))
			})

			It("should find strings", func() {
				factory.Finish()
				Expect(maps.Keys(cache.Contains("string1", []string{"group1"}))).Should(ConsistOf("group1"))
				Expect(maps.Keys(cache.Contains("string2", []string{"group1", "someOtherGroup"}))).Should(ConsistOf("group1"))
			})
		})
		When("Regex grouped cache is used", func() {
			BeforeEach(func() {
				cache = stringcache.NewInMemoryGroupedRegexCache()
				factory = cache.Refresh("group1")

				Expect(factory.AddEntry("string1")).Should(BeFalse())
				Expect(factory.AddEntry("/string2/")).Should(BeTrue())
				factory.Finish()
			})

			It("should ignore non-regex", func() {
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 1))
				Expect(maps.Keys(cache.Contains("string1", []string{"group1"}))).Should(BeEmpty())
				Expect(maps.Keys(cache.Contains("string2", []string{"group1"}))).Should(ConsistOf("group1"))
				Expect(maps.Keys(cache.Contains("shouldalsomatchstring2", []string{"group1"}))).Should(ConsistOf("group1"))
			})
		})
		When("Wildcard grouped cache is used", func() {
			BeforeEach(func() {
				cache = stringcache.NewInMemoryGroupedWildcardCache()
				factory = cache.Refresh("group1")

				Expect(factory.AddEntry("string1")).Should(BeFalse())
				Expect(factory.AddEntry("/string2/")).Should(BeFalse())
				Expect(factory.AddEntry("*.string3")).Should(BeTrue())
				factory.Finish()
			})

			It("should ignore non-wildcard", func() {
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 1))
				Expect(maps.Keys(cache.Contains("string1", []string{"group1"}))).Should(BeEmpty())
				Expect(maps.Keys(cache.Contains("string2", []string{"group1"}))).Should(BeEmpty())
				Expect(maps.Keys(cache.Contains("string3", []string{"group1"}))).Should(ConsistOf("group1"))
				Expect(maps.Keys(cache.Contains("shouldalsomatch.string3", []string{"group1"}))).Should(ConsistOf("group1"))
			})
		})
	})

	Describe("Cache refresh", func() {
		When("cache with 2 groups was created", func() {
			BeforeEach(func() {
				cache = stringcache.NewInMemoryGroupedStringCache()
				factory = cache.Refresh("group1")

				Expect(factory.AddEntry("g1")).Should(BeTrue())
				Expect(factory.AddEntry("both")).Should(BeTrue())
				factory.Finish()

				factory = cache.Refresh("group2")
				Expect(factory.AddEntry("g2")).Should(BeTrue())
				Expect(factory.AddEntry("both")).Should(BeTrue())
				factory.Finish()
			})

			It("should contain 4 elements in 2 groups", func() {
				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 2))
				Expect(cache.ElementCount("group2")).Should(BeNumerically("==", 2))
				Expect(maps.Keys(cache.Contains("g1", []string{"group1", "group2"}))).Should(ConsistOf("group1"))
				Expect(maps.Keys(cache.Contains("g2", []string{"group1", "group2"}))).Should(ConsistOf("group2"))
				Expect(maps.Keys(cache.Contains("both", []string{"group1", "group2"}))).Should(ConsistOf("group1", "group2"))
			})

			It("Should replace group content on refresh", func() {
				factory = cache.Refresh("group1")
				Expect(factory.AddEntry("newString")).Should(BeTrue())
				factory.Finish()

				Expect(cache.ElementCount("group1")).Should(BeNumerically("==", 1))
				Expect(cache.ElementCount("group2")).Should(BeNumerically("==", 2))
				Expect(maps.Keys(cache.Contains("g1", []string{"group1", "group2"}))).Should(BeEmpty())
				Expect(maps.Keys(cache.Contains("newString", []string{"group1", "group2"}))).Should(ConsistOf("group1"))
				Expect(maps.Keys(cache.Contains("g2", []string{"group1", "group2"}))).Should(ConsistOf("group2"))
				Expect(maps.Keys(cache.Contains("both", []string{"group1", "group2"}))).Should(ConsistOf("group2"))
			})
		})
	})
})
