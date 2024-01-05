package r

import (
	"context"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	test := []struct {
		RawString  string
		DestString string
		From, To   byte
	}{
		{"X3H383c383t3r3k3q3T3y3c3h3f3l3O39", "X6H686c686t6r6k6q6T6y6c6h6f6l6O69", '3', '6'},
		{"pxTxLx5xIxZxwxCx9xXx9xHxrxxxoxQ", "p6T6L656I6Z6w6C696X696H6r666o6Q", 'x', '6'},
		{"9jTjFj5jgjQjRjvjnj8jgjCjZjpjpjzjvj5", "9mTmFm5mgmQmRmvmnm8mgmCmZmpmpmzmvm5", 'j', 'm'},
		{"2rvrIrvrHrhror4rWrf", "2wvwIwvwHwhwow4wWwf", 'r', 'w'},
		{"SqsqjqYqSqsqpqA", "SCsCjCYCSCsCpCA", 'q', 'C'},
		{"kiViviPig", "kWVWvWPWg", 'i', 'W'},
		{"ZKyK1KFKiKzKsKAKbKmKWK9KkKnKLKRKGK9", "Zfyf1fFfifzfsfAfbfmfWf9fkfnfLfRfGf9", 'K', 'f'},
		{"dmAm7m4mKmCm0mdmdmfm9memSmnmnmjmH", "dGAG7G4GKGCG0GdGdGfG9GeGSGnGnGjGH", 'm', 'G'},
		{"MBvBpBSBPBHBBBZB4BdBTBTBw", "MsvspsSsPsHsssZs4sdsTsTsw", 'B', 's'},
		{"dTrTWTeTWTTTSTzTtT3TYTLTM", "dArAWAeAWAAASAzAtA3AYALAM", 'T', 'A'},
		{"AaZaQavaHaxa0apaUa5aoa5aFaoaEai", "AtZtQtvtHtxt0tptUt5tot5tFtotEti", 'a', 't'},
		{"gOJO9OJObOpOvOYOmOyOROD", "gAJA9AJAbApAvAYAmAyARAD", 'O', 'A'},
		{"H202m2w2z2b2J", "Hs0smswszsbsJ", '2', 's'},
		{"GqoqQqyqWq3qPqv", "GRoRQRyRWR3RPRv", 'q', 'R'},
		{"fPlPgPLPYPTP2PxP9PXPpPOP4PiPcPaPh", "fUlUgULUYUTU2UxU9UXUpUOU4UiUcUaUh", 'P', 'U'},
	}
	ctx := context.Background()
	for i := 0; i < len(test); i++ {
		r := R{}
		r.RawString, r.From, r.To = test[i].RawString, []byte{test[i].From},
			[]byte{test[i].To}
		r.Churn(ctx)
		if r.DestString != test[i].DestString {
			fmt.Printf("expected %s. got %s\n", test[i].DestString,
				r.DestString)
		}
	}
}

//func TestvalRegexRange(t *testing.T) {
//	test := []struct {
//		regex     string
//		allValues string
//	}{
//		{
//			"A-Z", "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
//		},
//		{
//			"A-Z0-9", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
//		},
//	}
//	ctx := context.Background()
//	for i := 0; i < len(test); i++ {
//		valRegexRange()
//		[]byte{test[i].To}
//		r.Churn(ctx)
//		if r.DestString != test[i].DestString {
//			fmt.Printf("expected %s. got %s\n", test[i].DestString,
//				r.DestString)
//		}
//	}
//}

func BenchmarkR_Read(t *testing.B) {
	test := []struct {
		RawString  string
		DestString string
		From, To   byte
	}{
		{"X3H383c383t3r3k3q3T3y3c3h3f3l3O39", "X6H686c686t6r6k6q6T6y6c6h6f6l6O69", '3', '6'},
		{"pxTxLx5xIxZxwxCx9xXx9xHxrxxxoxQ", "p6T6L656I6Z6w6C696X696H6r666o6Q", 'x', '6'},
		{"9jTjFj5jgjQjRjvjnj8jgjCjZjpjpjzjvj5", "9mTmFm5mgmQmRmvmnm8mgmCmZmpmpmzmvm5", 'j', 'm'},
		{"2rvrIrvrHrhror4rWrf", "2wvwIwvwHwhwow4wWwf", 'r', 'w'},
		{"SqsqjqYqSqsqpqA", "SCsCjCYCSCsCpCA", 'q', 'C'},
		{"kiViviPig", "kWVWvWPWg", 'i', 'W'},
		{"ZKyK1KFKiKzKsKAKbKmKWK9KkKnKLKRKGK9", "Zfyf1fFfifzfsfAfbfmfWf9fkfnfLfRfGf9", 'K', 'f'},
		{"dmAm7m4mKmCm0mdmdmfm9memSmnmnmjmH", "dGAG7G4GKGCG0GdGdGfG9GeGSGnGnGjGH", 'm', 'G'},
		{"MBvBpBSBPBHBBBZB4BdBTBTBw", "MsvspsSsPsHsssZs4sdsTsTsw", 'B', 's'},
		{"dTrTWTeTWTTTSTzTtT3TYTLTM", "dArAWAeAWAAASAzAtA3AYALAM", 'T', 'A'},
		{"AaZaQavaHaxa0apaUa5aoa5aFaoaEai", "AtZtQtvtHtxt0tptUt5tot5tFtotEti", 'a', 't'},
		{"gOJO9OJObOpOvOYOmOyOROD", "gAJA9AJAbApAvAYAmAyARAD", 'O', 'A'},
		{"H202m2w2z2b2J", "Hs0smswszsbsJ", '2', 's'},
		{"GqoqQqyqWq3qPqv", "GRoRQRyRWR3RPRv", 'q', 'R'},
		{"fPlPgPLPYPTP2PxP9PXPpPOP4PiPcPaPh", "fUlUgULUYUTU2UxU9UXUpUOU4UiUcUaUh", 'P', 'U'},
	}
	ctx := context.Background()
	for i := 0; i < len(test); i++ {
		r := R{}
		r.RawString, r.From, r.To = test[i].RawString, []byte{test[i].From},
			[]byte{test[i].To}
		r.Churn(ctx)
		if r.DestString != test[i].DestString {
			fmt.Printf("expected %s. got %s\n", test[i].DestString,
				r.DestString)
		}
	}
}

func BenchmarkR_ReadSep(t *testing.B) {
	test := []struct {
		RawString  string
		DestString string
		From, To   byte
	}{
		{"X3H383c383t3r3k3q3T3y3c3h3f3l3O39", "X6H686c686t6r6k6q6T6y6c6h6f6l6O69", '3', '6'},
		{"pxTxLx5xIxZxwxCx9xXx9xHxrxxxoxQ", "p6T6L656I6Z6w6C696X696H6r666o6Q", 'x', '6'},
		{"9jTjFj5jgjQjRjvjnj8jgjCjZjpjpjzjvj5", "9mTmFm5mgmQmRmvmnm8mgmCmZmpmpmzmvm5", 'j', 'm'},
		{"2rvrIrvrHrhror4rWrf", "2wvwIwvwHwhwow4wWwf", 'r', 'w'},
		{"SqsqjqYqSqsqpqA", "SCsCjCYCSCsCpCA", 'q', 'C'},
		{"kiViviPig", "kWVWvWPWg", 'i', 'W'},
		{"ZKyK1KFKiKzKsKAKbKmKWK9KkKnKLKRKGK9", "Zfyf1fFfifzfsfAfbfmfWf9fkfnfLfRfGf9", 'K', 'f'},
		{"dmAm7m4mKmCm0mdmdmfm9memSmnmnmjmH", "dGAG7G4GKGCG0GdGdGfG9GeGSGnGnGjGH", 'm', 'G'},
		{"MBvBpBSBPBHBBBZB4BdBTBTBw", "MsvspsSsPsHsssZs4sdsTsTsw", 'B', 's'},
		{"dTrTWTeTWTTTSTzTtT3TYTLTM", "dArAWAeAWAAASAzAtA3AYALAM", 'T', 'A'},
		{"AaZaQavaHaxa0apaUa5aoa5aFaoaEai", "AtZtQtvtHtxt0tptUt5tot5tFtotEti", 'a', 't'},
		{"gOJO9OJObOpOvOYOmOyOROD", "gAJA9AJAbApAvAYAmAyARAD", 'O', 'A'},
		{"H202m2w2z2b2J", "Hs0smswszsbsJ", '2', 's'},
		{"GqoqQqyqWq3qPqv", "GRoRQRyRWR3RPRv", 'q', 'R'},
		{"fPlPgPLPYPTP2PxP9PXPpPOP4PiPcPaPh", "fUlUgULUYUTU2UxU9UXUpUOU4UiUcUaUh", 'P', 'U'},
	}
	ctx := context.Background()
	for i := 0; i < len(test); i++ {
		r := R{}
		r.RawString, r.From, r.To = test[i].RawString, []byte{test[i].From},
			[]byte{test[i].To}
		r.Churn(ctx)
		if r.DestString != test[i].DestString {
			fmt.Printf("expected %s. got %s\n", test[i].DestString,
				r.DestString)
		}
	}
}

func BenchmarkByteSliceEqual(t *testing.B) {} //against reflect.DeepEqual
