// Package srs implements generic spaced repetition system functions.
package srs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

/// /pkg/time specifically says this is ugly haha
// The spaced repetition delay table.
var levelDelay = map[int]time.Duration{
	0: time.Minute * 10,                   // 10 minutes
	1: time.Hour * 1,                      // 1 hour
	2: time.Hour * 5,                      // 5 hours
	3: time.Hour * 24,                     // 1 day
	4: 5 * (time.Hour * 24),               // 5 days
	5: 25 * (time.Hour * 24),              // 25 days
	6: 4 * (30 * (time.Hour * 24)),        // 4 months
	7: 2 * (12 * (30 * (time.Hour * 24))), // 2 years
}

// Open takes a file path and returns a homework struct and potential errors.
func Open(filePath string) (h *Homework, err error) {
	h = new(Homework)
	h.path = filePath

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	/// Might bug
	err = json.Unmarshal(buf, h)
	return h, err
}

// Homework is a list of vocabularies. It also contains the path to the homework
// file.
type Homework struct {
	Vocabs []Vocab // List of vocabularies to train on.
	path   string  // Path to the homework file.
}

// Vocab is a vocabulary item to train on.
type Vocab struct {
	Question   string // Question that needs an answer.
	Answer     string // The answer to compare with the users answer.
	ReviewDate int64  // Last time this vocabulary item was reviewed.
	Level      int    // Which level the user has managed to achieve for this vocabulary item.
}

// save updates the homework file after a quiz.
func (h *Homework) save() (err error) {
	buf, err := json.Marshal(h)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(h.path, buf, 0777)
	if err != nil {
		return err
	}

	return nil
}

// Quiz is an interactive shell function which asks the user to answer all
// questions from the homework file.
func (h *Homework) Quiz() (err error) {
	// Number of questions correctly answered.
	var nSuccess int

	// A list of questions answered incorrectly.
	var fails []Vocab

	// If all reviews are in the future, save the earliest time.
	var earliest int64

	// Loop through all vocabulary items and ask the user to answer the
	// questions.
	for key, v := range h.Vocabs {
		// If review date is in the future continue, and save the date
		// to show the user when the next update is if all review dates are in
		// the future.
		if v.ReviewDate > time.Now().Unix() {
			if v.ReviewDate < earliest || earliest == 0 {
				earliest = v.ReviewDate
			}
			continue
		}

		fmt.Println(v.Question + "?")

		var answer string
		fmt.Print("Answer: ")
		fmt.Scanf("%s", &answer)

		// If the question was correctly answered increase it's SRS level.
		// Otherwise decrease it and print the correct answer.
		if answer == v.Answer {
			fmt.Println("Success!\n")

			h.Vocabs[key].Level++
			nSuccess++
		} else {
			fmt.Println("False.")
			fmt.Println("Correct answer: ", v.Answer, "\n")

			fails = append(fails, v)
			if v.Level > 0 {
				h.Vocabs[key].Level--
			}
		}
		// Update the review time proportional of the vocabulary items new
		// level.
		nextTime := time.Now().Add(levelDelay[v.Level])
		h.Vocabs[key].ReviewDate = nextTime.Unix()
	}
	// If all reviews are in the future tell the user when the next review is.
	if nSuccess == 0 && len(fails) == 0 {
		fmt.Printf("The next review is in: %s.\n", time.Unix(earliest, 0).Sub(time.Now()).String())
		return nil
	}
	fmt.Println("You were correct", nSuccess, "times out of", len(h.Vocabs))

	// If the user answered some questions incorrectly print them so the user
	// can look through them one more time.
	if len(fails) != 0 {
		fmt.Println("Train some more on these:")
		for _, v := range fails {
			fmt.Println("\t", v.Answer, " - ", v.Question)
		}
	}

	// Update the homework file.
	err = h.save()
	if err != nil {
		return err
	}

	return nil
}

// Reset will reset all levels and review dates for a homework.
func (h *Homework) Reset() (err error) {
	for key, _ := range h.Vocabs {
		h.Vocabs[key].Level = 0
		h.Vocabs[key].ReviewDate = 0
	}
	return h.save()
}

/// This function is incomplete since it will save the homework in the
/// randomized order.
// RandomizeOrder randomizes the vocabularies order.
func (h *Homework) RandomizeOrder() {
	for i := range h.Vocabs {
		j := rand.Intn(i + 1)
		h.Vocabs[i], h.Vocabs[j] = h.Vocabs[j], h.Vocabs[i]
	}
}
