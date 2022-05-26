package multimatcher

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Multimatcher of type [string,string]", func() {
	ginkgo.When("Creating call back function and a map labels", func() {
		cb := MatchCallback[string](func(string) string {
			return "test"
		})
		labels := map[string]string{
			"testKey1": "testValue1",
			"testKey2": "testValue2",
			"testKey3": "testValue3",
		}
		ginkgo.Context("creating new multimatcher and passing LabelExists func when adding entry", func() {
			mm := New[string, string]()
			err := mm.AddEntry("entry1", LabelExists("testKey1"), cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return 'test' upon match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
		ginkgo.Context("creating new multimatcher and passing LabelEquals func when adding entry", func() {
			mm := New[string, string]()
			err := mm.AddEntry("entry1", LabelEquals("testKey2", "testValue2"), cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return 'test' upon match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
		ginkgo.Context("creating new multimatcher and passing LabelMatchesRegex func when adding entry", func() {
			mm := New[string, string]()
			expr, err := LabelMatchesRegex("testKey3", "testValue3")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			err = mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return 'test' upon match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
		ginkgo.Context("creating new multimatcher and passing Not function when adding entry", func() {
			mm := New[string, string]()
			expr := Not(LabelExists("testKey1"))
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return '' upon negation", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal(""))
			})
		})
		ginkgo.Context("creating new multimatcher, array of expressions to test base case of Any func", func() {
			mm := New[string, string]()
			exprs := []Expr{
				LabelExists("testKey1"),
				LabelExists("testKey2"),
				LabelExists("testKey3"),
			}
			expr := Any(exprs)
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return 'test' upon match with any of the test keys", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
		ginkgo.Context("creating new multimatcher, array of expressions to test base case of All func", func() {
			mm := New[string, string]()
			exprs := []Expr{
				LabelExists("testKey1"),
				LabelExists("testKey2"),
				LabelExists("testKey3"),
				LabelExists("testKey4"),
			}
			expr := All(exprs)
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mm.Length()).To(gomega.Equal(1))

			ginkgo.It("should return '' upon missed match with testKey4", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal(""))
			})
		})
		ginkgo.Context("creating a new multimatcher and passing an empty array to check special case of All/Any func", func() {
			mm := New[string, string]()
			exprs := []Expr{}
			expr := All(exprs)
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ginkgo.It("should return 'test' upon empty match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal(""))
			})
			expr = Any(exprs)
			err = mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ginkgo.It("should return 'test' upon empty match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal(""))
			})
		})
		ginkgo.Context("creating a new multimatcher and passing a single element array to check special case of All/Any func", func() {
			mm := New[string, string]()
			exprs := []Expr{LabelExists("testKey1")}
			expr := All(exprs)
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ginkgo.It("should return 'test' upon empty match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
			expr = Any(exprs)
			err = mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			ginkgo.It("should return 'test' upon empty match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
		ginkgo.Context("creating a new multimatcher and passing an expression array to check All func ", func() {
			mm := New[string, string]()
			exprs := []Expr{LabelExists("testKey1"), LabelExists("testKey2"), LabelExists("testKey3")}
			expr := All(exprs)
			err := mm.AddEntry("entry1", expr, cb)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ginkgo.It("should return 'test' upon empty match with all the keys in labels", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal("test"))
			})
		})
	})
})

var _ = ginkgo.Describe("Multimatcher of type [string, int]", func() {
	ginkgo.When("Creating an int call back function and a map labels", func() {
		cb := MatchCallback[int](func(int) int {
			return 1
		})
		labels := map[string]string{
			"testKey1": "testValue1",
			"testKey2": "testValue2",
			"testKey3": "testValue3",
		}
		mm := New[string, int]()
		ginkgo.Context("passing multiple entries to multimatcher and multiple Label funcs to multimatcher", func() {
			ginkgo.It("should increase length when adding multiple entries", func() {
				err := mm.AddEntry("entry1", LabelExists("testKey1"), cb)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				err = mm.AddEntry("entry2", LabelEquals("testKey2", "testValue2"), cb)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mm.Length()).To(gomega.Equal(2))
			})
			ginkgo.It("should not change the length when overriding entry", func() {
				err := mm.AddEntry("entry1", LabelExists("testKey4"), cb)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mm.Length()).To(gomega.Equal(2))
			})
			ginkgo.It("should return 1 upon match", func() {
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal(1))
			})
			ginkgo.It("successfully remove an entry", func() {
				err := mm.RemoveEntry("entry1")
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(mm.Length()).To(gomega.Equal(1))
			})
		})
	})
})

var _ = ginkgo.Describe("Multimatcher of type [string, string]", func() {
	ginkgo.When("Creating a new multimatcher and a single label map", func() {
		labels := map[string]string{
			"testKey1": "testValue1",
		}
		mm := New[string, []string]()
		ginkgo.Context("When creating an entry with the append func", func() {
			ginkgo.It("should correctly test the appender func", func() {
				err := mm.AddEntry("entry1", LabelExists("testKey1"), Appender("test"))
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				s := mm.Match(Labels(labels))
				gomega.Expect(s).To(gomega.Equal([]string{"test"}))
			})
		})
		ginkgo.Context("When creating an irregular regex", func() {
			expr, err := LabelMatchesRegex("?=re", "?=re")
			ginkgo.It("should catch the error", func() {
				gomega.Expect(err).To(gomega.HaveOccurred())
				err := mm.AddEntry("entry1", expr, Appender("test"))
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})
		})
	})
})

var _ = ginkgo.Describe("Multimatcher of type [string, bool]", func() {
	ginkgo.When("Creating a new multimatcher and an empty labels map int", func() {
		mm := New[string, bool]()
		labels := map[string]string{}
		cb := MatchCallback[bool](func(bool) bool {
			return true
		})
		ginkgo.Context("When creating different types of entries", func() {
			mm.AddEntry("entry1", LabelEquals("testKey2", "testValue2"), cb)
			ginkgo.It("should return false upon mis match", func() {
				val := mm.Match(Labels(labels))
				gomega.Expect(val).To(gomega.Equal(false))
			})
			mm.AddEntry("", LabelExists(""), cb)
			ginkgo.It("should return true upon match", func() {
				val := mm.Match(Labels(labels))
				gomega.Expect(val).To(gomega.Equal(false))
			})
			exp, err := LabelMatchesRegex("testKey2", "testValue2")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			mm.AddEntry("entry2", exp, cb)
			ginkgo.It("should return false upon mis match", func() {
				val := mm.Match(Labels(labels))
				gomega.Expect(val).To(gomega.Equal(false))
			})
			exprs := []Expr{LabelExists("testKey1"), LabelEquals("testKey2", "testValue2")}
			expr := Any(exprs)
			mm.AddEntry("entry3", expr, cb)
			ginkgo.It("should return false upon mis match", func() {
				val := mm.Match(Labels(labels))
				gomega.Expect(val).To(gomega.Equal(false))
			})
			exprs = []Expr{exp}
			expr = All(exprs)
			mm.AddEntry("entry4", expr, cb)
			ginkgo.It("should return false upon mis match", func() {
				val := mm.Match(Labels(labels))
				gomega.Expect(val).To(gomega.Equal(false))
			})
		})
	})
})
