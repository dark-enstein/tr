package r

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
)

const (
	Action_DELETE = iota + 1
	Action_SQUEEZE
)

var (
	VALHELP = `
1). Ensure that your regex search and replace expressions are properly
delimited by hyphen (eg: A-Z a-z)
2) Ensure that you're only comparing one ASCII literall against the other
`
	PosixBracRegexMap = map[string]string{
		"[:upper:]": "A-Z",
		"[:lower:]": "a-z",
	}
	DelString = ""
)

// R defines the structure to be used for storing the state of tr execution
type R struct {
	// RawString represents the string as read from stdin / file input
	RawString string
	// DestString holds the processed string following replacement/processing
	DestString string
	// RawBytes hold the bytes representation of RawString
	RawBytes []byte
	// From and To, holds the binary representation of the search string and
	//the replace string respectively
	From, To []byte
	// FlagEnabled indicates whether a flag has been enabled or not.
	FlagEnabled bool
	// Flags defines the flags that can be set during starttime
	Flag *Flags
	// Embedded struct to control mutation of struct resource
	sync.Mutex
}

// Flags defines the flag object that holds all the binary options at runtime
type Flags struct {
	// DelString defines the string to delete
	DelString string
	// SqueezeByte defines the []byte char to squeeze
	SqueezeByte []byte
	// SqueezeString defines the string char to squeeze
	SqueezeString string
	// Action defines what mode of flag action is enabled
	Action int
}

// Churn processes the RawString in r,
// and perform the replacement operations as defined by the user
func (r *R) Churn(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if r.RawString == "" {
		log.Println("raw string not passed in yet")
	}
	if string(r.RawBytes) == "" {
		r.RawBytes = []byte(r.RawString)
	}
	if r.FlagEnabled {
		switch r.Flag.Action {
		case Action_DELETE:
			r.Delete(ctx)
			return
		case Action_SQUEEZE:
			r.Squeeze(ctx)
			return
		}
	}
	if len(r.From) > 1 {
		if (bytes.Contains(r.From, []byte("-")) && len(r.From) == 3) || (bytes.Contains(r.From, []byte(":")) && len(bytes.Split(r.From,
			[]byte(":"))) == 3) {
			if n := r.ResolveRegexArg(); n != 0 {
				return
			}
			if bytes.Contains(r.To, []byte("-")) && len(r.To) == 3 {
				r.ReplaceRange(ctx)
			} else {
				log.Printf("error parsing out replace range regex: %s\n", string(r.To))
			}
		} else {
			r.ReplaceSlice()
		}
	} else {
		r.Replace()
	}
}

// Replace replaces the portion of the input slice RawBytes that matches the search
// bytes From with the replace bytes To in-place.
// If a byte in RawBytes matches the first byte of From,
// it replaces that byte with To (considering To as a whole slice).
// DestString is updated with the new value of RawBytes.
func (r *R) Replace() {
	for i := 0; i < len(r.RawBytes); {
		if r.RawBytes[i] == r.From[0] {
			r.RawBytes = append(r.RawBytes[:i], append(r.To, r.RawBytes[i+1:]...)...)
			i += len(r.To)
		} else {
			i++
		}
	}
	r.DestString = string(r.RawBytes)
}

// ReplaceSlice is used when the From length is more than 1. ReplaceSlice replaces the portion of/the
// input/slice RawBytes/that matches the search slice bytes From with the replace bytes To in-place.
func (r *R) ReplaceSlice() {
	// Preallocate a buffer to avoid frequent reallocations
	buffer := make([]byte, 0, len(r.RawBytes)) // Initial capacity can be tuned based on expected final size

	i := 0
	for i < len(r.RawBytes) {
		if i+len(r.From) <= len(r.RawBytes) && ByteSliceEqual(r.RawBytes[i:i+len(r.From)],
			r.From) {
			buffer = append(buffer, r.To...)
			i += len(r.From)
		} else {
			buffer = append(buffer, r.RawBytes[i])
			i++
		}
	}
	r.RawBytes = buffer
	r.DestString = string(r.RawBytes)
}

// ReplaceRange
func (r *R) ReplaceRange(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var err error

	// elaborate ASCII compare ensuring that the range are within bounds of A-Z and a-z
	r.From, err = resolveRange(r.From)
	if err != nil {
		log.Printf("error occured: %s\n", err.Error())
		return
	}
	r.To, err = resolveRange(r.To)
	if err != nil {
		log.Printf("error occured:  %s\n", err.Error())
		return
	}
	//r.legacyRangeMutate()

	// I'm probably not handling this cancel op the proper way. TODO
	verdict := r.RangeMutate(func() {
		cancel()
	})

	if verdict == 0 {
		r.DestString = string(r.RawBytes)
	}
}

// DeleteRange
func (r *R) DeleteRange(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var err error

	// elaborate ASCII compare ensuring that the range are within bounds of A-Z and a-z
	r.From, err = resolveRange(r.From)
	if err != nil {
		log.Printf("error occured: %s\n", err.Error())
		return
	}

	r.DeleteOne(ctx)
}

// legacyRangeMutate implements ReplaceRange the way it is currently
// implemented in tr linux utility. It checks for the index of all matches with search string in original
// string. For each match index, it directly queries that index in the replace range To regardless of if
// the string replace ment is valid (eg: replacing t with D in toad >> Doad.)
// More robust implementation in RangeWinnow below
func (r *R) legacyRangeMutate() {
	for i := 0; i < len(r.RawBytes); i++ {
		go func(i int) {
			for j := 0; j < len(r.From); j++ {
				if r.RawBytes[i] == r.From[j] {
					if j > len(r.To) {
						r.Lock()
						r.RawBytes[i] = r.To[len(r.To)-1]
						r.Unlock()
					} else {
						r.Lock()
						r.RawBytes[i] = r.To[j]
						r.Unlock()
					}
				}
			}
		}(i)
	}
}

// RangeMutate implements ReplaceRange with certain safe restrictions.
// It enforces that the search and the replace ranges must be equal.
// Thus it can safely do a direct index search in the replacement array
func (r *R) RangeMutate(ctxFunc context.CancelFunc) int {
	// Check if the min and max of either range is the same
	if len(r.From) == 1 || len(r.To) == 1 {
		switch {
		case len(r.From) == 1:
			log.Printf("Incorrect search string: %s\n",
				fmt.Sprintf("%s-%s",
					string(r.From[0]),
					string(r.From[len(r.From)-1])))
			return 1
		case len(r.To) == 1:
			log.Printf("Incorrect replace string: %s\n",
				fmt.Sprintf("%s-%s",
					string(r.To[0]),
					string(r.To[len(r.To)-1])))
			return 1
		}
	}
	// Checks that the search and the replace string is of the same length
	if len(r.From) != len(r.To) {
		// this is possibly too much, optimize. TODO
		log.Printf("Search range %s, "+
			"is not the same length as replace range %s. "+
			"It must be of the same length\n", fmt.Sprintf("%s-%s",
			string(r.From[0]),
			string(r.From[len(r.From)-1])),
			fmt.Sprintf("%s-%s",
				string(r.To[0]),
				string(r.To[len(r.To)-1])))
		ctxFunc()
		return 1
	}
	for i := 0; i < len(r.RawBytes); i++ {
		go func(i int) {
			for j := 0; j < len(r.From); j++ {
				if r.RawBytes[i] == r.From[j] {
					if j > len(r.To) {
						r.Lock()
						r.RawBytes[i] = r.To[len(r.To)-1]
						r.Unlock()
					} else {
						r.Lock()
						r.RawBytes[i] = r.To[j]
						r.Unlock()
					}
				}
			}
		}(i)
	}
	return 0
}

// resolveRange resolve a range passed in into its individual bytes.
// Potentially be removed.
func resolveRange(b []byte) ([]byte, error) {
	var rangee = []byte{}
	if len(b) != 3 && b[1] != []byte("-")[0] {
		return []byte(""), fmt.Errorf("err: could not process byte, "+
			"not in right format: %s\n", b)
	}
	asciiDiff := rune(b[2]) - rune(b[0])
	in := rune(b[0])
	rangee = append(rangee, []byte(string(in))...)
	in = in + 1
	for i := 0; i < int(asciiDiff); i++ {
		rangee = append(rangee, []byte(string(in))...)
		in = in + 1
	}
	return rangee, nil
}

// valRegexRange handles general regex syntax errors: the inclusion of
// hyphen, a byte length of 3,
// and that the ascii value of the maxRange is greater than that of the minRange
func valRegexRange(b []byte) bool {
	if !(bytes.Contains(b, []byte("-")) && len(b) == 3 && (rune(
		b[2]) > rune(b[0]))) {
		log.Printf("error: incorrect regex range provided. "+
			"Check the following:\n%s\n", VALHELP)
		return false
	}
	return true
}

// ByteSliceEqual compares two byte slices, returning true if they are equal,
// and false if they aren't.
func ByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ResolveRegexArg consolidates validation of arguments for regex range
// options. It returns 0 if successful,
// and >0 if an error was encountered along the way
func (r *R) ResolveRegexArg() int {
	fmt.Printf("from: %v %d, to: %v %d\n", PosixBracRegexMap[string(r.From)], len(bytes.Split(r.From,
		[]byte(":"))),
		PosixBracRegexMap[string(r.To)], len(bytes.Split(r.To,
			[]byte(":"))))
	var errNo = 0
	if bytes.Contains(r.From, []byte(":")) && len(bytes.Split(r.From,
		[]byte(":"))) == 3 {
		log.Println("reached inner about map")
		val, ok := PosixBracRegexMap[string(r.From)]
		if !ok {
			errNo++
		}
		r.From = []byte(val)
	}
	if bytes.Contains(r.To, []byte(":")) && len(bytes.Split(r.To,
		[]byte(":"))) == 3 {
		log.Println("reached inner about map")
		val, ok := PosixBracRegexMap[string(r.To)]
		if !ok {
			errNo++
		}
		r.To = []byte(val)
	}
	if errNo == 0 && valRegexRange(r.From) && valRegexRange(r.To) {
		// passed validation. keep things the same.
		return 0
	} else {
		log.Printf("error validating regex")
		return 1
	}
}

// Delete deletes the specified string from the input text.
// This deletion happens in-place
func (r *R) Delete(ctx context.Context) {
	if bytes.Contains([]byte(r.Flag.DelString),
		[]byte(":")) && len(bytes.Split([]byte(r.Flag.DelString),
		[]byte(":"))) == 3 {
		var errNo int
		val, ok := PosixBracRegexMap[r.Flag.DelString]
		if !ok {
			errNo++
		} else {
			r.From = []byte(val)
		}

		r.DeleteRange(ctx)
	} else {
		i := 0
		for i < len(r.RawBytes) {
			if (len(r.Flag.DelString) < len(r.RawBytes[i:])) && (ByteSliceEqual(
				[]byte(r.Flag.DelString), r.RawBytes[i:i+len(r.Flag.
					DelString)])) {
				r.RawBytes = append(r.RawBytes[:i], r.RawBytes[i+len(r.Flag.
					DelString):]...)
				i += len(DelString)
			} else {
				i++
			}
		}
	}
	fmt.Printf("Deleted resp: %s\n", string(r.RawBytes))
	r.DestString = string(r.RawBytes)
}

// DeleteOne deletes every string char in the input text,
// as defined by the range expression.
func (r *R) DeleteOne(ctx context.Context) {
	// Preallocate a buffer to avoid frequent reallocations
	var buffer []byte
	delString := r.From
	for i := 0; i < len(r.RawBytes); i++ {
		shouldDelete := false
		for j := 0; j < len(delString); j++ {
			if r.RawBytes[i] == delString[j] {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			buffer = append(buffer, r.RawBytes[i])
		}
	}
	r.RawBytes = buffer
}

// Squeeze reduces repeated occurrences of defined char into once, within the
// input text.
func (r *R) Squeeze(ctx context.Context) {
	// Preallocate a buffer to avoid frequent reallocations
	var buffer []byte
	for i := 0; i < len(r.RawBytes); {
		for j := 0; j < len(r.Flag.SqueezeByte); j++ {
			if r.RawBytes[i] == r.Flag.SqueezeByte[j] {
				if r.RawBytes[i] == r.Flag.SqueezeByte[j] {
					n := endRepeated(r.RawBytes[i:])
					buffer = append(buffer, r.RawBytes[i])
					i += n
				}
			} else {
				buffer = append(buffer, r.RawBytes[i])
				i++
			}
		}
	}
}

// endRepeated checks how many indices from byt's first value has repeated
// values.
func endRepeated(byt []byte) int {
	cons := byt[0]
	for i := 1; i < len(byt); i++ {
		if cons == byt[i] {
			i++
		}
		if cons != byt[i+1] {
			return i
		}
	}
	// unlikely
	return 1
}
