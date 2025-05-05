package main

func main() {

}


//1. Build Deck
// a. build out card object
//	faces, front back, suit and rank.
// b. build out deck object
//	shuffle, deal(since euchre you deal either 3 or 2 cards at a time), 
//2. Build out players
// a. player has a hand, tricks won, their position, whether or not they are the caller
// b. build out probabilities, the point system for calling
//3. build out UI using fyne
// a. play area, deck and hand
// b. table for the probabilities of certain players having certain cards.
// c. add functionality to the cards to allow gameplay
//4. Time permitting - use probabilities to make a computer opponent.
// a. Play cards based on probability, play to win points
// b. call based on the point value in hand. Loner logic